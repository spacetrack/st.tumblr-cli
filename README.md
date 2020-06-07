# ST.Tumblr-CLI

very simple, first steps of using the Tumblr API with Go(lang)


## Bookmarks

- Tumblr REST API  
  https://www.tumblr.com/docs/en/api/v2

- Tumblr API with Go https://www.tumblr.com/docs/en/api/v2 | https://github.com/tumblr/tumblrclient.go


## Deployment
### Deployment: deploy on Development (local)

https://medium.com/better-programming/a-complete-go-development-environment-with-docker-and-vscode-a3e4410d27f7
https://code.visualstudio.com/docs/remote/containers

In your local Windows Git-Bash, run following commands:

For ongoing development, in your local Windows Git-Bash, run following commands:

```
# create the container with the source code mounted into the container
$ cd /c/projects/st.tumblr-cli
$ winpty powershell -Command 'docker create -i -t -v ${PWD}:/workspaces/st.tumblr-cli -w /workspaces/st.tumblr-cli -e GOOS=linux -e GOARCH=amd64 --name st.tumblr-cli.dev golang'

# start the container
$ docker start st.tumblr-cli.dev

# exec go version
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'go version'

# exec go get
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'go get -d -v ./...'

# exec go fmt
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'go fmt'

# exec go run
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'go fmt && go run . version'
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'go fmt && go run . posts daily-passcode'
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'go fmt && go run . new daily-passcode 2019-12-24T15:00:00 "hello world b" "hello world t" "tag1,tag2,tag3"'

# exec go build Windows
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'env GOOS=windows GOARCH=amd64 go build .'

# exec go build RasPi
$ winpty docker exec -i -t st.tumblr-cli.dev sh -c 'env GOOS=linux GOARCH=arm GOARM=5 go build .'


```
