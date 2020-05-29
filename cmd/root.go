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
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/gookit/color.v1"
)

var cfgFile string
var safiraInit = color.Bold.Sprintf("sudo -E safira init")
var notExistBinary = fmt.Sprintf("\nDependência(s) em falta, execute: %s", safiraInit)

const (
	kubectlBinaryName = "kubectl"
	k3dBinaryName = "k3d"
	helmBinaryName = "helm"
	faasBinaryName = "faas-cli"
)

var rootCmd = &cobra.Command{
	Use:           "safira",
	Short:         "O Safira é uma ferramenta de auxílio ao Vertigo iPaaS",
	Long:          "O Safira é uma ferramenta para auxiliar os desenvolvedores no Vertigo iPaaS",
	Version:       "v0.0.1-beta.2",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		if len(e) != 0 {
			fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		}
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".safira" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".safira")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
