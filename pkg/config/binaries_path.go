// Copyright © 2020 Vertigo Tecnologia. All rights reserved.
// Licensed under the Apache License, Version 2.0. See LICENSE file in the project root for full license information.
package config

import "fmt"

func GetKubectlPath() string {
	return fmt.Sprintf("%sbin/kubectl", GetSafiraDir())
}

func GetK3dPath() string {
	return fmt.Sprintf("%sbin/k3d", GetSafiraDir())
}

func GetHelmPath() string {
	return fmt.Sprintf("%sbin/helm", GetSafiraDir())
}

func GetFaasCliPath() string {
	return fmt.Sprintf("%sbin/faas-cli", GetSafiraDir())
}
