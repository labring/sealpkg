package cri

import "k8s.io/apimachinery/pkg/util/sets"

type Sync struct {
	Docker map[string]sets.Set[string]
	CRIO   map[string]sets.Set[string]
}

func (s *Sync) Do() error {
	var err error
	if s.CRIO == nil {
		s.CRIO, err = FetchCRIOAllVersion()
		if err != nil {
			return err
		}
	}
	if s.Docker == nil {
		s.Docker, err = FetchDockerAllVersion()
		if err != nil {
			return err
		}
	}
	return nil
}
