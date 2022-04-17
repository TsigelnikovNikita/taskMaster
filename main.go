package main

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

type T struct {
	AwesomeSection AwesomeSection `ini:"AwesomeSection"`
}

type AwesomeSection struct {
	StringValue string
	IntValue int
}

func main() {
	data, err := os.ReadFile("config.ini")
	cfg, err := ini.Load(data)
	if err != nil {
		log.Fatal(err)
	}

	t := T{}
	err = cfg.MapTo(&t)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(t.AwesomeSection)
}
