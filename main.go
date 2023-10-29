package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Entry struct {
	Header      string
	Message     []string
	HeaderDone  bool
	MessageDone bool
}

func NewEntry(h string, m []string) *Entry {
	return &Entry{
		Header:      h,
		Message:     m,
		HeaderDone:  false,
		MessageDone: false,
	}
}

func (e *Entry) setHeader(s string) {
	var str string
	arr := splitLine(s)

	re, err := regexp.Compile("\\*\\*\\*")
	if err != nil {
		log.Println(err)
	}

	if arr[0] == "TASK" && re.Match([]byte(arr[len(arr)-1])) {
		for _, i := range arr {
			str += i + " "
		}
		e.Header = str
		e.HeaderDone = true
	} else if arr[0] == "TASK" && !re.Match([]byte(arr[len(arr)-1])) {
		for _, i := range arr {
			str += i + " "
		}
		e.Header = str
	} else if !e.HeaderDone && e.Header != "" {
		for _, i := range arr {
			str += i + " "
		}
		e.Header += str
		if arr[len(arr)-1] == "***" {
			e.HeaderDone = true
		}
	}
}

func (e *Entry) setMessage(s string) {
	arr := splitLine(s)

	states := []string{"ok:", "skipping:", "changed:", "included:"}

	if len(arr) == 1 {
		e.MessageDone = true
	}

	for _, state := range states {
		if arr[0] == state && arr[len(arr)-1] != "{" {
			e.Message = append(e.Message, s)
		} else {
			e.Message = append(e.Message, s)
		}
		return
	}
}

func (e *Entry) clearEntry() {
	e.Header = ""
	e.Message = []string{}
	e.HeaderDone = false
	e.MessageDone = false
}

type Entries []Entry

func (es Entries) FilterChanged() {
	var indexes []int

	for _, e := range es {
		for i, m := range e.Message {
			arr := splitLine(m)
			if arr[0] == "changed:" {
				indexes = append(indexes, i)
			}
		}
		if len(indexes) > 0 {
			fmt.Println(e.Header)
			for _, i := range indexes {
				fmt.Println(e.Message[i])
			}
			fmt.Println()
		}
		indexes = []int{}
	}
}

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

func (app *Config) ParseFile() Entries {
	var entries []Entry
	e := NewEntry("", []string{})
	sc := bufio.NewScanner(app.SourceFile)
	for sc.Scan() {
		if !e.HeaderDone {
			e.setHeader(sc.Text())
			continue
		}
		if !e.MessageDone {
			e.setMessage(sc.Text())
			continue
		}

		entries = append(entries, *e)

		e.clearEntry()
	}

	return entries
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

	entries := app.ParseFile()
	entries.FilterChanged()
}
