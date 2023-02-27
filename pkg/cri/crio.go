/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
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

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/labring-actions/runtime-ctl/pkg/utils"
	v1 "github.com/labring-actions/runtime-ctl/types/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"strings"
)

func FetchCRIOAllVersion() (map[string]sets.Set[string], error) {
	const crioAddress = "https://cri-o.github.io/cri-o/"
	versions := make(map[string]sets.Set[string])
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		data, err := utils.Request(crioAddress, "GET", []byte(""), 0)
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
				if strings.Contains(attr[0].Val, "cri-o") && !strings.Contains(attr[0].Val, "dependencies.html") {
					//docker-18.09.2.tgz
					tmpVal := strings.ReplaceAll(attr[0].Val, "/cri-o/", "")
					tmpVal = strings.ReplaceAll(tmpVal, ".html", "")
					tmpVal = strings.ReplaceAll(tmpVal, "v", "")
					if strings.Contains(tmpVal, "-") {
						continue
					}
					if !v1.Compare(tmpVal, "1.20") {
						continue
					}
					if len(strings.Split(tmpVal, ".")) < 3 {
						continue
					}
					//versions = append(versions, tmpVal)
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
		klog.Error("get docker version error: %s", err.Error())
		return nil, err
	}
	return versions, nil
}
