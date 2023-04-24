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

package v1

import (
	"k8s.io/apimachinery/pkg/util/yaml"
	"testing"
)

func Test_NewRuntimeConfigFromFile(t *testing.T) {
	yConfig := `
config:
  cri: docker
  runtime: k8s
  runtimeVersion: 1.23.0
`
	c := &RuntimeConfig{}
	err := yaml.Unmarshal([]byte(yConfig), c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func Test_NewRuntimeDefaultFromFile(t *testing.T) {
	yConfig := `
default:
  containerd: 1.5.0
  docker: 1.5.0
  cri-docker-v2: 1.5.0
  sealos: 4.1.5-rc1
  crio: 1.2.0
  crio-crun: 1.0.0
`
	c := &RuntimeConfig{}
	err := yaml.Unmarshal([]byte(yConfig), c)
	if err != nil {
		t.Errorf(err.Error())
	}
}
