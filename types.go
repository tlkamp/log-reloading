package main

import "fmt"

type config struct {
	Logging logging
}

func (c config) ToString() string {
	return fmt.Sprintf("%v", c)
}

type logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Colors bool   `yaml:"colors"`
	Hash   string `yaml:"-"` // This field is never read or displayed except in debug messages
}
