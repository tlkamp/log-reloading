package main

import "fmt"

type config struct {
	Logging logging
	Hash    string `yaml:"-"` // Never display the hash to the user
}

func (c config) ToString() string {
	return fmt.Sprintf("%v", c)
}

type logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Colors bool   `yaml:"colors"`
}
