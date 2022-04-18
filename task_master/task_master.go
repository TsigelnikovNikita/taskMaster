package task_master

import (
	"fmt"
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

func (t *TaskMaster) ReloadConfig() error {
	newPrograms, err := t.configParser.Parse()
	if err != nil {
		return err
	}

	prevPrograms := t.programs
	t.programs = make(map[string]*program.Program)
	for _, p := range newPrograms {
		t.programs[p.Name] = p
		if prevProgram, ok := prevPrograms[p.Name]; !ok {
		//	t.startProgram(prevProgram)
		} else if !p.EqualTo(prevProgram) {
		//  t.reloadProgram(prevProgram)
		}
	}
	return nil
}

func (t *TaskMaster) RunProgram() {
	err := t.ReloadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
