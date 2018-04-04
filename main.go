package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

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

		formattedLogLine := FormatLogLine(t, usr, music)

		err = WriteToLogFile(c.String("file"), formattedLogLine)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "subscribe",
			Aliases: []string{"s"},
			Usage:   "start a daemon to listen for updates",
			Flags:   app.Flags,
			Action: func(c *cli.Context) error {
				// block here for listening
				ListenForRecommendations(c.String("file"))

				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all recommended music",
			Flags:   app.Flags,
			Action: func(c *cli.Context) error {
				err := ListAllRecommendations(c.String("file"))
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
