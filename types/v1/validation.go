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

import (
	"fmt"
	"k8s.io/klog/v2"
	"strings"
)

func ValidationDefaultComponent(c *RuntimeConfigDefaultComponent) error {
	if c.CRIO == "" {
		return fmt.Errorf("crio default version is empty,please retry config it")
	}
	if c.Crun == "" {
		return fmt.Errorf("crio-runc default version is empty,please retry config it")
	}
	if c.Docker == "" {
		return fmt.Errorf("docker default version is empty,please retry config it")
	}
	if c.CRIDockerV2 == "" {
		return fmt.Errorf("cri-docker v2 default version is empty,please retry config it")
	}
	if c.CRIDockerV3 == "" {
		return fmt.Errorf("cri-docker v3 default version is empty,please retry config it")
	}
	if c.Containerd == "" {
		return fmt.Errorf("containerd default version is empty,please retry config it")
	}
	if c.Sealos == "" {
		return fmt.Errorf("sealos default version is empty,please retry config it")
	}
	klog.Info("validate default component success")
	return nil
}

func ValidationConfigData(c *RuntimeConfigData) error {
	if c.CRI == "" {
		return fmt.Errorf("cri not set,please retry config it")
	}
	if c.Runtime == "" {
		return fmt.Errorf("runtime not set,please retry config it")
	}
	if c.RuntimeVersion == "" {
		return fmt.Errorf("runtime version not set,please retry config it")
	}
	klog.Info("validate config data success")
	return nil
}

func ValidationRuntimeConfig(c *RuntimeConfig) error {
	if c.Config.Runtime == "k8s" {
		//kubernetes gt 1.26
		if Compare(c.Config.RuntimeVersion, "v1.26") && !Compare(c.Default.Sealos, "v4.1.3") {
			// echo "INFO::skip $KUBE(kube>=1.26) when $SEALOS(sealos<=4.1.3)"
			//  echo https://kubernetes.io/blog/2022/11/18/upcoming-changes-in-kubernetes-1-26/#cri-api-removal
			klog.Info("Please see https://kubernetes.io/blog/2022/11/18/upcoming-changes-in-kubernetes-1-26/#cri-api-removal")
			return fmt.Errorf("skip $KUBE(kube>=1.26) when $SEALOS(sealos<=4.1.3)")
		}
	}
	klog.Info("validate runtime success")
	return nil
}

// Compare is version compare
// if v1 >= v2 return true, else return false
func Compare(v1, v2 string) bool {
	v1 = strings.Replace(v1, "v", "", -1)
	v2 = strings.Replace(v2, "v", "", -1)
	v1 = strings.Split(v1, "-")[0]
	v2 = strings.Split(v2, "-")[0]
	v1List := strings.Split(v1, ".")
	v2List := strings.Split(v2, ".")
	if v1List[0] > v2List[0] {
		return true
	} else if v1List[0] < v2List[0] {
		return false
	}
	if v1List[1] > v2List[1] {
		return true
	} else if v1List[1] < v2List[1] {
		return false
	}
	if len(v1List) == 3 && len(v2List) == 3 {
		if v1List[2] >= v2List[2] {
			return true
		}
		return false
	}

	return true
}
