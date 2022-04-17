package config_parser

import (
	"gopkg.in/ini.v1"
	"log"
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
		log.Fatal(err)
	}

	cfg.DeleteSection("DEFAULT") // DEFAULT is section by default. We don't need it
	var result map[string]*program.Program
	for _, section := range cfg.Sections() {
		newProgram := program.NewProgram()
		if err := section.MapTo(&newProgram); err != nil {
			return nil, err
		}
		result[section.Name()] = newProgram
	}
	return result, nil
}
