package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"strings"
	"time"
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
		data, err := request(fetchURL, "GET", []byte(""), 0)
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

func request(url, method string, requestData []byte, timeout int64) ([]byte, error) {
	if timeout == 0 {
		timeout = 60
	}
	trans := http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: 100,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	// https://github.com/golang/go/issues/13801
	client := &http.Client{
		Transport: &trans,
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout*int64(time.Second)))
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(requestData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-agent", "SealosRuntime")
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request %s resp code is %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	defer io.Copy(ioutil.Discard, resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
