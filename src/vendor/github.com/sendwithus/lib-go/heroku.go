package swu

import (
	"time"

	"github.com/bgentry/heroku-go"
	"log"
	"sync"
)

type HerokuService interface {
	StartDedicatedDyno(appName string, command string, size string, attach bool) (*heroku.Dyno, error)
	GetDynoFormationCached(herokuAppName string, dynoName string) (*heroku.Formation, error)
	GetDynoFormation(herokuAppName string, dynoName string) (*heroku.Formation, error)
	GetAllDynoFormations(herokuAppName string) ([]heroku.Formation, error)
	UpdateDynoFormation(appName string, dynoName string, quantity *int, size *string) (*heroku.Formation, error)
	UpdateDynoFormationCache(herokuAppName string) error
}

type herokuServiceImpl struct {
	apiKey string
}

type cachedFormationsList struct {
	LastUpdated time.Time
	Formations  []heroku.Formation
}

var hasCacheInited = false
var cachedFormationsLock = &sync.Mutex{}
var cachedFormations map[string]cachedFormationsList
var cacheTimeoutDuration = time.Duration(time.Second * 5)

func NewHerokuService(apiKey string) HerokuService {
	return herokuServiceImpl{
		apiKey: apiKey,
	}
}

func (s herokuServiceImpl) UpdateDynoFormationCache(herokuAppName string) error {
	cachedFormationsLock.Lock()
	defer cachedFormationsLock.Unlock()

	// initial build of cache
	if !hasCacheInited {
		cachedFormations = make(map[string]cachedFormationsList)
		hasCacheInited = true
	}

	cachedApp := cachedFormationsList{}
	if val, exists := cachedFormations[herokuAppName]; exists {
		cachedApp = val
	}

	if !cachedApp.LastUpdated.IsZero() {
		if time.Since(cachedApp.LastUpdated) < cacheTimeoutDuration {
			return nil
		}
	}

	// update the cached formations
	formations, err := s.GetAllDynoFormations(herokuAppName)
	if err != nil {
		cachedFormationsLock.Unlock()
		return err
	}

	cachedApp.Formations = formations
	cachedApp.LastUpdated = time.Now()
	cachedFormations[herokuAppName] = cachedApp

	return nil
}

func (s herokuServiceImpl) GetDynoFormationCached(herokuAppName string, dynoName string) (*heroku.Formation, error) {
	err := s.UpdateDynoFormationCache(herokuAppName)
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}

	// find formation in the cache
	thisCachedFormationsList, exists := cachedFormations[herokuAppName]
	if exists {
		for _, cachedFormation := range thisCachedFormationsList.Formations {
			if cachedFormation.Type == dynoName {
				return &cachedFormation, nil
			}
		}
	}

	// fallback to non-cached query
	log.Fatalf("No dyno found in cache named %v, in app %v", dynoName, herokuAppName)
	return s.GetDynoFormation(herokuAppName, dynoName)
}

func (s herokuServiceImpl) GetDynoFormation(herokuAppName string, dynoName string) (*heroku.Formation, error) {
	client := heroku.Client{Username: "", Password: s.apiKey}

	formation, err := client.FormationInfo(herokuAppName, dynoName)
	if err != nil {
		return nil, err
	}

	return formation, nil
}

func (s herokuServiceImpl) GetAllDynoFormations(herokuAppName string) ([]heroku.Formation, error) {
	client := heroku.Client{Username: "", Password: s.apiKey}
	return client.FormationList(herokuAppName, &heroku.ListRange{Field: "name", Max: 1000})
}

func (s herokuServiceImpl) UpdateDynoFormation(appName string, dynoName string, quantity *int, size *string) (*heroku.Formation, error) {
	client := heroku.Client{Username: "", Password: s.apiKey}

	update := heroku.FormationUpdateOpts{
		Quantity: quantity,
		Size:     size,
	}

	return client.FormationUpdate(appName, dynoName, &update)
}

func (s herokuServiceImpl) StartDedicatedDyno(appName string, command string, size string, attach bool) (*heroku.Dyno, error) {
	client := heroku.Client{Username: "", Password: s.apiKey}

	create_options := heroku.DynoCreateOpts{
		Attach: &attach,
		Size:   &size,
	}

	return client.DynoCreate(appName, command, &create_options)
}
