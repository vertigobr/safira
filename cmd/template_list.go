// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/gookit/color.v1"
)

type TemplateInfo struct {
	TemplateName string `json:"template"`
	Platform     string `json:"platform"`
	Language     string `json:"language"`
	Source       string `json:"source"`
	Description  string `json:"description"`
	Repository   string `json:"repo"`
	Official     string `json:"official"`
}

var templateListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lists the official Vertigo iPaaS templates",
	Long:    "Lists the official Vertigo iPaaS templates",
	Example: `To list all the official templates, run:

    $ safira template list`,
	RunE:                       runTemplateList,
	SuggestionsMinimumDistance: 1,
}

func init() {
	templateCmd.AddCommand(templateListCmd)
}

func runTemplateList(cmd *cobra.Command, _ []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	templateInfo, err := getTemplateInfo(verboseFlag)
	if err != nil {
		return err
	}

	outputTemplateInfo(templateInfo, verboseFlag)

	return nil
}

func getTemplateInfo(verboseFlag bool) ([]TemplateInfo, error) {
	if verboseFlag {
		fmt.Printf("%s Fetching template information\n", color.Blue.Text("[v]"))
	}

	req, reqErr := http.NewRequest(http.MethodGet, faasTemplateStoreURL, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("%s Error in requesting template information", color.Red.Text("[!]"))
	}

	reqContext, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	req = req.WithContext(reqContext)
	client := http.DefaultClient
	res, clientErr := client.Do(req)
	if clientErr != nil || res.Body == nil {
		return nil, fmt.Errorf("%s Error getting the list of templates", color.Red.Text("[!]"))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s Unexpected request status %d", color.Red.Text("[!]"), res.StatusCode)
	}

	body, bodyErr := ioutil.ReadAll(res.Body)
	if bodyErr != nil {
		return nil, fmt.Errorf("%s Error when getting the request body", color.Red.Text("[!]"))
	}

	var templatesInfo []TemplateInfo
	if err := json.Unmarshal(body, &templatesInfo); err != nil {
		return nil, fmt.Errorf("%s Error when unmarshalling template information", color.Red.Text("[!]"))
	}

	return templatesInfo, nil
}

func outputTemplateInfo(templates []TemplateInfo, verboseFlag bool) {
	if len(templates) == 0 {
		fmt.Printf("%s No templates\n", color.Green.Text("[+]"))
		return
	}

	if verboseFlag {
		fmt.Printf("%s Processing template information\n", color.Blue.Text("[v]"))
	}

	var buff bytes.Buffer
	lineWriter := tabwriter.NewWriter(&buff, 0, 0, 3, ' ', 0)

	_, _ = fmt.Fprintln(lineWriter)
	_, _ = fmt.Fprintf(lineWriter, "NAME\tLANGUAGE\tPLATFORM\tSOURCE\tDESCRIPTION\n")
	for _, template := range templates {
		_, _ = fmt.Fprintf(lineWriter, "%s\t%s\t%s\t%s\t%s\n",
			template.TemplateName,
			template.Language,
			template.Platform,
			template.Source,
			template.Description)
	}

	_, _ = fmt.Fprint(lineWriter)

	_ = lineWriter.Flush()

	fmt.Println(buff.String())
}
