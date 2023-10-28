package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

type Entry struct {
	Header  string
	Message string
}

func NewEntry(h, m string) *Entry {
	return &Entry{
		Header:  h,
		Message: m,
	}
}

type Entries []Entry

type Config struct {
	SourcePath string
	TargetPath string
	SourceFile *os.File
}

func NewConfig(s, t string) *Config {
	return &Config{
		SourcePath: s,
		TargetPath: t,
	}
}

func (app *Config) LoadSource() error {
	f, err := os.Open(app.SourcePath)
	if err != nil {
		return err
	}

	app.SourceFile = f
	return nil
}

func (app *Config) ParseFile() {
	sc := bufio.NewScanner(app.SourceFile)
	for sc.Scan() {
	}
}

func checkArguments() error {
	if len(os.Args) != 3 {
		return errors.New("You need to specify source and target file.\n")
	}

	return nil
}

func main() {
	if err := checkArguments(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	app := NewConfig(os.Args[1], os.Args[2])

	err := app.LoadSource()
	if err != nil {
		log.Panic(err)
	}

	app.ParseFile()
}
