package config_parser

import "os"

type fileConfigReader struct {
	fileName string
}

func (f fileConfigReader) getData() ([]byte, error) {
	if data, err := os.ReadFile(f.fileName); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
