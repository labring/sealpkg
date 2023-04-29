package k8s

import (
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/labring/sealpkg/pkg/retry"
	"github.com/labring/sealpkg/pkg/utils"
	v1 "github.com/labring/sealpkg/types/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"sort"
	"strings"
	"sync"
)

func fetchVersion(kubeVersion string) (string, error) {

	fetchURL := fmt.Sprintf("https://dl.k8s.io/release/stable-%s.txt", kubeVersion)
	latestVersion := ""
	logger.Debug("current version is %s, fetch url is %s", kubeVersion, fetchURL)
	if err := retry.Retry(func() error {
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
func fetchAllVersion(kubeVersion string) []string {
	versionSplits := strings.Split(kubeVersion, ".")
	bigVersion := strings.Join(versionSplits[:2], ".")

	allTagsChan := make(chan []Tag, 10)
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go fetchTags(i, allTagsChan, &wg)
	}
	wg.Wait()
	close(allTagsChan)
	var allTags []Tag
	for tags := range allTagsChan {
		allTags = append(allTags, tags...)
	}
	var returnTags []string
	for _, tag := range allTags {
		if strings.HasPrefix(tag.Name, "v"+bigVersion) && !strings.Contains(tag.Name, "-") {
			returnTags = append(returnTags, tag.Name)
		}
	}

	sort.Slice(returnTags, func(i, j int) bool {
		return !v1.Compare(returnTags[i], returnTags[j])
	})

	return returnTags
}

type Tag struct {
	Name string `json:"name"`
}

func fetchTags(page int, allTagsChan chan<- []Tag, wg *sync.WaitGroup) {
	defer wg.Done()

	fetchURL := fmt.Sprintf("https://api.github.com/repos/kubernetes/kubernetes/tags?page=%d&per_page=100", page)
	var tags []Tag
	if err := retry.Retry(func() error {
		data, err := utils.Request(fetchURL, "GET", []byte(""), 0)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &tags)
		if err != nil {
			return err
		}
		return nil

	}); err != nil {
		logger.Error("fetch tags failed, err is %s", err)
	}
	allTagsChan <- tags
}

func FetchK8sAllVersion(kubeVersion string) []string {
	kubeVersion = strings.ReplaceAll(kubeVersion, "v", "")
	versionSplits := strings.Split(kubeVersion, ".")
	if len(versionSplits) == 3 {
		if versionSplits[2] == "*" {
			vs := fetchAllVersion(kubeVersion)
			return vs
		}
		return []string{kubeVersion}
	}
	logger.Debug("before version is %s,new version is %s", kubeVersion, kubeVersion)
	v, _ := fetchVersion(kubeVersion)
	return []string{v}
}
