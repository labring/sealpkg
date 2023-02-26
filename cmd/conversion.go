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

package cmd

import (
	"github.com/labring-actions/runtime-ctl/pkg/apply"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"github.com/spf13/cobra"
)

var conversionFiles []string
var defaultFile string
var yamlEnable bool
var applier *apply.Applier

const printInfo = `All Version:
	cri-docker: https://github.com/Mirantis/cri-dockerd/releases
	docker: https://download.docker.com/linux/static/stable/
	containerd: https://github.com/containerd/containerd/releases
	crun: https://github.com/containers/crun/releases
	runc: https://github.com/opencontainers/runc/releases
	sealos: https://github.com/labring/sealos/releases
	crio: https://github.com/cri-o/cri-o/releases
`

// conversionCmd represents the conversion command
var conversionCmd = &cobra.Command{
	Use:   "conversion",
	Short: "conversion runtime cri and release version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return applier.Apply()
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		klog.Info(printInfo)
		applier = apply.NewApplier()
		if err := applier.WithDefaultFile(defaultFile); err != nil {
			return errors.WithMessage(err, "validate default error")
		}
		if err := applier.WithConfigFiles(conversionFiles...); err != nil {
			return errors.WithMessage(err, "validate config error")
		}
		return applier.WithYaml(yamlEnable)
	},
}

func init() {
	rootCmd.AddCommand(conversionCmd)
	conversionCmd.Flags().StringSliceVarP(&conversionFiles, "files", "f", []string{}, "config files")
	conversionCmd.Flags().StringVarP(&defaultFile, "default", "d", "", "default file location")
	conversionCmd.Flags().BoolVar(&yamlEnable, "yaml", false, "print yaml")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// conversionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// conversionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
