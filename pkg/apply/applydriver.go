package apply

import (
	"errors"
	"github.com/labring-actions/runtime-ctl/pkg/merge"
	v1 "github.com/labring-actions/runtime-ctl/types/v1"
)

type Applier struct {
	Runtime *v1.RuntimeConfig
}

func NewApplier(files ...string) (*Applier, error) {
	if len(files) > 0 {
		return nil, errors.New("files not set,please retry")
	}
	cfg, err := merge.Merge(files...)
	if err != nil {
		return nil, err
	}
	return &Applier{Runtime: cfg}, nil
}

func (a *Applier) Apply() error {
	if err := a.Validation(); err != nil {
		return err
	}

	return nil
}

func (a *Applier) Validation() error {
	if err := v1.ValidationDefaultComponent(a.Runtime.Default); err != nil {
		return err
	}
	if err := v1.ValidationConfigData(a.Runtime.Config); err != nil {
		return err
	}
	if err := v1.ValidationRuntimeConfig(a.Runtime); err != nil {
		return err
	}
	return nil
}
