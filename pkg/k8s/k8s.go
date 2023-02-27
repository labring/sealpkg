package k8s

import (
	"fmt"
	"github.com/labring-actions/runtime-ctl/pkg/utils"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"strings"
)

func FetchFinalVersion(kubeVersion string) (string, error) {
	kubeVersion = strings.ReplaceAll(kubeVersion, "v", "")
	versionSplits := strings.Split(kubeVersion, ".")
	if len(versionSplits) == 3 {
		klog.Infof("before version is %s,new version is %s", kubeVersion, kubeVersion)
		return kubeVersion, nil
	}
	fetchURL := fmt.Sprintf("https://dl.k8s.io/release/stable-%s.txt", strings.Join(versionSplits, "."))
	latestVersion := ""
	klog.Infof("current version is %s, fetch url is %s", kubeVersion, fetchURL)
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		data, err := utils.Request(fetchURL, "GET", []byte(""), 0)
		if err != nil {
			return err
		}
		latestVersion = strings.ReplaceAll(string(data), "v", "")
		return nil

	}); err != nil {
		return "", err
	}
	return latestVersion, nil
}
