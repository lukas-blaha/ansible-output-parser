package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Entry struct {
	Header      string
	Message     string
	HeaderDone  bool
	MessageDone bool
}

func NewEntry(h, m string) *Entry {
	return &Entry{
		Header:      h,
		Message:     m,
		HeaderDone:  false,
		MessageDone: false,
	}
}

func (e *Entry) createEntry(s string) {
	var str string
	arr := splitLine(s)

	// Set header
	if arr[0] == "TASK" && arr[len(arr)-1] == "***" {
		for _, i := range arr {
			str += i + " "
		}
		e.Header = str
		e.HeaderDone = true
	} else if arr[0] == "TASK" && arr[len(arr)-1] != "***" {
		for _, i := range arr {
			str += i + " "
		}
		e.Header = str
	} else if !e.HeaderDone && e.Header != "" {
		for _, i := range arr {
			str += i + " "
		}
		e.Header += str
		for _, i := range arr {
			str += i + " "
		}
		e.Header += str
		if arr[len(arr)-1] == "***" {
			e.HeaderDone = true
		}
	}

	// Set Message
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
		checkHeader(sc.Text())
	}
}

func checkArguments() error {
	if len(os.Args) != 3 {
		return errors.New("You need to specify source and target file.\n")
	}

	return nil
}

func splitLine(s string) []string {
	arr := strings.Split(s, " ")
	return arr
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
