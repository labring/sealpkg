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

package apply

import (
	"errors"
	"fmt"
	"github.com/labring-actions/runtime-ctl/pkg/cri"
	"github.com/labring-actions/runtime-ctl/pkg/k8s"
	"github.com/labring-actions/runtime-ctl/pkg/merge"
	v1 "github.com/labring-actions/runtime-ctl/types/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

type Applier struct {
	StatusComponents []v1.RuntimeStatusComponent
	RuntimeConfigs   []v1.RuntimeConfig
	DefaultFile      string
	Yaml             bool
}

func NewApplier() *Applier {
	return &Applier{}
}

func (a *Applier) WithDefaultFile(file string) error {
	if file == "" {
		return errors.New("file not set,please set file retry")
	}
	cfg, err := merge.Merge(file)
	if err != nil {
		return err
	}
	if err = v1.ValidationDefaultComponent(cfg.Version); err != nil {
		return err
	}
	a.DefaultFile = file
	return nil
}

func (a *Applier) WithYaml(yamlEnable bool) error {
	a.Yaml = yamlEnable
	return nil
}

func (a *Applier) WithConfigFiles(files ...string) error {
	if len(files) <= 0 {
		return errors.New("files not set,please set retry")
	}
	validationFunc := func(index int, r *v1.RuntimeConfig) error {
		if err := v1.ValidationConfigData(r.Config); err != nil {
			return err
		}
		if err := v1.ValidationRuntimeConfig(r); err != nil {
			return err
		}
		klog.Infof("validate index=%d config data and runtime success", index)
		return nil
	}
	versions := sets.NewString()
	for i, f := range files {
		cfg, err := merge.Merge(f, a.DefaultFile)
		if err != nil {
			return err
		}
		if err = validationFunc(i, cfg); err != nil {
			return fmt.Errorf("file is %s is validation error: %+v", f, err)
		}
		if cfg.Config.CRI == nil || len(cfg.Config.CRI) == 0 {
			cfg.Config.CRI = []string{v1.CRIContainerd, v1.CRIDocker, v1.CRICRIO}
		}
		for _, v := range cfg.Config.RuntimeVersion {
			for _, r := range cfg.Config.CRI {
				setKey := fmt.Sprintf("%s-%s-%s", r, cfg.Config.Runtime, v)
				if !versions.Has(setKey) {
					versions.Insert(setKey)
					rt := v1.RuntimeStatusComponent{
						CRIType:        r,
						Runtime:        cfg.Config.Runtime,
						RuntimeVersion: v,
					}
					rt.CRIRuntime, rt.CRIRuntimeVersion = cri.CRIRuntime(r, *cfg.Version)
					rt.Sealos = cfg.Version.Sealos
					a.StatusComponents = append(a.StatusComponents, rt)
					a.RuntimeConfigs = append(a.RuntimeConfigs, *cfg)
				}
			}
		}

	}
	return nil
}

func (a *Applier) Apply() error {
	statusList := &v1.RuntimeStatusList{}
	for i, rt := range a.StatusComponents {
		localRuntime := a.StatusComponents[i]
		switch rt.Runtime {
		case v1.RuntimeK8s:
			switch rt.CRIType {
			case v1.CRIDocker:
				dockerVersion, criDockerVersion := cri.FetchDockerVersion(rt.RuntimeVersion)
				if dockerVersion != "" {
					localRuntime.CRIVersion = dockerVersion
				} else {
					localRuntime.CRIVersion = a.RuntimeConfigs[i].Version.Docker
				}
				localRuntime.CRIDockerd = criDockerVersion
			case v1.CRIContainerd:
				localRuntime.CRIVersion = a.RuntimeConfigs[i].Version.Containerd
			}

			newVersion, err := k8s.FetchFinalVersion(rt.RuntimeVersion)
			if err != nil {
				return fmt.Errorf("runtime is %s,runtime version is %s,get new version is error: %+v", rt.Runtime, rt.RuntimeVersion, err)
			}
			a.StatusComponents[i].RuntimeVersion = newVersion
		default:
			return fmt.Errorf("not found runtime,current version not support")
		}
		a.StatusComponents[i] = localRuntime
	}
	statusList.Include = a.StatusComponents
	if a.Yaml {
		actionYAML, err := yaml.Marshal(statusList)
		if err != nil {
			return err
		}
		fmt.Println(string(actionYAML))
		return nil
	}
	actionJSON, err := json.Marshal(statusList)
	if err != nil {
		return err
	}
	fmt.Println(string(actionJSON))
	return nil
}
