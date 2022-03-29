// TODO:
//   - add proper error handling
//   - add support accepting `json` or `yaml` from `stdin`?

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

var git *gitlab.Client
var singleClient *bool

func getClient() *gitlab.Client {
	var err error
	if singleClient == nil {
		git, err = gitlab.NewClient(getToken())
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		b := true
		singleClient = &b
	}
	return git
}

func getToken() string {
	apiToken, isSet := os.LookupEnv("GITLAB_API_PRIVATE_TOKEN")
	if apiToken == "" || !isSet {
		panic("[ERROR] Must set $GITLAB_API_PRIVATE_TOKEN")
	}
	return apiToken
}

func parseFile(filename string) ([]Group, error) {
	// Move into a readfile function?
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var groups []Group
	extension := filepath.Ext(filename)
	if extension == ".json" {
		err = json.Unmarshal(content, &groups)
	} else if extension == ".yaml" {
		err = yaml.Unmarshal(content, &groups)
	} else {
		panic("[ERROR] File extension not recognized, must be either `json` or `yaml`.")
	}

	if err != nil {
		panic(err)
	}

	return groups, err
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
		groups, err := parseFile(*filename)
		if err != nil {
			panic(err)
		}

		processProjects(groups, *destroy)
	}
}
