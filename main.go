// TODO:
//   - add proper error handling
//	 - use a MultiWriter for the logging, i.e.,
//			mw := io.MultiWriter(os.Stdout, logFile)
//			logrus.SetOutput(mw)

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

var gitlabClient *gitlab.Client
var singleClient *bool

func getClient() *gitlab.Client {
	var err error
	if singleClient == nil {
		gitlabClient, err = gitlab.NewClient(getToken())
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		b := true
		singleClient = &b
	}
	return gitlabClient
}

func getFileContents(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func getToken() string {
	apiToken, isSet := os.LookupEnv("GITLAB_API_PRIVATE_TOKEN")
	if apiToken == "" || !isSet {
		panic("[ERROR] Must set $GITLAB_API_PRIVATE_TOKEN")
	}
	return apiToken
}

func main() {
	filename := flag.String("file", "", "Path to GitLab config file (json or yaml).")
	user := flag.String("user", "", "List everything for the given user.")
	destroy := flag.Bool("destroy", false, "Should destroy all projects listed in the given Gitlab config file.")
	flag.Parse()

	if *user != "" {
		// TODO
		user, err := getUser(*user)
		if err != nil {
			panic(err)
		}
		fmt.Println("user", user.ID)
		getUserProjects(user.ID)
	} else if *filename != "" {
		groups, err := getGroups(*filename)
		if err != nil {
			panic(err)
		}

		processProjects(groups, *destroy)
	}
}
