package config_parser

import (
	"bytes"
	"gopkg.in/ini.v1"
	"main/program"
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
		*program.Program `comment:"Program"`
	}
	expectedResult :=
		&program.Program{
			Name:             "Program",
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
	programs, _ := cp.Parse()

	expectedLength := 1
	if len(programs) != expectedLength {
		t.Fatalf(`len(programs) equal to "%v", should be equal to "%d"`, len(programs), expectedLength)
	}
	if result, ok := programs[expectedResult.Name]; !ok {
		t.Fatalf(`result.Name equal to "%s", should be equal to "%s"'`, result.Name, expectedResult.Name)
	} else if !result.EqualTo(expectedResult) {
		t.Fatalf(`result doesn't equal to expectedResult'`)
	}
}

func TestWithTManyPrograms(t *testing.T) {
	type T struct {
		Program           *program.Program `comment:"Program"`
		AnotherProgram    *program.Program `comment:"AnotherProgram"`
		AnotherOneProgram *program.Program `comment:"AnotherOneProgram"`
	}
	expectedPrograms := []*program.Program{
		{
			Name:    "Program",
			Command: "ls -la",
		},
		{
			Name:    "AnotherProgram",
			Command: "cd ../directory/",
		},
		{
			Name:    "AnotherOneProgram",
			Command: "rm -rf *",
		},
	}

	data := serializeStructToBytes(&T{
		expectedPrograms[0],
		expectedPrograms[1],
		expectedPrograms[2],
	})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	resultPrograms, _ := cp.Parse()

	expectedLength := len(expectedPrograms)
	if len(resultPrograms) != expectedLength {
		t.Fatalf(`len(resultPrograms) equal to "%v", should be equal to "%d"`, len(resultPrograms), expectedLength)
	}
	for _, p := range expectedPrograms {
		if result, ok := resultPrograms[p.Name]; !ok {
			t.Fatalf(`result.Name equal to "%s", should be equal to "%s"`, result.Name, p.Name)
		} else if result.Command != p.Command {
			t.Fatalf(`result.Command equal to "%s", should be equal to "%s"`, result.Command, p.Command)
		}
	}
}

func TestWithIncorrectData(t *testing.T) {
	data := []byte("IncorrectData")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "program must have 'command' value"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}

func TestWithEmptyProgram(t *testing.T) {
	data := []byte("[Program]")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "program must have 'command' value"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}

func TestWithoutCommandValue(t *testing.T) {
	data := []byte("[Program]\n" +
		"auto_restart=unexpected")

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "program must have 'command' value"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(
			`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}

func TestWithIncorrectAutoStartValue(t *testing.T) {
	type T struct {
		*program.Program `comment:"Program"`
	}
	programData :=
		&program.Program{
			Name:        "Program",
			Command:     "ls -la",
			AutoRestart: "qwerty",
		}
	data := serializeStructToBytes(&T{programData})

	cr := mockConfigReader{data: data}
	cp := IniConfigParser{configReader: cr}
	expectedResult := "auto_restart option should be 'never', 'unexpected' or 'always' value"
	if _, err := cp.Parse(); err == nil {
		t.Fatalf(`expected err from cp.Parse() method`)
	} else if err.Error() != expectedResult {
		t.Fatalf(`err.Error() equal to "%s", should be equal to "%s"`, err.Error(), expectedResult)
	}
}
