// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/gookit/color.v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func OutputFunctionLog(functionName, kubeconfig, namespace, outputFlag string) error {
	client, err := GetClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("%s Not was possible communication with the cluster", color.Red.Text("[!]"))
	}

	pods := client.CoreV1().Pods(namespace)
	podList, _ := pods.List(context.TODO(), v1.ListOptions{})

	for _, pod := range podList.Items {
		if strings.HasPrefix(pod.Name, functionName) {
			podLogOpts := corev1.PodLogOptions{}
			podLog := client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)

			log, err := podLog.Stream(context.TODO())
			if err != nil {
				return err
			}
			defer log.Close()

			logBuffer := new(bytes.Buffer)
			_, err = io.Copy(logBuffer, log)
			if err != nil {
				return fmt.Errorf("\nError ao imprimir log!")
			}

			if len(outputFlag) > 0 {
				f, err := os.OpenFile(outputFlag, os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					return fmt.Errorf("error ao abrir o arquivo %s: %s", outputFlag, err.Error())
				}
				defer f.Close()

				f.Truncate(0)

				_, err = f.Write(logBuffer.Bytes())
				if err != nil {
					return fmt.Errorf("error ao escrever no arquivo %s: %s", outputFlag, err.Error())
				}
			} else {
				fmt.Println(logBuffer.String())
			}

			return nil
		}
	}

	return fmt.Errorf("%s Function %s not found", color.Red.Text("[!]"), functionName)
}
