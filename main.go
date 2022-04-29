package main

import (
	"main/config_parser"
	"main/input_reader"
	"main/output_printer"
	"main/task_master"
)

func main() {
	inputReader := input_reader.NewInputReader()
	outputPrinter := output_printer.NewOutputPrinter()
	taskMaster := task_master.GetTaskMaster()
	configParser := config_parser.NewConfigParser("config.ini")

	taskMaster.SetOutputPrinter(outputPrinter)
	taskMaster.SetInputReader(inputReader)
	taskMaster.SetConfigParser(configParser)
	taskMaster.RunProgram()
}
