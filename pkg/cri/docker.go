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

package cri

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/sealpkg/pkg/retry"
	"github.com/labring-actions/sealpkg/pkg/utils"
	"github.com/labring-actions/sealpkg/pkg/version"
	v1 "github.com/labring-actions/sealpkg/types/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
)

func FetchDockerVersion(kubeVersion string) (string, string) {
	var dockerVersion string
	switch {
	//# kube 1.16(docker-18.09)
	case !v1.Compare(kubeVersion, "v1.17"):
		dockerVersion = version.Docker18
	//# kube 1.17-20(docker-19.03)
	case v1.Compare(kubeVersion, "v1.17") && !v1.Compare(kubeVersion, "v1.21"):
		dockerVersion = version.Docker19
	//kube 1.21-23(docker-20.10)
	case v1.Compare(kubeVersion, "v1.21") && !v1.Compare(kubeVersion, "v1.24"):
		dockerVersion = version.Docker20
	default:
		dockerVersion = ""
	}

	var dockerCRIVersion string
	switch {
	//# kube 1.1x-25(cri-dockerd v0.2.x)
	case !v1.Compare(kubeVersion, "v1.26"):
		dockerCRIVersion = version.CRIDockerd2x
	//# kube 1.26-2x(cri-dockerd v0.3.x)
	case v1.Compare(kubeVersion, "v1.26"):
		dockerCRIVersion = version.CRIDockerd3x
	}

	return dockerVersion, dockerCRIVersion
}

func FetchDockerAllVersion() (map[string]sets.Set[string], error) {
	fetchURL := "https://download.docker.com/linux/static/stable/x86_64/"
	versions := make(map[string]sets.Set[string])
	if err := retry.Retry(func() error {
		data, err := utils.Request(fetchURL, "GET", []byte(""), 0)
		if err != nil {
			return err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
		if err != nil {
			return err
		}
		ahtml := doc.Find("a")
		for _, html := range ahtml.Nodes {
			attr := html.Attr
			if len(attr) > 0 {
				if strings.Contains(attr[0].Val, "docker") && !strings.Contains(attr[0].Val, "rootless") && !strings.Contains(attr[0].Val, "ce") {
					//docker-18.09.2.tgz
					tmpVal := strings.ReplaceAll(attr[0].Val, "docker-", "")
					tmpVal = strings.ReplaceAll(tmpVal, ".tgz", "")
					if len(strings.Split(tmpVal, ".")) < 3 {
						continue
					}
					bigVersion := strings.Join(strings.Split(tmpVal, ".")[:2], ".")
					if _, ok := versions[bigVersion]; !ok {
						versions[bigVersion] = sets.New(tmpVal)
					} else {
						versions[bigVersion].Insert(tmpVal)
					}
				}
			}
		}
		return nil
	}); err != nil {
		logger.Error("get docker version error: %s", err.Error())
		return nil, err
	}
	return versions, nil
}
