// TODO:
//   - add proper error handling
//	 - use a MultiWriter for the logging, i.e.,
//			mw := io.MultiWriter(os.Stdout, logFile)
//			logrus.SetOutput(mw)

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getConfigs(filename string) ([]Group, error) {
	content, err := getFileContents(filename)
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
		err = errors.New("[ERROR] File extension not recognized, must be either `json` or `yaml`.")
	}

	return groups, err
}

func getFileContents(filename string) ([]byte, error) {
	f, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}
	return ioutil.ReadFile(f)
}

func main() {
	filename := flag.String("file", "", "Path to GitLab config file (json or yaml).")
	user := flag.String("user", "", "List everything for the given user.")
	destroy := flag.Bool("destroy", false, "Should destroy all projects listed in the given Gitlab config file.")
	flag.Parse()

	if *user != "" {
		// TODO
		//		user, err := getUser(*user)
		//		if err != nil {
		//			panic(err)
		//		}
		//		fmt.Println("user", user.ID)
		//		getUserProjects(user.ID)
	} else if *filename != "" {
		configs, err := getConfigs(*filename)
		if err != nil {
			panic(err)
		}
		p := NewProvisioner(configs)
		p.ProcessConfigs(*destroy)
	}
}
