// Copyright © 2023 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

type RuntimeDefaultComponent struct {
	Containerd string `json:"containerd"`
	Docker     string `json:"docker"`
	Sealos     string `json:"sealos"`
	Crun       string `json:"crun"`
	Runc       string `json:"runc"`
}

type RuntimeStatusComponent struct {
	CRIType           string `json:"criType"`
	CRIVersion        string `json:"criVersion"`
	CRIDockerd        string `json:"criDockerd,omitempty"`
	CRIRuntime        string `json:"criRuntime"`
	CRIRuntimeVersion string `json:"criRuntimeVersion"`
	Sealos            string `json:"sealos"`
	Runtime           string `json:"runtime"`
	RuntimeVersion    string `json:"runtimeVersion"`
}

const (
	RuntimeK8s    string = "k8s"
	CRIDocker     string = "docker"
	CRIContainerd string = "containerd"
	CRICRIO       string = "crio"
)

type RuntimeConfigData struct {
	CRI            []string `json:"cri,omitempty"`
	Runtime        string   `json:"runtime"`
	RuntimeVersion []string `json:"runtimeVersion,omitempty"`
}

type RuntimeConfig struct {
	Config  *RuntimeConfigData       `json:"config,omitempty"`
	Version *RuntimeDefaultComponent `json:"version,omitempty"`
}

type RuntimeStatusList struct {
	Include []RuntimeStatusComponent `json:"include,omitempty"`
}
