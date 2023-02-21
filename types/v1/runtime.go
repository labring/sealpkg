package v1

type RuntimeConfigDefault struct {
	RuntimeConfigDefaultComponent
	K8s *RuntimeConfigDefaultComponent
}

type RuntimeConfigDefaultComponent struct {
	Containerd string `json:"containerd"`
	Docker     string `json:"docker"`
	CRIDocker  string `json:"cri-docker"`
	Sealos     string `json:"sealos"`
	CRIO       string `json:"crio"`
	CRIOCrun   string `json:"crio-crun"`
}

type RuntimeConfigData struct {
	CRI            string `json:"cri"`
	Runtime        string `json:"runtime"`
	RuntimeVersion string `json:"runtimeVersion"`
}

type RuntimeConfig struct {
	Config  *RuntimeConfigData    `json:"config,omitempty"`
	Default *RuntimeConfigDefault `json:"default,omitempty"`
}
