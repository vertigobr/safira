// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"text/tabwriter"

	"gopkg.in/gookit/color.v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runInfraStatus,
	SuggestionsMinimumDistance: 1,
}

func init() {
	infraCmd.AddCommand(statusCmd)
}

func runInfraStatus(cmd *cobra.Command, args []string) error {
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	k8sClient, err := getClient("/home/lcortes/.config/k3d/vertigo-ipaas/kubeconfig.yaml")
	if err != nil {
		return err
	}

	if err := outputStatus(k8sClient, verboseFlag); err != nil {
		return err
	}

	return nil
}

func getClient(kubeconfig string) (*kubernetes.Clientset, error){
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error ao criar client, verifique o kubeconfig: %s", err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error ao verificar kubeconfig: %s", err.Error())
	}

	return client, nil
}

func outputStatus(client *kubernetes.Clientset, verboseFlag bool) error {
	deploymentsClient := client.AppsV1().Deployments("default")
	list, _ := deploymentsClient.List(context.TODO(), v1.ListOptions{})

	deploymentsFunctions := client.AppsV1().Deployments("openfaas-fn")
	listFunction, _ := deploymentsFunctions.List(context.TODO(), v1.ListOptions{})

	if verboseFlag {
		fmt.Println("[+] Processando informações dos deployments")
	}

	var buff bytes.Buffer
	lineWriter := tabwriter.NewWriter(&buff, 0, 0, 3, ' ', 0)

	fmt.Fprintln(lineWriter)
	fmt.Fprintf(lineWriter, color.Bold.Sprintf("SERVICES\n"))
	fmt.Fprintf(lineWriter, "NAME\t\t STATUS\n")
	for _, d := range list.Items {
		checkStatus := d.Status.AvailableReplicas == d.Status.Replicas
		var status string
		if checkStatus {
			status = color.Green.Sprintf("Ready")
		} else {
			status = color.Red.Sprintf("Not Ready")
		}

		fmt.Fprintf(lineWriter, "%s\t%s\n",
			d.Name,
			status,
		)
	}

	if len(listFunction.Items) > 0 {
		fmt.Fprintln(lineWriter)
		fmt.Fprintf(lineWriter, color.Bold.Sprintf("FUNCTIONS\n"))
		fmt.Fprintf(lineWriter, "NAME\t\t STATUS\n")
		for _, d := range listFunction.Items {
			checkStatus := d.Status.AvailableReplicas == d.Status.Replicas
			var status string
			if checkStatus {
				status = color.Green.Sprintf("Ready")
			} else {
				status = color.Red.Sprintf("Not Ready")
			}

			fmt.Fprintf(lineWriter, "%s\t%s\n",
				d.Name,
				status,
			)
		}
	}

	//fmt.Fprintln(lineWriter)

	lineWriter.Flush()

	fmt.Println(buff.String())

	return nil
}
