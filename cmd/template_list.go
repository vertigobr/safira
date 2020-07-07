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
	RunE:    runTemplateList,
	SuggestionsMinimumDistance: 1,
}

func init() {
	templateCmd.AddCommand(templateListCmd)
}

func runTemplateList(cmd *cobra.Command, args []string) error {
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
		fmt.Println("[+] Requisitando informações dos templates")
	}

	req, reqErr := http.NewRequest(http.MethodGet, faasTemplateStoreURL, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("error ao criar uma solicitação para pegar informações do template: %s", reqErr.Error())
	}

	reqContext, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	req = req.WithContext(reqContext)
	client := http.DefaultClient
	res, clientErr := client.Do(req)
	if clientErr != nil {
		return nil, fmt.Errorf("erro ao solicitar a lista de templates: %s", clientErr.Error())
	}

	if res.Body == nil {
		return nil, fmt.Errorf("erro ao solicitar a lista de templates, boby vazio: %s", faasTemplateStoreURL)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status da requisição inesperado: %d got: %d", http.StatusOK, res.StatusCode)
	}

	body, bodyErr := ioutil.ReadAll(res.Body)
	if bodyErr != nil {
		return nil, fmt.Errorf("error ao ler conteúdo da resposta: %s", bodyErr.Error())
	}

	var templatesInfo []TemplateInfo
	unmarshallErr := json.Unmarshal(body, &templatesInfo)
	if unmarshallErr != nil {
		return nil, fmt.Errorf("error ao executar o unmarshalling da estrutura do template: %s", unmarshallErr.Error())
	}

	return templatesInfo, nil
}

func outputTemplateInfo(templates []TemplateInfo, verboseFlag bool) {
	if len(templates) == 0 {
		fmt.Println("Sem dados!")
		return
	}

	if verboseFlag {
		fmt.Println("[+] Processando informações dos templates")
	}

	var buff bytes.Buffer
	lineWriter := tabwriter.NewWriter(&buff, 0, 0, 3, ' ', 0)

	fmt.Fprintln(lineWriter)
	fmt.Fprintf(lineWriter, "NAME\tLANGUAGE\tPLATFORM\tSOURCE\tDESCRIPTION\n")
	for _, template := range templates {
		fmt.Fprintf(lineWriter, "%s\t%s\t%s\t%s\t%s\n",
			template.TemplateName,
			template.Language,
			template.Platform,
			template.Source,
			template.Description)
	}

	fmt.Fprintln(lineWriter)

	lineWriter.Flush()

	fmt.Println(buff.String())
}
