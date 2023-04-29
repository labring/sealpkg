// Copyright © 2023 sealos.
//
// Licensed under the Apache License, DefaultVersion 2.0 (the "License");
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
	"github.com/cuisongliu/logger"
	"github.com/labring/sealpkg/pkg/cri"
	"github.com/labring/sealpkg/pkg/k8s"
	"github.com/labring/sealpkg/pkg/sync"
	"github.com/labring/sealpkg/pkg/utils"
	v1 "github.com/labring/sealpkg/types/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
	"strings"
)

type Applier struct {
	Status  []v1.ComponentAndVersion
	Configs []v1.RuntimeConfig
	Yaml    bool
	Sync    *sync.Sync
}

func NewApplier() *Applier {
	return &Applier{}
}

func (a *Applier) defaultFile(cfg *v1.ComponentDefaultVersion) error {
	if err := v1.ValidationDefaultComponent(cfg); err != nil {
		return err
	}
	const printInfo = `All Default Version:
	docker: %s
	containerd: %s
	sealos: %s
	crun: %s
	runc: %s
`
	logger.Info(printInfo, cfg.Docker, cfg.Containerd, cfg.Sealos, cfg.Crun, cfg.Runc)
	return nil
}

func (a *Applier) WithCRISync(sync *sync.Sync) error {
	a.Sync = sync
	return nil
}

func (a *Applier) WithYaml(yamlEnable bool) error {
	a.Yaml = yamlEnable
	return nil
}

func (a *Applier) WithConfigFiles(file string) error {
	if file == "" || !utils.IsFileExist(file) {
		return errors.New("files not set or file is not exist,please set retry")
	}
	versions := sets.NewString()
	var cfg *v1.RuntimeConfig
	var err error
	cfg, err = v1.ReadFileToObject(file)
	if err != nil {
		return err
	}
	if err = v1.ValidationConfigData(cfg.Config); err != nil {
		return fmt.Errorf("file is %s is validation error: %+v", file, err)
	}
	logger.Debug("validate config data and runtime success")
	if cfg.Config.CRI == nil || len(cfg.Config.CRI) == 0 {
		cfg.Config.CRI = []string{v1.CRIContainerd, v1.CRIDocker, v1.CRICRIO}
	}
	const printInfo = `All Config:
	cri: %s
	runtime: %s
	versions: %s
`
	cris := fmt.Sprintf("[%s]", strings.Join(cfg.Config.CRI, ","))
	runtimeVersions := fmt.Sprintf("[%s]", strings.Join(cfg.Config.RuntimeVersion, ","))
	logger.Info(printInfo, cris, cfg.Config.Runtime, runtimeVersions)

	for _, v := range cfg.Config.RuntimeVersion {
		for _, r := range cfg.Config.CRI {
			setKey := fmt.Sprintf("%s-%s-%s", r, cfg.Config.Runtime, v)
			if !versions.Has(setKey) {
				versions.Insert(setKey)
				rt := v1.ComponentAndVersion{
					CRIType:        r,
					Runtime:        cfg.Config.Runtime,
					RuntimeVersion: v,
				}
				rt.CRIRuntime, rt.CRIRuntimeVersion = cri.DetectCRIRuntime(r, *cfg.DefaultVersion)
				rt.Sealos = cfg.DefaultVersion.Sealos
				if err = v1.CheckSealosAndRuntime(cfg.Config, cfg.DefaultVersion); err != nil {
					logger.Warn("check sealos and runtime error: %+v", err)
					continue
				}

				newVersions := k8s.FetchK8sAllVersion(rt.RuntimeVersion)
				for _, vv := range newVersions {
					rt.RuntimeVersion = vv
					a.Status = append(a.Status, rt)
					a.Configs = append(a.Configs, *cfg)
				}
			}
		}
	}

	return nil
}

func (a *Applier) Apply() error {
	statusList := &v1.RuntimeList{}
	for i, rt := range a.Status {
		localRuntime := a.Status[i]
		switch rt.Runtime {
		case v1.RuntimeK8s:
			kubeBigVersion := v1.ToBigVersion(rt.RuntimeVersion)
			switch rt.CRIType {
			case v1.CRIDocker:
				dockerVersion, criDockerVersion := cri.FetchDockerVersion(rt.RuntimeVersion)
				if dockerVersion != "" {
					localRuntime.CRIVersion = dockerVersion
				} else {
					cfgDocker := a.Configs[i].DefaultVersion.Docker
					localRuntime.CRIVersion = v1.ToBigVersion(cfgDocker)
				}
				versions := a.Sync.Docker[localRuntime.CRIVersion]
				sortList := utils.List(versions)
				newVersion := sortList[len(sortList)-1]
				logger.Debug("docker version is %s, docker using latest version: %s", localRuntime.CRIVersion, newVersion)
				localRuntime.CRIVersion = newVersion
				localRuntime.CRIDockerd = criDockerVersion
			case v1.CRIContainerd:
				localRuntime.CRIVersion = a.Configs[i].DefaultVersion.Containerd
				containerdBigVersion := v1.ToBigVersion(a.Configs[i].DefaultVersion.Containerd)
				if v1.Compare(kubeBigVersion, "1.26") && !v1.Compare(containerdBigVersion, "1.6") {
					localRuntime.CRIVersion = ""
					logger.Warn("if kubernetes version gt 1.26 your containerd must be gt 1.6")
				}
			case v1.CRICRIO:
				versions := a.Sync.CRIO[kubeBigVersion]
				sortList := utils.List(versions)
				newVersion := sortList[len(sortList)-1]
				logger.Debug("kube version is %s, crio using latest version: %s", rt.RuntimeVersion, newVersion)
				localRuntime.CRIVersion = newVersion
			}
		default:
			return fmt.Errorf("not found runtime,current version not support")
		}
		if localRuntime.CRIVersion == "" {
			continue
		}
		statusList.Include = append(statusList.Include, localRuntime)
	}
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
