/*
Copyright 2014 The Kubernetes Authors.

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

package version

import (
	"fmt"
)

// Info contains versioning information.
// TODO: Add []string of api versions supported? It's still unclear
// how we'll want to distribute that information.
type Info struct {
	GitVersion string `json:"gitVersion" yaml:"gitVersion"`
	GitCommit  string `json:"gitCommit,omitempty" yaml:"gitCommit,omitempty"`
	BuildDate  string `json:"buildDate" yaml:"buildDate"`
	GoVersion  string `json:"goVersion" yaml:"goVersion"`
	Compiler   string `json:"compiler" yaml:"compiler"`
	Platform   string `json:"platform" yaml:"platform"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return fmt.Sprintf("%s-%s", info.GitVersion, info.GitCommit)
}

type DefaultVersion struct {
	CRIDockerd3x string `json:"criDockerd3x,omitempty"`
	CRIDockerd2x string `json:"criDockerd2x,omitempty"`
	Dockerd18    string `json:"docker18,omitempty"`
	Dockerd19    string `json:"docker19,omitempty"`
	Dockerd20    string `json:"docker20,omitempty"`
}

type Output struct {
	Version        Info           `json:"version" yaml:"version"`
	DefaultVersion DefaultVersion `json:"defaultVersion" yaml:"defaultVersion"`
}

// String returns info as a human-friendly version string.
func (info Output) String() string {
	return info.Version.String()
}
