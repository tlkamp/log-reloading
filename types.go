package main

import "fmt"

type config struct {
	Server  server
	Logging logging
}

func (c config) ToString() string {
	return fmt.Sprintf("%v", c)
}

type logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Colors bool   `yaml:"colors"`
	Hash   string `yaml:"-"`
}

type server struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}
