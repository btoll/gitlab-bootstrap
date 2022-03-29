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

type Group struct {
	Name     string    `json:"name,omitempty" yaml:"name,omitempty"`
	Projects []Project `json:"projects,omitempty" yaml:"projects,omitempty"`
}

type Project struct {
	Name       string      `json:"name,omitempty" yaml:"name,omitempty" omitempty:"name"`
	TplName    string      `json:"tpl_name,omitempty" yaml:"tpl_name,omitempty" omitempty:"tpl_name"`
	Visibility string      `json:"visibility,omitempty" yaml:"visibility,omitempty" omitempty:"visibility"`
	Invites    []Invite    `json:"invites,omitempty" yaml:"invites,omitempty" omitempty:"invites"`
	Issues     []IssueType `json:"issues,omitempty" yaml:"issues,omitempty" omitempty:"issues"`
	Releases   []Release   `json:"releases,omitempty" yaml:"releases" omitempty:"releases"`
}

type Release struct {
	Name    string `json:"name,omitempty" yaml:"name" omitempty:"name"`
	Ref     string `json:"ref,omitempty" yaml:"ref" omitempty:"ref"`
	TagName string `json:"tag_name,omitempty" yaml:"tag_name" omitempty:"tag_name"`
}

type Invite struct {
	AccessLevel string `json:"access_level,omitempty" yaml:"access_level,omitempty"`
	Email       string `json:"email,omitempty" yaml:"email,omitempty"`
}

type IssueType struct {
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	Type  string `json:"type,omitempty" yaml:"type,omitempty"`
}

const (
	Incident string = "incident"
	Issue           = "issue"
	TestCase        = "test_case"
)

func getIssueType(issueType string) *string {
	var s string
	switch issueType {
	case Incident:
		s = "incident"
	case TestCase:
		s = "test_case"
	default:
		s = "issue"
	}
	return gitlab.String(s)
}

// https://docs.gitlab.com/ee/development/permissions.html#members
const (
	None       gitlab.AccessLevelValue = 0
	Minimal                            = 5
	Guest                              = 10
	Reporter                           = 20
	Developer                          = 30
	Maintainer                         = 40
	Owner                              = 50
)

func getAccessLevel(accessLevel string) *gitlab.AccessLevelValue {
	var v gitlab.AccessLevelValue
	switch accessLevel {
	case "None":
		v = None
	case "Minimal":
		v = Minimal
	case "Guest":
		v = Guest
	case "Reporter":
		v = Reporter
	case "Maintainer":
		v = Maintainer
	case "Owner":
		v = Owner
	default:
		v = Developer
	}

	return gitlab.AccessLevel(v)
}

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
