package main

import (
	"sync"
	"time"
)

const MAX_ROUTINES = 4

type ServiceI interface {
	DownloadVideos(videosIDs []string) error
}

type Service struct {
	components *Components
}

func NewService(components *Components) *Service {
	return &Service{
		components: components,
	}
}

func (s *Service) DownloadVideos(videosIDs []string) error {
	// max go routines
	maxRoutines := len(videosIDs)
	if maxRoutines > MAX_ROUTINES {
		maxRoutines = MAX_ROUTINES
	}

	sem := make(chan struct{}, maxRoutines)
	work := func(wg *sync.WaitGroup, id string) {
		defer wg.Done()

		sem <- struct{}{}
		err := YoutubeDownloadClient(id)
		if err != nil {
			s.components.Set(id, "error", time.Hour)
		} else {
			s.components.Set(id, "ok", time.Hour)
		}

		<-sem
	}

	var waitGroup sync.WaitGroup
	for _, id := range videosIDs {
		// Check cache
		cacheVal, err := s.components.Get(id)
		if err != nil {
			return err
		}

		if cacheVal == "ok" {
			println(id + " alredy downloaded.")
			continue
		}
		if cacheVal == "created" {
			println(id + " currently downloading.")
			continue
		}

		err = s.components.Set(id, "created", time.Hour)
		if err != nil {
			return err
		}

		waitGroup.Add(1)
		go work(&waitGroup, id)
	}

	waitGroup.Wait()
	return nil
}
