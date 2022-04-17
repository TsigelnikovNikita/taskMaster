package task_master

import (
	"main/config_parser"
	"main/input_reader"
	"main/output_printer"
	"main/program"
)


type ConfigParser interface {
	Parse() (map[string]*program.Program, error)
}

type TaskMaster struct {
	inputReader input_reader.InputReader
	outputPrinter output_printer.OutputPrinter
	configParser ConfigParser

	programs map[string]*program.Program
}

var taskMaster TaskMaster

func GetTaskMaster() *TaskMaster {
	return &taskMaster
}

func (t *TaskMaster) SetInputReader(inputReader input_reader.InputReader) {
	t.inputReader = inputReader
}

func (t *TaskMaster) SetOutputPrinter(outputPrinter output_printer.OutputPrinter) {
	t.outputPrinter = outputPrinter
}

func (t *TaskMaster) SetConfigParser(configParser config_parser.ConfigParser) {
	t.configParser = configParser
}

func (t *TaskMaster) RunProgram() {
	if programs, err := t.configParser.Parse(); err != nil {
		t.programs = programs
	}
}
