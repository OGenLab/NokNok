package battlemanager

import "github.com/lonng/nano/component"

type BattleManager struct {
	component.Base
}

var DefaultBattleManager = NewBattleManager()

func NewBattleManager() *BattleManager {
	return &BattleManager{}
}
