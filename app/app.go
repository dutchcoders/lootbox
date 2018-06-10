package app

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"go.dutchsec.com/lootbox/app/twitter"

	"github.com/fatih/color"
	logging "github.com/op/go-logging"

	"github.com/mvdan/xurls"
)

var log = logging.MustGetLogger("app")

type Config struct {
	destinationDir string
	threadCount    int

	consumerKey    string
	consumerSecret string

	token       string
	tokenSecret string

	keywords []string

	urlSubstrings []string
}

type App struct {
	Config
}

type OptionFn func(*App) error

func WithKeywords(val string) (OptionFn, error) {
	return func(b *App) error {
		b.keywords = strings.Split(val, ",")
		return nil
	}, nil
}

func WithToken(val string) (OptionFn, error) {
	return func(b *App) error {
		b.token = val
		return nil
	}, nil
}

func WithTokenSecret(val string) (OptionFn, error) {
	return func(b *App) error {
		b.tokenSecret = val
		return nil
	}, nil

}

func WithConsumerSecret(val string) (OptionFn, error) {
	return func(b *App) error {
		b.consumerSecret = val
		return nil
	}, nil

}

func WithConsumerKey(val string) (OptionFn, error) {
	return func(b *App) error {
		b.consumerKey = val
		return nil
	}, nil

}

func WithDestinationDir(dir string) (OptionFn, error) {
	return func(b *App) error {
		b.destinationDir = dir
		return nil
	}, nil

}

func WithUrlSubstrings(val string) (OptionFn, error) {
	return func(b *App) error {
		b.urlSubstrings = strings.Split(val,"|")
		return nil
	}, nil

}

func New(options ...OptionFn) (*App, error) {
	app := &App{
		Config: Config{
			destinationDir: "./loot/",
			threadCount:    20,
		},
	}

	for _, fn := range options {
		if err := fn(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *App) Run() error {
	d := Downloader(
		context.Background(), 
		WithThreads(a.threadCount),
		WithFileStorage(a.destinationDir), 
		WithFilter(a.urlSubstrings))

	re, err := xurls.StrictMatchingScheme("hxxps?")
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	tc := twitter.New(
		twitter.WithConsumerKey(a.consumerKey),
		twitter.WithConsumerSecret(a.consumerSecret),
		twitter.WithToken(a.token),
		twitter.WithTokenSecret(a.tokenSecret),
	)

	for _, keyword := range a.keywords {
		wg.Add(1)

		go func() {
			defer wg.Done()

			feed := tc.Feed(context.Background(), keyword)
			for tweet := range feed {
				// fmt.Println(color.YellowString(tweet))

				links := re.FindAllString(tweet, -1)
				for _, link := range links {
					re := regexp.MustCompile(`[\[\]]`)

					link = re.ReplaceAllString(link, "")

					u, err := url.Parse(link)
					if err != nil {
						fmt.Println(color.RedString(link, err.Error()))
						continue
					}

					if u.Scheme == "" {
						continue
					}

					if u.Scheme == "hxxp" {
						u.Scheme = "http"
					} else if u.Scheme == "hxxps" {
						u.Scheme = "https"
					}

					d.Download(u)
				}
			}
		}()
	}

	wg.Wait()

	return nil
}
