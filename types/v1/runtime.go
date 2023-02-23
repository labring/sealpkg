package v1

type RuntimeConfigDefaultComponent struct {
	Containerd  string `json:"containerd"`
	Docker      string `json:"docker"`
	CRIDockerV3 string `json:"cri-docker-v3"`
	CRIDockerV2 string `json:"cri-docker-v2"`
	Sealos      string `json:"sealos"`
	CRIO        string `json:"crio"`
	CRIOCrun    string `json:"crio-crun"`
}

type RuntimeConfigData struct {
	CRI            string `json:"cri"`
	Runtime        string `json:"runtime"`
	RuntimeVersion string `json:"runtimeVersion"`
}

type RuntimeConfig struct {
	Config  *RuntimeConfigData             `json:"config,omitempty"`
	Default *RuntimeConfigDefaultComponent `json:"default,omitempty"`
}
