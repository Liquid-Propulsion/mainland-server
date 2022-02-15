package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var CurrentSessionManager = NewSessionManager()

var timeTillAutoLock = time.Minute
var timeTillAutoExpire = time.Hour * 4

func NewSessionManager() *SessionManager {
	manager := new(SessionManager)
	manager.sessions = make(map[uuid.UUID]Session)
	return manager
}

type SessionManager struct {
	sessions map[uuid.UUID]Session
}

type Session struct {
	UUID       uuid.UUID
	UserID     uint
	Locked     bool
	LockTime   time.Time
	ExpiryTime time.Time
}

func (manager *SessionManager) Add(userID uint) Session {
	for {
		uid := uuid.New()
		if _, ok := manager.sessions[uid]; !ok {
			session := Session{
				UUID:       uid,
				UserID:     userID,
				Locked:     false,
				LockTime:   time.Now().Add(timeTillAutoLock),
				ExpiryTime: time.Now().Add(timeTillAutoExpire),
			}
			manager.sessions[uid] = session
			return session
		}
	}
}

func (manager *SessionManager) Remove(sessionID string) error {
	uid, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}
	if _, ok := manager.sessions[uid]; ok {
		delete(manager.sessions, uid)
		return nil
	}
	return errors.New("session doesnt exist")
}

func (manager *SessionManager) PreventLock(sessionID string) (*Session, error) {
	uid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, err
	}
	if _, ok := manager.sessions[uid]; ok {
		session := manager.sessions[uid]
		session.LockTime = time.Now().Add(timeTillAutoLock)
		session.ExpiryTime = time.Now().Add(timeTillAutoExpire)
		manager.sessions[uid] = session
		return &session, nil
	}
	return nil, errors.New("session doesnt exist")
}

func (manager *SessionManager) SetLock(sessionID string, locked bool) (*Session, error) {
	uid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, err
	}
	if _, ok := manager.sessions[uid]; ok {
		session := manager.sessions[uid]
		session.Locked = locked
		if !locked {
			session.LockTime = time.Now().Add(timeTillAutoLock)
			session.ExpiryTime = time.Now().Add(timeTillAutoExpire)
		}
		manager.sessions[uid] = session
		return &session, nil
	}
	return nil, errors.New("session doesnt exist")
}

func (manager *SessionManager) VerifySession(sessionID string) (*Session, error) {
	uid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, err
	}
	if _, ok := manager.sessions[uid]; ok {
		session := manager.sessions[uid]
		if !session.Locked && session.LockTime.After(time.Now()) {
			return &session, nil
		}
		session.Locked = true
		manager.sessions[uid] = session
		return &session, errors.New("session is locked")
	}
	return nil, errors.New("session doesnt exist")
}
