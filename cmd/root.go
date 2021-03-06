// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vertigobr/safira/pkg/config"
	"gopkg.in/gookit/color.v1"
)

var (
	cfgFile        string
	safiraInit     = color.Bold.Sprintf("sudo -E safira init")
	notExistBinary = fmt.Sprintf("\n%s Missing dependencies, run: %s", color.Red.Text("[!]"), safiraInit)
	kubeconfigPath = fmt.Sprintf("%s/.config/k3d/%s/kubeconfig.yaml", os.Getenv("HOME"), clusterName)
)

const (
	faasTemplateStoreURL = "https://raw.githubusercontent.com/vertigobr/openfaas-templates/master/templates.json"
	faasTemplateRepo     = "https://github.com/vertigobr/openfaas-templates.git"
	kubectlBinaryName    = "kubectl"
	k3dBinaryName        = "k3d"
	helmBinaryName       = "helm"
	faasBinaryName       = "faas-cli"
	oktetoBinaryName     = "okteto"
	clusterName          = "vertigo-ipaas"
	functionsNamespace   = "ipaas-fn"
)

var rootCmd = &cobra.Command{
	Use:           "safira",
	Short:         "Safira is a toolkit for Vertigo iPaaS",
	Long:          "Safira is a toolkit for Vertigo iPaaS",
	Version:       "v0.0.20",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		if len(e) != 0 {
			fmt.Println(e)
		}
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	setPath()
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

func setPath() {
	path := config.GetSafiraDir() + "bin:" + os.Getenv("PATH")
	_ = os.Setenv("PATH", path)
}
