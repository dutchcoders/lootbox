package app

import (
	"bytes"
	"context"
	"io"
	"time"
	"fmt"
	"strings"

	"go.dutchsec.com/lootbox/app/queue"

	"github.com/fatih/color"

	"net/http"
	"net/url"
)

func WithFileStorage(dst string) func(*downloader) {
	return func(d *downloader) {
		d.s = FileStorage(dst)
	}
}

func WithThreads(count int) func(*downloader) {
	return func(d *downloader) {
		d.numThreads = count
	}
}

func WithFilter(arr []string) func(*downloader) {
	return func(d *downloader) {
		d.urlSubstrings = arr
	}
}

func UrlContains(url *url.URL , needles []string) bool{
	if len(needles) == 0 {
		return true
	}
	var u = strings.ToLower(url.String())
	for _, n := range needles {
		if strings.Contains(u, n){
			fmt.Println(color.BlueString("%s contains %s", u,n))
			return true
		}
	}
	return false
}

func Downloader(ctx context.Context, options ...func(*downloader)) *downloader {
	d := &downloader{
		q:          make(chan *url.URL),
		s:          DummyStorage(),
		numThreads: 20,
	}

	for _, fn := range options {
		fn(d)
	}

	go func() {
		for {
			u := queue.Pop()
			if u == nil {
				time.Sleep(10 * time.Second)
				continue
			}

			log.Infof("Downloading: %s\n", u.String())

			d.q <- u
		}

	}()

	for i := 0; i < d.numThreads; i++ {
		go d.run()
	}

	return d
}

type downloader struct {
	q chan *url.URL
	s Storage

	numThreads int

	urlSubstrings []string
}

func (d *downloader) Download(u *url.URL) {
	queue.Push(u)
}

func (d *downloader) run() {
	client := http.DefaultClient

	for {
		u := <-d.q

		resp, err := client.Get(u.String())
		if err != nil {
			log.Error(color.RedString("Error downloading url=%s, error=%s", u.String(), err.Error()))
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		} else {
			log.Error(color.RedString("Error downloading url=%s, statuscode=%d", u.String(), resp.StatusCode))
			continue
		}

		defer resp.Body.Close()

		buff := &bytes.Buffer{}

		body := io.TeeReader(resp.Body, buff)

		f, err := d.s.Create(u)
		if err != nil {
			log.Error(color.RedString(u.String(), err.Error()))
			continue
		}

		if _, err := io.Copy(f, body); err != nil {
			log.Error(color.RedString(u.String(), err.Error()))
			continue
		}

		f.Close()

		doc, err := NewDocumentFromReader(u, buff)
		if err != nil {
			log.Error(color.RedString(u.String(), err.Error()))
			continue
		}

		links := doc.ExtractLinks("a", "href")

		for _, link := range links {
			if UrlContains(link, d.urlSubstrings){
				d.Download(link)
			}
		}
	}
}
