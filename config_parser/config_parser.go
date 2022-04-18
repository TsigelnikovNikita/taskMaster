package config_parser

import (
	"gopkg.in/ini.v1"
	"main/task"
)

type configReader interface {
	getData() ([]byte, error)
}

type ConfigParser interface {
	Parse() (map[string]*task.Task, error)
	setConfigReader(cr configReader)
}

type IniConfigParser struct {
	configReader configReader
}

func (i *IniConfigParser) setConfigReader(cr configReader) {
	i.configReader = cr
}

func NewConfigParser(filename string) ConfigParser {
	return &IniConfigParser{fileConfigReader{filename}}
}

func (i *IniConfigParser) Parse() (map[string]*task.Task, error) {
	if data, err := i.configReader.getData(); err != nil {
		return nil, err
	} else {
		return i.parse(data)
	}
}

func (i *IniConfigParser) parse(data []byte) (map[string]*task.Task, error) {
	cfg, err := ini.Load(data)
	if err != nil {
		return nil, err
	}

	cfg.DeleteSection("DEFAULT") // DEFAULT is section by default. We don't need it
	result := make(map[string]*task.Task)
	for _, section := range cfg.Sections() {
		task := task.NewTask()
		if err := section.StrictMapTo(&task); err != nil {
			return nil, err
		} else if err = task.IsCorrect(); err != nil {
			return nil, err
		}
		task.Name = section.Name()
		result[task.Name] = task
	}
	return result, nil
}
