package config_parser

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	t.Parallel()

	type T struct {
		*task.Task `comment:"Task"`
	}
	expectedResult := &task.Task{
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
	expectedLength := 1
	data := serializeStructToBytes(&T{expectedResult})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	tasks, _ := cp.Parse()
	result := tasks[0]

	assert.Equal(t, expectedLength, len(tasks))
	assert.Equal(t, expectedResult, result)
}

func TestDefaultValuesShouldBeCorrect(t *testing.T) {
	t.Parallel()

	type T struct {
		*task.Task `comment:"Task"`
	}
	expectedResult := &task.Task{
		Command: "ls -la",
	}
	data := serializeStructToBytes(&T{expectedResult})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	tasks, _ := cp.Parse()
	result := tasks[0]

	assert.Equal(t, "unexpected", result.AutoRestart)
	assert.Equal(t, []int{0}, result.ExitCodes)
}

func TestWithManyTasks(t *testing.T) {
	t.Parallel()

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
	expectedLength := len(expectedTasks)

	data := serializeStructToBytes(&T{
		expectedTasks[0],
		expectedTasks[1],
		expectedTasks[2],
	})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	resultTasks, _ := cp.Parse()

	assert.Equal(t, expectedLength, len(resultTasks))
	for i, expectedTask := range expectedTasks {
		t.Run(fmt.Sprintf("Result #%d should be equal to ExpectedResult #%d", i, i), func(t *testing.T) {
			assert.Equal(t, expectedTask.Name, resultTasks[i].Name)
			assert.Equal(t, expectedTask.Command, resultTasks[i].Command)
		})
	}
}

func TestWithIncorrectData(t *testing.T) {
	t.Parallel()

	data := []byte("IncorrectData")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedError := "key-value delimiter not found: IncorrectData"

	_, err := cp.Parse()

	assert.EqualError(t, err, expectedError)
}

func TestWithEmptyTask(t *testing.T) {
	t.Parallel()

	data := []byte("[Task]")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedError := "task must have 'command' value"

	_, err := cp.Parse()

	assert.EqualError(t, err, expectedError)
}

func TestWithoutCommandValue(t *testing.T) {
	t.Parallel()

	data := []byte("[Task]\n" +
		"auto_restart=unexpected")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedError := "task must have 'command' value"

	_, err := cp.Parse()

	assert.EqualError(t, err, expectedError)
}

func TestWithIncorrectAutoStartValue(t *testing.T) {
	t.Parallel()

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
	expectedError := "auto_restart option should be 'never', 'unexpected' or 'always' value"

	_, err := cp.Parse()

	assert.EqualError(t, err, expectedError)
}
