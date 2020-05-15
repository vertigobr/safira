/*
Copyright © Vertigo Tecnologia

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
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		installInfra()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func installInfra() {
	helmUrl := "https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3"

	// Get the data
	resp, err := http.Get(helmUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create("get_helm.sh")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	if path, _ := exec.LookPath("helm"); len(path) < 1 {
		execCommand("chmod", "700", "./get_helm.sh")
		execCommand("bash", "./get_helm.sh")
	} else {
		fmt.Println("Helm já instalado.")
	}

	if path, _ := exec.LookPath("minikube"); len(path) < 1 {
		execCommand("curl", "-Lo", "minikube", "https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64")
		execCommand("chmod", "+x", "minikube")
		execCommand("sudo", "install", "minikube", "/usr/local/bin/")
	} else {
		fmt.Println("Minikube já instalado.")
	}

	if path, _ := exec.LookPath("kubectl"); len(path) < 1 {
		execCommand("curl", "-LO", "https://storage.googleapis.com/kubernetes-release/release/v1.18.2/bin/linux/amd64/kubectl")
		execCommand("chmod", "+x", "kubectl")
		execCommand("sudo", "mv", "kubectl", "/usr/local/bin/kubectl")
	} else {
		fmt.Println("Kubectl já instalado.")
	}

	execCommand("minikube", "start", "--vm-driver=docker")
	execCommand("helm", "repo", "add", "vtg-ipaas", "https://vertigobr.gitlab.io/ipaas/vtg-ipaas-chart")
	execCommand("helm", "repo", "update")
	execCommand("helm", "upgrade", "-i", "-f", "https://gist.githubusercontent.com/kyfelipe/81a783a48025277e126b1dc8684fb420/raw/141b694e7bb605ec0e458b4c18115da02ddb9432/minikube.yaml", "vtg-ipaas", "vtg-ipaas/vtg-ipaas")
}

func execCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
