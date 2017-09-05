// API blog info
// https://www.tumblr.com/docs/en/api/v2#blog-info

package main

import (
	"net/url"
)

func info(blog string) ([]byte, error) {
	requestURL := "https://api.tumblr.com/v2/blog/" + blog + "/info"
	httpContents, err := doApiRequest("GET", requestURL, url.Values{})

	return httpContents, err
}
