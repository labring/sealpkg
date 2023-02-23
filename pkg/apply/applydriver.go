package apply

import (
	"errors"
	"fmt"
	"github.com/labring-actions/runtime-ctl/pkg/merge"
	v1 "github.com/labring-actions/runtime-ctl/types/v1"
)

type Applier struct {
	Runtimes    []v1.RuntimeConfig
	DefaultFile string
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
	if err = v1.ValidationDefaultComponent(cfg.Default); err != nil {
		return err
	}
	a.DefaultFile = file
	return nil

}

func (a *Applier) WithConfigFiles(files ...string) error {
	if len(files) <= 0 {
		return errors.New("files not set,please set retry")
	}
	validationFunc := func(r *v1.RuntimeConfig) error {
		if err := v1.ValidationConfigData(r.Config); err != nil {
			return err
		}
		if err := v1.ValidationRuntimeConfig(r); err != nil {
			return err
		}
		return nil
	}

	for _, f := range files {
		cfg, err := merge.Merge(f, a.DefaultFile)
		if err != nil {
			return err
		}
		if err = validationFunc(cfg); err != nil {
			return fmt.Errorf("file is %s is validation error: %+v", f, err)
		}
		a.Runtimes = append(a.Runtimes, *cfg)
	}
	return nil
}

func (a *Applier) Apply() error {
	//if len(files) > 0 {
	//	return nil, errors.New("files not set,please retry")
	//}
	//cfg, err := merge.Merge(files...)
	//if err != nil {
	//	return nil, err
	//}

	return nil
}
