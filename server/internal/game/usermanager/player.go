package usermanager

import "github.com/lonng/nano/session"

const (
	kCurPlayer = "player"
)

type Player struct {
	uid int64 // 用户ID
	// 玩家数据
	session *session.Session
}

func newPlayer(s *session.Session, uid int64) *Player {
	p := &Player{
		uid: uid,
	}
	p.bindSession(s)

	return p
}

func (p *Player) UID() int64 {
	return p.uid
}

func (p *Player) Session() *session.Session {
	return p.session
}

func (p *Player) bindSession(s *session.Session) {
	p.session = s
	p.session.Set(kCurPlayer, p)
}

func (p *Player) removeSession() {
	p.session.Remove(kCurPlayer)
	p.session = nil
}
