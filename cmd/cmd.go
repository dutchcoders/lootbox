package cmd

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/minio/cli"
	logging "github.com/op/go-logging"

	app "go.dutchsec.com/lootbox/app"
)

var log = logging.MustGetLogger("cmd")

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rand.Seed(time.Now().UTC().UnixNano())
}

func run(c *cli.Context) error {
	log.Info("Lootbox started.")
	log.Info("Version: %v (%s)", app.Version, app.ShortCommitID)
	log.Info("---------------------------")

	var options []app.OptionFn

	if d := c.String("destination-dir"); d == "" {
		ec := cli.NewExitError("destination-dir not set", 1)
		return ec
	} else if fn, err := app.WithDestinationDir(d); err != nil {
		ec := cli.NewExitError(err.Error(), 1)
		return ec
	} else {
		options = append(options, fn)
	}

	if d := c.String("token"); d == "" {
		ec := cli.NewExitError("token not set", 1)
		return ec
	} else if fn, err := app.WithToken(d); err != nil {
		ec := cli.NewExitError(err.Error(), 1)
		return ec
	} else {
		options = append(options, fn)
	}

	if d := c.String("token-secret"); d == "" {
		ec := cli.NewExitError("token-secret not set", 1)
		return ec
	} else if fn, err := app.WithTokenSecret(d); err != nil {
		ec := cli.NewExitError(err.Error(), 1)
		return ec
	} else {
		options = append(options, fn)
	}

	if d := c.String("consumer-key"); d == "" {
		ec := cli.NewExitError("consumer-key not set", 1)
		return ec
	} else if fn, err := app.WithConsumerKey(d); err != nil {
		ec := cli.NewExitError(err.Error(), 1)
		return ec
	} else {
		options = append(options, fn)
	}

	if d := c.String("consumer-secret"); d == "" {
		ec := cli.NewExitError("consumer-secret not set", 1)
		return ec
	} else if fn, err := app.WithConsumerSecret(d); err != nil {
		ec := cli.NewExitError(err.Error(), 1)
		return ec
	} else {
		options = append(options, fn)
	}

	if d := c.String("keywords"); d == "" {
	} else if fn, err := app.WithKeywords(d); err != nil {
		ec := cli.NewExitError(err.Error(), 1)
		return ec
	} else {
		options = append(options, fn)
	}

	a, err := app.New(options...)
	if err != nil {
		ec := cli.NewExitError(err.Error(), 1)
		return ec
	}

	return a.Run()
}

func New() *cli.App {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer,
			`Version: %s
Release-Tag: %s
Commit-ID: %s
`, color.YellowString(app.Version), color.YellowString(app.ReleaseTag), color.YellowString(app.CommitID))
	}

	app := cli.NewApp()
	app.Name = "lootbox"
	app.Author = ""
	app.Usage = "lootbox"
	app.Flags = append(app.Flags, []cli.Flag{
		cli.StringFlag{Name: "destination-dir", EnvVar: "DESTINATION_DIR", Value: "./loot/", Usage: "Where to download the loot to"},
		cli.StringFlag{Name: "keywords", EnvVar: "KEYWORDS", Value: "hxxp://,#opendir", Usage: "Twitter token"},
		cli.StringFlag{Name: "consumer-key", EnvVar: "TWITTER_CONSUMER_KEY", Value: "", Usage: "Twitter consumer key"},
		cli.StringFlag{Name: "consumer-secret", EnvVar: "TWITTER_CONSUMER_SECRET", Value: "", Usage: "Twitter consumer secret"},
		cli.StringFlag{Name: "token", EnvVar: "TWITTER_TOKEN", Value: "", Usage: "Twitter token"},
		cli.StringFlag{Name: "token-secret", EnvVar: "TWITTER_TOKEN_SECRET", Value: "", Usage: "Twitter token secret"},
	}...)
	app.Description = `lootbox: twitter loot downloader`
	app.Commands = []cli.Command{}
	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = run

	return app
}
