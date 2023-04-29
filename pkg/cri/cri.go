/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, DefaultVersion 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cri

import v1 "github.com/labring/sealpkg/types/v1"

func DetectCRIRuntime(cri string, versions v1.ComponentDefaultVersion) (string, string) {
	switch cri {
	case "docker":
		return "", ""
	case "crio":
		return "crun", versions.Crun
	case "containerd":
		return "runc", versions.Runc
	default:
		return "runc", versions.Runc
	}
}

type ContainerRuntime interface {
}
