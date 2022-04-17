package config_parser

import (
	"gopkg.in/ini.v1"
	"main/program"
)

type configReader interface {
	getData() ([]byte, error)
}

type ConfigParser interface {
	Parse() (map[string]*program.Program, error)
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

func (i *IniConfigParser) Parse() (map[string]*program.Program, error) {
	if data, err := i.configReader.getData(); err != nil {
		return nil, err
	} else {
		return i.parse(data)
	}
}

func (i *IniConfigParser) parse(data []byte) (map[string]*program.Program, error) {
	cfg, err := ini.Load(data)
	if err != nil {
		return nil, err
	}

	cfg.DeleteSection("DEFAULT") // DEFAULT is section by default. We don't need it
	result := make(map[string]*program.Program)
	for _, section := range cfg.Sections() {
		newProgram := program.NewProgram()
		if err := section.StrictMapTo(&newProgram); err != nil {
			return nil, err
		} else if err = newProgram.IsCorrect(); err != nil {
			return nil, err
		}
		newProgram.Name = section.Name()
		result[newProgram.Name] = newProgram
	}
	return result, nil
}
