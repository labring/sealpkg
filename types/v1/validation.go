package v1

import "fmt"

func ValidationDefaultComponent(c *RuntimeConfigDefaultComponent) error {
	if c.CRIO == "" {
		return fmt.Errorf("crio default version is empty,please retry config it")
	}
	if c.CRIOCrun == "" {
		return fmt.Errorf("crio-runc default version is empty,please retry config it")
	}
	if c.Docker == "" {
		return fmt.Errorf("docker default version is empty,please retry config it")
	}
	if c.CRIDocker == "" {
		return fmt.Errorf("cri-docker default version is empty,please retry config it")
	}
	if c.Containerd == "" {
		return fmt.Errorf("containerd default version is empty,please retry config it")
	}
	if c.Sealos == "" {
		return fmt.Errorf("sealos default version is empty,please retry config it")
	}
	return nil
}

func ValidationConfigData(c *RuntimeConfigData) error {
	if c.CRI == "" {
		return fmt.Errorf("cri not set,please retry config it")
	}
	if c.Runtime == "" {
		return fmt.Errorf("runtime not set,please retry config it")
	}
	if c.RuntimeVersion == "" {
		return fmt.Errorf("runtime version not set,please retry config it")
	}
	return nil
}
