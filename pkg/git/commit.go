// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"gopkg.in/gookit/color.v1"
)

func GetImageWithCommitSha(image string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%s Error getting path to current folder", color.Red.Text("[!]"))
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	h, err := repo.ResolveRevision("HEAD")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	untaggedImage := strings.Split(image, ":")

	return strings.Replace(image, untaggedImage[len(untaggedImage)-1], h.String()[:7], 1), nil
}
