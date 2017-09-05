// API post delete
// https://www.tumblr.com/docs/en/api/v2#deleting-posts

package main

import (
	"net/url"
)

func delete(blog string, id string) ([]byte, error) {
	requestURL := "https://api.tumblr.com/v2/blog/" + blog + "/post/delete"

	values := url.Values{}
	values.Set("id", id)

	httpContents, err := doApiRequest("POST", requestURL, values)

	return httpContents, err
}
