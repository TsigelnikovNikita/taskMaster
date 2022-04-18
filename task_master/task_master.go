package task_master

import (
	"fmt"
	"main/config_parser"
	"main/input_reader"
	"main/output_printer"
	"main/task"
)


type ConfigParser interface {
	Parse() (map[string]*task.Task, error)
}

type TaskMaster struct {
	inputReader input_reader.InputReader
	outputPrinter output_printer.OutputPrinter
	configParser ConfigParser

	tasks map[string]*task.Task
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
	newTasks, err := t.configParser.Parse()
	if err != nil {
		return err
	}

	prevTasks := t.tasks
	t.tasks = make(map[string]*task.Task)
	for _, task := range newTasks {
		t.tasks[task.Name] = task
		if prevTask, ok := prevTasks[task.Name]; !ok {
		//	t.startTask(prevTask)
		} else if !task.EqualTo(prevTask) {
		//  t.reloadTask(prevTask)
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
