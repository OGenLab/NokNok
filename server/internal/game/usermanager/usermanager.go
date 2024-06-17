package usermanager

import (
	"errors"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	"qoobing.com/gomod/log"
)

var DefaultUserManger = NewUserManager()

type UserManager struct {
	component.Base
	players map[int64]*Player
}

func NewUserManager() *UserManager {
	return &UserManager{
		players: map[int64]*Player{},
	}
}

func (u *UserManager) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		if s.UID() > 0 {
			if err := u.onPlayerDisconnect(s); err != nil {
				log.Errorf("player quit: UID=%d, Error=%s", s.UID, err.Error())
			}
		}
	})
}

func (u *UserManager) onPlayerDisconnect(s *session.Session) error {
	uid := s.UID()
	p, err := playerWithSession(s)
	if err != nil {
		return err
	}
	log.Debugf("DeskManager.onPlayerDisconnect network disconnect uid:%v", uid)

	// 移除session
	p.removeSession()

	DefaultUserManger.offline(uid)
	return nil
}

func playerWithSession(s *session.Session) (*Player, error) {
	p, ok := s.Value(kCurPlayer).(*Player)
	if !ok {
		return nil, errors.New("player not found")
	}
	return p, nil
}

func (u *UserManager) player(uid int64) (*Player, bool) {
	p, ok := u.players[uid]
	return p, ok
}

func (u *UserManager) setPlayer(uid int64, p *Player) {
	if _, ok := u.players[uid]; ok {
		log.Infof("the player already exists, overwriting the player, UID=%d", uid)
	}
	u.players[uid] = p
}

func (u *UserManager) sessionCount() int {
	return len(u.players)
}

func (u *UserManager) offline(uid int64) {
	delete(u.players, uid)
	log.Infof("player: %d Removed from online list, remaining: %d", uid, len(u.players))
}
