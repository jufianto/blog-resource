// Package base holds the concept-neutral plumbing shared by every Day 2 example:
// the in-memory go-oidc storage and small web/OAuth-client helpers. Each cmd/ example
// keeps its own Authorization Server config and Resource Server so the concept it
// teaches (JWT validation, DPoP, PAR, FAPI, ...) stays visible in one place.
package base

import (
	"context"
	"sync"

	"github.com/luikyv/go-oidc/pkg/goidc"
)

// MemoryStore implements the go-oidc storage interfaces with plain maps. It is
// identical to the Day 1 store — storage is not what Day 2 teaches, so it is shared.
type MemoryStore struct {
	mu       sync.RWMutex
	grants   map[string]*goidc.Grant
	tokens   map[string]*goidc.Token
	sessions map[string]*goidc.AuthnSession
	clients  map[string]*goidc.Client
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		grants:   make(map[string]*goidc.Grant),
		tokens:   make(map[string]*goidc.Token),
		sessions: make(map[string]*goidc.AuthnSession),
		clients:  make(map[string]*goidc.Client),
	}
}

func (s *MemoryStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grants = make(map[string]*goidc.Grant)
	s.tokens = make(map[string]*goidc.Token)
	s.sessions = make(map[string]*goidc.AuthnSession)
	// Pre-registered clients are kept; they are seeded once at startup.
}

func (s *MemoryStore) SaveGrant(_ context.Context, grant *goidc.Grant) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.grants[grant.ID] = grant
	return nil
}

func (s *MemoryStore) Grant(_ context.Context, id string) (*goidc.Grant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.grants[id]
	if !ok {
		return nil, goidc.ErrNotFound
	}
	return g, nil
}

func (s *MemoryStore) SaveToken(_ context.Context, token *goidc.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token.ID] = token
	return nil
}

func (s *MemoryStore) Token(_ context.Context, id string) (*goidc.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tokens[id]
	if !ok {
		return nil, goidc.ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) SaveSession(_ context.Context, session *goidc.AuthnSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
	return nil
}

func (s *MemoryStore) Session(_ context.Context, id string) (*goidc.AuthnSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[id]
	if !ok {
		return nil, goidc.ErrNotFound
	}
	return sess, nil
}

func (s *MemoryStore) GrantByAuthCode(_ context.Context, code string) (*goidc.Grant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, g := range s.grants {
		if g.AuthCode == code {
			return g, nil
		}
	}
	return nil, goidc.ErrNotFound
}

func (s *MemoryStore) GrantByRefreshToken(_ context.Context, token string) (*goidc.Grant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, g := range s.grants {
		if g.RefreshToken == token {
			return g, nil
		}
	}
	return nil, goidc.ErrNotFound
}

func (s *MemoryStore) SaveClient(_ context.Context, client *goidc.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.ID] = client
	return nil
}

func (s *MemoryStore) Client(_ context.Context, id string) (*goidc.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.clients[id]
	if !ok {
		return nil, goidc.ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) DeleteClient(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, id)
	return nil
}

func (s *MemoryStore) SessionByUserCode(_ context.Context, code string) (*goidc.AuthnSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sess := range s.sessions {
		if sess.UserCode == code {
			return sess, nil
		}
	}
	return nil, goidc.ErrNotFound
}

func (s *MemoryStore) SessionByDeviceCode(_ context.Context, code string) (*goidc.AuthnSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sess := range s.sessions {
		if sess.DeviceCode == code {
			return sess, nil
		}
	}
	return nil, goidc.ErrNotFound
}

// SessionByPushedAuthReqID satisfies goidc.PARManager: the /authorize endpoint uses it to
// resolve the session created by a pushed authorization request (PAR).
func (s *MemoryStore) SessionByPushedAuthReqID(_ context.Context, id string) (*goidc.AuthnSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sess := range s.sessions {
		if sess.PushedAuthReqID == id {
			return sess, nil
		}
	}
	return nil, goidc.ErrNotFound
}

func (s *MemoryStore) GrantByDeviceCode(_ context.Context, code string) (*goidc.Grant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, g := range s.grants {
		if g.DeviceCode == code {
			return g, nil
		}
	}
	return nil, goidc.ErrNotFound
}
