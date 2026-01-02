package storage

import (
	"sync"
	"time"
)

type CaptchaSession struct {
	Answer    string
	CreatedAt time.Time
}

type Storage struct {
	data sync.Map
	ttl  time.Duration
}

func NewStorage(ttl time.Duration) *Storage {
	s := &Storage{
		ttl: ttl,
	}
	go s.cleanupLoop()
	return s
}

func (s *Storage) Set(userID int64, answer string) {
	s.data.Store(userID, CaptchaSession{
		Answer:    answer,
		CreatedAt: time.Now(),
	})
}

func (s *Storage) Verify(userID int64, answer string) bool {
	val, ok := s.data.Load(userID)
	if !ok {
		return false
	}
	session := val.(CaptchaSession)
	s.data.Delete(userID)

	if time.Since(session.CreatedAt) > s.ttl {
		return false
	}
	return session.Answer == answer
}

func (s *Storage) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.data.Range(func(key, value interface{}) bool {
			session := value.(CaptchaSession)
			if now.Sub(session.CreatedAt) > s.ttl {
				s.data.Delete(key)
			}
			return true
		})
	}
}