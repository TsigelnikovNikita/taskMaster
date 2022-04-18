package config_parser

import (
	"bytes"
	"gopkg.in/ini.v1"
	"main/task"
	"testing"
)

type mockConfigReader struct {
	data []byte
}

func (m mockConfigReader) getData() ([]byte, error) {
	return m.data, nil
}

func serializeStructToBytes(v interface{}) []byte {
	cfg := ini.Empty()
	_ = ini.ReflectFrom(cfg, v)

	data := new(bytes.Buffer)
	cfg.WriteTo(data)
	return data.Bytes()
}

func TestWithAllInformation(t *testing.T) {
	type T struct {
		*task.Task `comment:"Task"`
	}
	expectedResult :=
		&task.Task{
			Name:             "Task",
			Command:          "ls -la",
			ProcessNumber:    10,
			AutoStart:        true,
			StartTimeSec:     1,
			StartRetries:     3,
			AutoRestart:      "unexpected",
			ExitCodes:        []int{1, 2},
			StopSignal:       0,
			StopWaitSecs:     10,
			StdErrLogfile:    "/dev/null",
			StdOutLogfile:    "/dev/null",
			Environments:     []string{"KEY1=VALUE1", "KEY2=VALUE2"},
			WorkingDirectory: ".",
			Umask:            "002",
		}
	data := serializeStructToBytes(&T{expectedResult})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	tasks, _ := cp.Parse()

	expectedLength := 1
	if len(tasks) != expectedLength {
		t.Fatalf(`len(tasks) equal to "%v", should be equal to "%d"`, len(tasks), expectedLength)
	}
	result := tasks[0]
	if result.Name != expectedResult.Name {
		t.Fatalf(`result.Name equal to "%s", should be equal to "%s"'`, result.Name, expectedResult.Name)
	} else if !result.EqualTo(expectedResult) {
		t.Fatalf(`result doesn't equal to expectedResult'`)
	}
}

func TestWithTManyTasks(t *testing.T) {
	type T struct {
		Task           *task.Task `comment:"Task"`
		AnotherTask    *task.Task `comment:"AnotherTask"`
		AnotherOneTask *task.Task `comment:"AnotherOneTask"`
	}
	expectedTasks := []*task.Task{
		{
			Name:    "Task",
			Command: "ls -la",
		},
		{
			Name:    "AnotherTask",
			Command: "cd ../directory/",
		},
		{
			Name:    "AnotherOneTask",
			Command: "rm -rf *",
		},
	}

	data := serializeStructToBytes(&T{
		expectedTasks[0],
		expectedTasks[1],
		expectedTasks[2],
	})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	resultTasks, _ := cp.Parse()

	expectedLength := len(expectedTasks)
	if len(resultTasks) != expectedLength {
		t.Fatalf(`len(resultTasks) equal to "%v", should be equal to "%d"`, len(resultTasks), expectedLength)
	}
	for i, _ := range expectedTasks {
		result := resultTasks[i]
		expectedResult := expectedTasks[i]
		if expectedResult.Name != result.Name {
			t.Fatalf(`result.Name equal to "%s", should be equal to "%s"`, result.Name, expectedResult.Name)
		} else if result.Command != expectedTasks[i].Command {
			t.Fatalf(`result.Command equal to "%s", should be equal to "%s"`, result.Command, expectedResult.Command)
		}
	}
}

func TestWithIncorrectData(t *testing.T) {
	data := []byte("IncorrectData")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "key-value delimiter not found: IncorrectData"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}

func TestWithEmptyTask(t *testing.T) {
	data := []byte("[Task]")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "task must have 'command' value"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}

func TestWithoutCommandValue(t *testing.T) {
	data := []byte("[Task]\n" +
		"auto_restart=unexpected")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "task must have 'command' value"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(
			`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}

func TestWithIncorrectAutoStartValue(t *testing.T) {
	type T struct {
		*task.Task `comment:"Task"`
	}
	taskData :=
		&task.Task{
			Name:        "Task",
			Command:     "ls -la",
			AutoRestart: "qwerty",
		}
	data := serializeStructToBytes(&T{taskData})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "auto_restart option should be 'never', 'unexpected' or 'always' value"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}
