package app

import (
	"io"
	"net/url"

	"github.com/puerkitobio/goquery"
)

type Document struct {
	*goquery.Document

	baseURL *url.URL
}

func NewDocumentFromReader(baseURL *url.URL, r io.Reader) (*Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	return &Document{
		Document: doc,
		baseURL:  baseURL,
	}, nil
}

func (d *Document) ExtractLinks(elem string, attr string) []*url.URL {
	links := []*url.URL{}

	d.Find(elem).Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr(attr); !ok {
		} else if rel, err := url.Parse(val); err != nil {
		} else if u := d.baseURL.ResolveReference(rel); err != nil {
		} else if d.baseURL.Host != u.Host {
		} else {
			// remove query
			u.RawQuery = ""
			u.Fragment = ""

			links = append(links, u)
		}
	})

	return links
}
