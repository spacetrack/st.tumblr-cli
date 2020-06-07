// tumblr-cli
// see Tumblr API docs: https://www.tumblr.com/docs/en/api/v2

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tumblr/tumblrclient.go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: please provide a command! Run \"tumblr-cli help\" for getting list of commands")
		os.Exit(1)
	}

	// ----------------------------------------
	// read credentials
	//

	fileContents, err := ioutil.ReadFile("CREDENTIALS")

	if err != nil {
		fileContents, err = ioutil.ReadFile(".tumblr/CREDENTIALS")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	lines := strings.Split(string(fileContents), "\n")

	// ----------------------------------------
	// init the tumblr client
	//

	client := tumblrclient.NewClientWithToken(strings.TrimSpace(lines[0]), strings.TrimSpace(lines[1]), strings.TrimSpace(lines[2]), strings.TrimSpace(lines[3]))

	var thisBlog string

	switch os.Args[1] {
	// help
	case "?", "-?", "-h", "--help", "help":
		fmt.Println("tumblr-cli - Tumblr command line interface")
		fmt.Println("command: " + os.Args[1])
		os.Exit(0)

	// create a new post
	case "new", "create":
		// [1] new | create
		// [2] <BLOG_IDENTIFIER>
		// [3] <YYYY-MM-DDTHH:MM:SS>
		// [4] <BODY>
		// [5] [<TITLE>]
		// [6] [<TAG>[,<TAG>,[<TAG>]]]

		p := Post{}

		// BLOG
		thisBlog = os.Args[2]

		// TIMESTAMP
		var err error
		p.Time, err = time.Parse("2006-01-02T15:04:05", os.Args[3])

		if err != nil {
			panic(err)
		}

		// BODY
		p.Body = os.Args[4]

		// TITLE
		if len(os.Args) > 5 {
			if len(os.Args[5]) > 0 {
				p.Title = os.Args[5]
			}
		}

		// TAGS
		if len(os.Args) > 6 {
			if len(os.Args[6]) > 0 {
				p.Tags = strings.Split(os.Args[6], ",")
			}
		}

		// the new API ".../posts" does not allow
		// to set the date of the posting, so let's
		// use the legacy API ".../post"
		// (and ... we ditch the p variable ;-))

		r, err := client.GetHttpClient().Post("https://api.tumblr.com/v2/blog/"+thisBlog+"/post", "application/json", strings.NewReader(`{
			"type": "text",
			"date": "`+os.Args[3]+`",
			"title": "`+os.Args[5]+`",
			"body": "`+os.Args[4]+`",
			"tags": "`+os.Args[6]+`"
		}`))

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		if err != nil {
			fmt.Fprintf(os.Stderr, buf.String())
			panic(err)
		}

		fmt.Println(buf.String())

	// update existing posting:
	// tumblr-cli update <blog> <id> <status> <time>
	case "update":
		// yet not implemented
		fmt.Println("update yet not implemented")
		os.Exit(0)

	// delete a post
	// tumblr-cli delete <blog> <id>
	case "delete":
		r, err := client.PostWithParams("blog/"+os.Args[2]+"/post/delete", url.Values{"id": []string{os.Args[3]}})

		if err != nil {
			panic(err)
		}

		fmt.Println(string(r.GetBody()))

	// get list of draft posts
	case "posts":
		// [1] posts
		// [2] <BLOG>

		// BLOG
		fmt.Println("drafts and posts")
		thisBlog = os.Args[2]

		r, err := client.Get("blog/" + os.Args[2] + "/posts")

		if err != nil {
			panic(err)
		}

		fmt.Println(string(r.GetBody()))
		os.Exit(0)

	// get info
	case "info":
		r, err := client.Get("blog/" + os.Args[2] + "/info")

		if err != nil {
			panic(err)
		}

		fmt.Println(string(r.GetBody()))

	// get user info
	case "user-info":
		r, err := client.Get("user/info")

		if err != nil {
			panic(err)
		}

		fmt.Println(string(r.GetBody()))

	// get version
	case "version":
		fmt.Println("st-tumblr-cli verson 0.2.1 (2020-05-09)")
		os.Exit(0)

	// debugging
	case "debug":
		fmt.Println("debugging ...")
		fmt.Println(os.Args[0])
		fmt.Println(os.Args[1])
		fmt.Println(os.Args[2])
		os.Exit(0)

	// no matching command
	default:
		fmt.Fprintf(os.Stderr, "ERROR! unknown command \""+os.Args[1]+"\"\n")
		os.Exit(1)
	}

}
