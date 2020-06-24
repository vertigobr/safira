// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vertigobr/safira/pkg/k8s"
	"gopkg.in/gookit/color.v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var infraStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Vizualiza o status dos serviços e funções do cluster local",
	Long:  "Vizualiza o status dos serviços e funções do cluster local",
	RunE:  runInfraStatus,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(infraStatusCmd)
}

func runInfraStatus(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	k8sClient, err := k8s.GetClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("cluster local não encontrado!\n")
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
		fmt.Println("[+] Processando informações dos deployments")
	}

	var buff bytes.Buffer
	lineWriter := tabwriter.NewWriter(&buff, 0, 0, 3, ' ', 0)

	fmt.Fprintln(lineWriter)
	fmt.Fprintf(lineWriter, color.Bold.Sprintf("SERVICES\n"))
	fmt.Fprintf(lineWriter, "NAME\t\t    STATUS\t\t AVAILABILITY\n")
	for _, d := range list.Items {
		checkStatus := d.Status.AvailableReplicas == d.Status.Replicas
		var status string
		if checkStatus {
			status = color.Green.Sprintf("Ready")
		} else {
			status = color.Red.Sprintf("Not Ready")
		}

		fmt.Fprintf(lineWriter, "%s\t\t%s\t%s\n",
			d.Name,
			fmt.Sprintf("%v/%v", d.Status.AvailableReplicas, d.Status.Replicas),
			status,
		)
	}

	if len(listFunction.Items) > 0 {
		fmt.Fprintln(lineWriter)
		fmt.Fprintf(lineWriter, color.Bold.Sprintf("FUNCTIONS\n"))
		fmt.Fprintf(lineWriter, "NAME\t\t    STATUS\t\t AVAILABILITY\n")
		for _, d := range listFunction.Items {
			checkStatus := d.Status.AvailableReplicas == d.Status.Replicas
			var status string
			if checkStatus {
				status = color.Green.Sprintf("Ready")
			} else {
				status = color.Red.Sprintf("Not Ready")
			}

			fmt.Fprintf(lineWriter, "%s\t\t%s\t%s\n",
				d.Name,
				fmt.Sprintf("%v/%v", d.Status.AvailableReplicas, d.Status.Replicas),
				status,
			)
		}
	}

	//fmt.Fprintln(lineWriter)

	lineWriter.Flush()

	fmt.Println(buff.String())

	return nil
}
