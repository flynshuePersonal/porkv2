/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	nap "github.com/flynshuePersonal/napv2"
	"github.com/spf13/cobra"
)

type DocsResponse struct {
	Content string `json:"content"`
}

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("must supply repository")
		}
		if err := GetReadme(args[0]); err != nil {
			log.Fatalln("error getting readme for repository:", err)
		}
	},
}

func GetReadme(repo string) error {
	values := strings.Split(repo, "/")
	if len(values) != 2 {
		return fmt.Errorf("must supply repository in owner/project format")
	}
	params := map[string]string{"owner": values[0], "project": values[1]}
	return GithubAPI().Call("docs", params, nil)
}

func DocsSuccess(resp *http.Response) error {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := &DocsResponse{}
	if err := json.Unmarshal(b, response); err != nil {
		return err
	}
	docs, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		return err
	}
	fmt.Println(string(docs))
	return nil
}

func DocsResource() nap.RestResource {
	router := nap.NewRouter()
	router.RegisterFunc(200, DocsSuccess)
	docs := nap.NewResource("GET", "/repos/{{.owner}}/{{.project}}/readme", router)
	return docs
}

func init() {
	rootCmd.AddCommand(docsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
