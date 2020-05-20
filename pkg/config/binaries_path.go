package config

import (
	"fmt"
)

func GetKubectlPath() string {
	return fmt.Sprintf("%sbin/kubectl", GetUserDir())
}

func GetK3dPath() string {
	return fmt.Sprintf("%sbin/k3d", GetUserDir())
}

func GetHelmPath() string {
	return fmt.Sprintf("%sbin/helm", GetUserDir())
}
