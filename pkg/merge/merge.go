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

package merge

import (
	"github.com/imdario/mergo"
	"github.com/labring-actions/runtime-ctl/types/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
)

func Merge(files ...string) (*v1.RuntimeConfig, error) {
	c := &v1.RuntimeConfig{}
	for _, f := range files {
		cfg, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		unmarshalConfig := new(v1.RuntimeConfig)
		err = yaml.Unmarshal(cfg, unmarshalConfig)
		if err != nil {
			return nil, err
		}
		if err = mergo.Merge(c, unmarshalConfig); err != nil {
			return nil, err
		}
	}
	return c, nil
}
