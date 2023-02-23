package docker

import v1 "github.com/labring-actions/runtime-ctl/types/v1"

const (
	CRIDockerV2 = "v0.2.x"
	CRIDockerV3 = "v0.3.x"
)

func FetchVersion(kubeVersion string) (string, string) {
	var dockerVersion string
	switch {
	//# kube 1.16(docker-18.09)
	case !v1.Compare(kubeVersion, "v1.17"):
		dockerVersion = "18.09"
	//# kube 1.17-20(docker-19.03)
	case v1.Compare(kubeVersion, "v1.17") && !v1.Compare(kubeVersion, "v1.21"):
		dockerVersion = "19.03"
	//kube 1.21-23(docker-20.10)
	case v1.Compare(kubeVersion, "v1.21") && !v1.Compare(kubeVersion, "v1.24"):
		dockerVersion = "20.10"
	default:
		dockerVersion = ""
	}

	var dockerCRIVersion string
	switch {
	//# kube 1.1x-25(cri-dockerd v0.2.x)
	case !v1.Compare(kubeVersion, "v1.26"):
		dockerCRIVersion = CRIDockerV2
	//# kube 1.26-2x(cri-dockerd v0.3.x)
	case v1.Compare(kubeVersion, "v1.26"):
		dockerCRIVersion = CRIDockerV3
	}

	return dockerVersion, dockerCRIVersion
}
