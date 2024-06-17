package entities

import (
	"sync"
	"time"

	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_entities_properties"
)

const (
	HAMMER_TYPE = 1 + iota
	GOPHER_TYPE
)

func init() {
	refreshEntities()
	go refresh()
}

type entities struct {
	gopherIds []uint64
	gopherSet map[uint64]*t_entities_properties.T_entities_properties

	hammerIds           []uint64
	hammerProbabilities []int
	hammerSet           map[uint64]*t_entities_properties.T_entities_properties

	sync.RWMutex
}

var EntSet = &entities{}

func (e *entities) GopherIds() []uint64 {
	e.RLock()
	defer e.RUnlock()
	return e.gopherIds
}

func (e *entities) GopherSet() map[uint64]*t_entities_properties.T_entities_properties {
	e.RLock()
	defer e.RUnlock()
	return e.gopherSet
}

func (e *entities) HammerIds() []uint64 {
	e.RLock()
	defer e.RUnlock()
	return e.hammerIds
}

func (e *entities) HammerProbabilities() []int {
	e.RLock()
	defer e.RUnlock()
	return e.hammerProbabilities
}

func (e *entities) HammerSet() map[uint64]*t_entities_properties.T_entities_properties {
	e.RLock()
	defer e.RUnlock()
	return e.hammerSet
}

func refreshEntities() {
	var (
		m                   = model.NewModel()
		gopherIds           = make([]uint64, 0)
		gopherSet           = make(map[uint64]*t_entities_properties.T_entities_properties)
		hammerIds           = make([]uint64, 0)
		hammerProbabilities = make([]int, 0)
		hammerSet           = make(map[uint64]*t_entities_properties.T_entities_properties)
	)

	defer m.Close()

	entities, err := t_entities_properties.Select(m.DB, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	for i := range entities {
		switch entities[i].F_type {
		case HAMMER_TYPE:
			hammerIds = append(hammerIds, entities[i].F_entities_id)
			hammerSet[entities[i].F_entities_id] = entities[i]
			hammerProbabilities = append(hammerProbabilities, entities[i].F_probability)
		case GOPHER_TYPE:
			gopherIds = append(gopherIds, entities[i].F_entities_id)
			gopherSet[entities[i].F_entities_id] = entities[i]
		}
	}

	EntSet.Lock()
	defer EntSet.Unlock()

	EntSet.gopherIds = gopherIds
	EntSet.gopherSet = gopherSet
	EntSet.hammerIds = hammerIds
	EntSet.hammerProbabilities = hammerProbabilities
	EntSet.hammerSet = hammerSet
}

func refresh() {
	for {
		time.Sleep(time.Duration(config.Instance().NokConfig.RefreshInterval) * time.Minute)
		refreshEntities()
	}
}
