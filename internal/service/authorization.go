package service

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type Token struct {
	Value  string
	User   models.User
	Expiry time.Time
}

type TokenBucket map[string]Token

func (b TokenBucket) New(user models.User) Token {
	// clear existing tokens before adding a new one
	for t, tk := range b {
		if tk.User == user {
			delete(b, t)
		}
	}

	token := Token{
		Value:  models.UUID(),
		User:   user,
		Expiry: time.Now().Add(24 * time.Hour),
	}

	b[token.Value] = token
	return token
}

func (b TokenBucket) Token(token string) (Token, bool) {
	if t, ok := b[token]; ok {
		return t, t.Expiry.After(time.Now())
	}

	return Token{}, false
}

// Service for managing authorization requests.
type AuthorizationService struct {
	tokens *TokenBucket
}

func NewAuthorizationService(
	tokens *TokenBucket,
) *AuthorizationService {
	return &AuthorizationService{
		tokens: tokens,
	}
}

func (s *AuthorizationService) HasAdminRights(token string) bool {
	if t, ok := s.tokens.Token(token); ok {
		return t.User.IsAdmin()
	}

	return false
}
