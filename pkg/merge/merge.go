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
