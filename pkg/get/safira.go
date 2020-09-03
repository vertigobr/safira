// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vertigobr/safira/pkg/config"
	"github.com/vertigobr/safira/pkg/execute"
	"gopkg.in/gookit/color.v1"
)

const safiraUrl = "https://github.com/vertigobr/safira"

func DownloadSafira(cliVersion string, verboseFlag bool) (string, error) {
	arch := getARCH()
	os := getOS()

	if verboseFlag {
		fmt.Printf("%s Checking platform support\n", color.Blue.Text("[v]"))
	}

	if err := verifySafiraSupport(arch, os); err != nil {
		return "", err
	}

	if verboseFlag {
		fmt.Printf("%s Checking for a new update\n", color.Blue.Text("[v]"))
	}

	tag, err := checkLatestVersion(cliVersion)
	if err != nil {
		return "", err
	}

	safiraUrl := getSafiraUrl(arch, os, tag)

	if verboseFlag {
		fmt.Printf("%s Downloading Safira\n", color.Blue.Text("[v]"))
	}

	if err := downloadBinary(safiraUrl, "safira", false); err != nil {
		return "", err
	}

	binaryTempFolder := config.GetSafiraDir() + "bin/safira"
	mv := execute.Task{
		Command: "sudo",
		Args: []string{
			"mv", binaryTempFolder, "/usr/local/bin/safira",
		},
	}

	res, err := mv.Execute()
	if err != nil {
		return "", err
	}

	if res.ExitCode != 0 {
		return "", fmt.Errorf(res.Stderr)
	}

	return tag, nil
}

func getSafiraUrl(arch, os, tag string) string {
	return fmt.Sprintf("%s/releases/download/%s/safira-%s-%s-%s.tar.gz", safiraUrl, tag, tag, os, arch)
}

func verifySafiraSupport(arch, os string) error {
	if os != "linux" {
		return fmt.Errorf("%s Safira is only available for linux", color.Red.Text("[!]"))
	}

	distribution := strings.Join([]string{os, arch}, "-")
	availables := []string{"linux-386", "linux-amd64", "linux-arm", "linux-arm64"}

	for _, available := range availables {
		if distribution == available {
			return nil
		}
	}

	return fmt.Errorf("%s Safira is not available for %s arch", color.Red.Text("[!]"), arch)
}

func checkLatestVersion(cliVersion string) (string, error) {
	repoUrlLatest := safiraUrl + "/releases/latest"

	res, err := http.DefaultClient.Get(repoUrlLatest)
	if err != nil {
		return "", fmt.Errorf("%s Error when obtaining request body: %s", color.Red.Text("[!]"), repoUrlLatest)
	}
	defer res.Body.Close()

	repoTagLatestSplit := strings.Split(res.Request.URL.String(), "/")
	latestTag := repoTagLatestSplit[len(repoTagLatestSplit)-1]
	if latestTag != cliVersion {
		return latestTag, nil
	}

	return "", fmt.Errorf("%s The latest version of Safira is already installed", color.Red.Text("[!]"))
}
