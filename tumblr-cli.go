// tumblr-cli
// see Tumblr API docs: https://www.tumblr.com/docs/en/api/v2

package main

import (
	"fmt"
	"github.com/kurrik/oauth1a"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func doApiRequest(method string, url string, values url.Values) ([]byte, error) {
	contents, err := ioutil.ReadFile("CREDENTIALS")

	if err != nil {
		contents, err = ioutil.ReadFile(".tumblr/CREDENTIALS")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	lines := strings.Split(string(contents), "\n")

	service := &oauth1a.Service{
		RequestURL:   "https://www.tumblr.com/oauth/request_token",
		AuthorizeURL: "https://www.tumblr.com/oauth/authorize",
		AccessURL:    "https://www.tumblr.com/oauth/access_token",

		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    lines[0],
			ConsumerSecret: lines[1],
			CallbackURL:    "",
		},

		Signer: new(oauth1a.HmacSha1Signer),
	}

	httpClient := new(http.Client)
	//userConfig := &oauth1a.UserConfig{}
	//userConfig.GetRequestToken(service, httpClient)
	//url, err := userConfig.GetAuthorizeURL(service)

	userConfig := oauth1a.NewAuthorizedConfig(lines[2], lines[3])

	httpRequest, err := http.NewRequest(method, url, strings.NewReader(values.Encode()))

	if err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(1)
	}

	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	service.Sign(httpRequest, userConfig)
	httpResponse, err := httpClient.Do(httpRequest)

	if err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(1)
	}

	defer httpResponse.Body.Close()

	return ioutil.ReadAll(httpResponse.Body)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: please provide a command! Run \"tumblr-cli help\" for getting list of commands")
		os.Exit(1)
	}

	var thisBlog string

	switch os.Args[1] {
	// help
	case "?", "-?", "-h", "--help", "help":
		fmt.Println("tumblr-cli - Tumblr command line interface")
		fmt.Println("command: " + os.Args[1])
		os.Exit(0)

	// create a new post
	case "new", "create":
		// [1] create
		// [2] <BLOG>
		// [3] YYYY-MM-DDTHH:MM:SS
		// [4] <TITLE>
		// [5] <BODY>
		// [6] tags separated by comma

		p := Post{}

		// BLOG
		thisBlog = os.Args[2]

		// TIMESTAMP
		var err error
		p.Time, err = time.Parse("2006-01-02T15:04:05", os.Args[3])

		if err != nil {
			fmt.Println("tumblr-cli - Tumblr command line interface")
			fmt.Println("command: " + os.Args[1])
			fmt.Println("invalid timestamp format")
			os.Exit(1)
		}

		// TITLE
		p.Title = os.Args[4]

		if len(os.Args) > 4 {
			if len(os.Args[5]) > 0 {
				p.Body = os.Args[5]
			}
		}

		if len(os.Args) > 5 {
			if len(os.Args[6]) > 0 {
				p.Tags = strings.Split(os.Args[6], ",")
			}
		}

		fmt.Printf("sending data: %+v\n", p)

		///* debug */ os.Exit(0)

		// API
		apiRequestURL := "https://api.tumblr.com/v2/blog/" + thisBlog + "/post"
		apiValues := p.GetTumblrApiValues()

		httpContents, err := doApiRequest("POST", apiRequestURL, apiValues)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR! can't read http response body | %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(httpContents))
		os.Exit(0)

	// update existing posting:
	// tumblr-cli update <blog> <id> <status> <time>
	case "update":
		// yet not implemented
		fmt.Println("update yet not implemented")
		os.Exit(0)

	// delete a post
	// tumblr-cli delete <blog> <id>
	case "delete":
		httpContents, err := delete(os.Args[2], os.Args[3])

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR! can't read http response body | %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(httpContents))
		os.Exit(0)

	// get list of draft posts
	case "drafts", "posts":
		// BLOG
		thisBlog = os.Args[2]
		
		// API
		requestURL := "https://api.tumblr.com/v2/blog/" + thisBlog + "/posts"

		if os.Args[1] == "drafts" {
			requestURL = requestURL + "/draft"
		}

		httpContents, err := doApiRequest("GET", requestURL, url.Values{})

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR! can't read http response body | %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(httpContents))
		os.Exit(0)

	// get info
	case "info":
		httpContents, err := info(os.Args[2]);

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR! API request failed | %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(httpContents))
		os.Exit(0)

	// get version
	case "version":
		fmt.Println("tumblr-cli verson 0.1.4 (2017-09-05)")
		os.Exit(0)

	// debugging 
	case "debug":
		fmt.Println("debugging ...");
		fmt.Println(os.Args[0]);
		fmt.Println(os.Args[1]);
		fmt.Println(os.Args[2]);
		os.Exit(0)

	// no matching command
	default:
		fmt.Fprintf(os.Stderr, "ERROR! unknown command \""+os.Args[1]+"\"\n")
		os.Exit(1)
	}

}
