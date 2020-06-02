// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the  Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package get

import (
	"runtime"
	"strings"
)

func getARCH() string {
	arch := runtime.GOARCH

	if strings.HasPrefix(arch, "armv7") {
		return "arm"
	} else if strings.HasPrefix(arch, "aarch64") {
		return "arm64"
	}

	return strings.ToLower(arch)
}

func getOS() string {
	os := runtime.GOOS

	if strings.Contains(strings.ToLower(os), "mingw") {
		return "windows"
	}

	return strings.ToLower(os)
}

