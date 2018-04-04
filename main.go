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
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "listen-to"
	app.Usage = "<3 recommend an artist, song, or album to fellow server pals <3"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "file, f",
			Value:  "/etc/listen-to/music.log",
			Usage:  "log recommendations here",
			EnvVar: "LISTEN_TO_LOG_FILE",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "subscribe",
			Aliases: []string{"s"},
			Usage:   "start a daemon to listen for updates",
			Flags:   app.Flags,
			Action: func(c *cli.Context) error {
				// start listening for changes to filesystem
				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					log.Fatal(err)
				}

				// watch the log file
				watcher.Add(c.String("file"))

				// block and listen for file changes
				for {
					select {
					case e := <-watcher.Events:
						if e.Op == fsnotify.Write {
							// read the last line of the file
							f, err := os.OpenFile(c.String("file"), os.O_RDONLY, 0666)
							if err != nil {
								return err
							}
							defer f.Close()

							scanner := bufio.NewScanner(f)
							scanner.Split(bufio.ScanLines)

							lastLine := ""
							for scanner.Scan() {
								lastLine = scanner.Text()
							}

							// now we have the last line

							data := strings.Split(lastLine, " ")

							music := ""
							for i := 2; i < len(data); i++ {
								music += " " + data[i]
							}

							fmt.Println(
								fmt.Sprintf(
									"%s says you should listen to %s",
									color.GreenString(data[1]),
									color.YellowString(music),
								),
							)

						}
					}
				}
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all recommended music",
			Flags:   app.Flags,
			Action: func(c *cli.Context) error {
				// read the last line of the file
				f, err := os.OpenFile(c.String("file"), os.O_RDONLY, 0666)
				if err != nil {
					return err
				}
				defer f.Close()

				scanner := bufio.NewScanner(f)
				scanner.Split(bufio.ScanLines)

				for scanner.Scan() {
					line := scanner.Text()
					data := strings.Split(line, " ")

					music := ""
					for i := 2; i < len(data); i++ {
						music += " " + data[i]
					}

					fmt.Println(
						fmt.Sprintf(
							"%s says you should listen to %s",
							color.GreenString(data[1]),
							color.YellowString(music),
						),
					)

				}

				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Println("wait, you gotta recommend something ;__;")
		}

		// get user who made the recommendation
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		// get the music they recommended
		music := ""
		for i := 0; i < c.NArg(); i++ {
			music += " " + c.Args().Get(i)
		}

		// get the time they recommended it
		t := time.Now()

		fd, err := os.OpenFile(c.String("file"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()

		_, err = fd.WriteString(
			fmt.Sprintf(
				"%s %s %s\n",
				t.Format(time.RFC3339),
				usr.Username,
				music,
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
