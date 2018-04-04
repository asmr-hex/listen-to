package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

type LogLine struct {
	Raw     string
	Columns []string
}

func NewLogLine(raw string) *LogLine {
	return &LogLine{
		Raw:     raw,
		Columns: strings.Split(raw, " "),
	}
}

func (l *LogLine) GetMusic() string {
	var (
		music = ""
	)

	// the log lines are formatted s.t. everything after the 2nd column
	// is part of the music the user is recomending.
	for i := 2; i < len(l.Columns); i++ {
		music += " " + l.Columns[i]
	}

	return music
}

func (l *LogLine) GetUser() string {
	return l.Columns[1]
}

func (l *LogLine) PrintRecommendation() {
	fmt.Println(
		fmt.Sprintf(
			"%s says you should listen to %s",
			color.GreenString(l.GetUser()),
			color.YellowString(l.GetMusic()),
		),
	)
}

func FormatLogLine(t time.Time, usr *user.User, music string) string {
	return fmt.Sprintf(
		"%s %s %s\n",
		t.Format(time.RFC3339),
		usr.Username,
		music,
	)
}

func WriteToLogFile(fname, logLine string) error {
	fd, err := os.OpenFile(
		fname,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0777, // we want all user to be able to modify this file.
	)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.WriteString(logLine)
	if err != nil {
		return err
	}

	return nil
}

func ListenForRecommendations(fname string) {
	// start listening for changes to filesystem
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// watch the log file
	watcher.Add(fname)

	// block and listen for file changes
	for {
		select {
		case e := <-watcher.Events:
			if e.Op == fsnotify.Write {
				// read the last line of the file
				f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				scanner := bufio.NewScanner(f)
				scanner.Split(bufio.ScanLines)

				lastLine := &LogLine{}
				for scanner.Scan() {
					lastLine = NewLogLine(scanner.Text())
				}

				lastLine.PrintRecommendation()
			}
		}
	}
}

func ListAllRecommendations(fname string) error {
	// read the last line of the file
	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := NewLogLine(scanner.Text())
		line.PrintRecommendation()
	}

	return nil
}
