// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/k8s"
	"gopkg.in/gookit/color.v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var infraStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "View the status of cluster services and functions",
	Long:  "View the status of cluster services and functions",
	Example: `To view the status of services and functions in the cluster, run:

    $ safira infra status`,
	RunE:                       runInfraStatus,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(infraStatusCmd)
}

func runInfraStatus(cmd *cobra.Command, _ []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	k8sClient, err := k8s.GetClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("%s No clusters found", color.Red.Text("[!]"))
	}

	if err := outputStatus(k8sClient, verboseFlag); err != nil {
		return err
	}

	return nil
}

func outputStatus(client *kubernetes.Clientset, verboseFlag bool) error {
	deploymentsClient := client.AppsV1().Deployments("default")
	list, _ := deploymentsClient.List(context.TODO(), v1.ListOptions{})

	deploymentsFunctions := client.AppsV1().Deployments(functionsNamespace)
	listFunction, _ := deploymentsFunctions.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Printf("%s Processing cluster information\n", color.Blue.Text("[v]"))
	}

	var buff bytes.Buffer
	lineWriter := tabwriter.NewWriter(&buff, 0, 0, 3, ' ', 0)

	_, _ = fmt.Fprintln(lineWriter)
	_, _ = fmt.Fprintf(lineWriter, color.Bold.Sprintf("SERVICES\n"))
	_, _ = fmt.Fprintf(lineWriter, "NAME\t\t    STATUS\t    AVAILABILITY\t\t URL\n")
	for _, d := range list.Items {
		checkStatus := d.Status.AvailableReplicas == d.Status.Replicas
		var status string
		if checkStatus {
			status = color.Green.Sprintf("Ready")
		} else {
			status = color.Red.Sprintf("Not Ready")
		}

		deployName := d.Name
		if strings.HasPrefix(deployName, "vtg-ipaas-") {
			deployName = strings.Split(deployName, "vtg-ipaas-")[1]
		}

		_, _ = fmt.Fprintf(lineWriter, "%s\t\t%s\t%s\t\t\t\t%s\n",
			deployName,
			fmt.Sprintf("%v/%v", d.Status.AvailableReplicas, d.Status.Replicas),
			status,
			getUrl(deployName, false),
		)
	}

	if len(listFunction.Items) > 0 {
		_, _ = fmt.Fprintln(lineWriter)
		_, _ = fmt.Fprintf(lineWriter, color.Bold.Sprintf("FUNCTIONS\n"))
		_, _ = fmt.Fprintf(lineWriter, "NAME\t\t    STATUS\t    AVAILABILITY\t\t URL\n")
		for _, d := range listFunction.Items {
			checkStatus := d.Status.AvailableReplicas == d.Status.Replicas
			var status string
			if checkStatus {
				status = color.Green.Sprintf("Ready")
			} else {
				status = color.Red.Sprintf("Not Ready")
			}

			_, _ = fmt.Fprintf(lineWriter, "%s\t\t%s\t%s\t\t\t\t%s\n",
				d.Name,
				fmt.Sprintf("%v/%v", d.Status.AvailableReplicas, d.Status.Replicas),
				status,
				getUrl(d.Name, true),
			)
		}
	}

	_ = lineWriter.Flush()

	fmt.Println(buff.String())

	return nil
}

func getUrl(deployName string, function bool) string {
	switch deployName {
	case "swaggereditor":
		return "editor.localdomain:8080"
	case "gateway":
		return "gateway.ipaas.localdomain:8080"
	case "kong":
		return "ipaas.localdomain:8080"
	case "konga":
		return "konga.localdomain:8080"
	default:
		break
	}

	if function {
		return "ipaas.localdomain:8080/function/" + deployName
	} else if strings.HasSuffix(deployName, "swagger-ui") {
		return "ipaas.localdomain:8080/swagger-ui/" + strings.Split(deployName, "-swagger-ui")[0]
	}

	return ""
}
