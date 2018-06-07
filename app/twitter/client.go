package twitter

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/fatih/color"
)

func WithConsumerKey(s string) func(*config) {
	return func(t *config) {
		t.consumerKey = s
	}
}

func WithConsumerSecret(s string) func(*config) {
	return func(t *config) {
		t.consumerSecret = s
	}
}

func WithToken(s string) func(*config) {
	return func(t *config) {
		t.token = s
	}
}

func WithTokenSecret(s string) func(*config) {
	return func(t *config) {
		t.tokenSecret = s
	}
}

func New(opts ...func(*config)) *Client {
	c := &config{}

	for _, opt := range opts {
		opt(c)
	}

	config := oauth1.NewConfig(c.consumerKey, c.consumerSecret)

	httpClient := config.Client(
		oauth1.NoContext,
		oauth1.NewToken(c.token, c.tokenSecret),
	)

	client := twitter.NewClient(httpClient)

	return &Client{
		client,
	}
}

type Client struct {
	*twitter.Client
}

type config struct {
	consumerKey    string
	consumerSecret string
	token          string
	tokenSecret    string
}

func (t *Client) Feed(ctx context.Context, query string) chan string {
	feed := make(chan string)

	go func() {
		defer close(feed)

		for {

			search, _, err := t.Search.Tweets(&twitter.SearchTweetParams{
				Query: query,
				Count: 100,
			})

			if err != nil {
				fmt.Println(color.RedString(err.Error()))
				time.Sleep(time.Second * 60 * 10)
				continue
			}

			for _, tweet := range search.Statuses {
				select {
				case feed <- tweet.Text:
				case <-ctx.Done():
					return
				}
			}

			time.Sleep(time.Second * 60 * 1)
		}
	}()

	return feed
}
