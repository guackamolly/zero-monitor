package service

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Service for managing authentication requests.
type AuthenticationService struct {
	authRepo repositories.AuthenticationRepository
	userRepo repositories.UserRepository
	tokens   map[string]models.User

	cacheNeedsAdminRegistration *bool
}

func NewAuthenticationService(
	authRepo repositories.AuthenticationRepository,
	userRepo repositories.UserRepository,
) *AuthenticationService {
	s := &AuthenticationService{
		authRepo: authRepo,
		userRepo: userRepo,
		tokens:   map[string]models.User{},
	}

	return s
}

func (s *AuthenticationService) Authenticate(
	username string,
	password string,
) (string, error) {
	u, err := s.authRepo.SignIn(username, password)
	if err != nil {
		return "", err
	}

	return s.tokenize(u), nil
}

func (s *AuthenticationService) RegisterAdmin(
	username string,
	password string,
) (string, error) {
	if !s.NeedsAdminRegistration() {
		return "", fmt.Errorf("one admin account is already registered")
	}

	password = s.hash(password)
	u, err := s.authRepo.RegisterAdmin(username, password)
	if err != nil {
		return "", err
	}

	*s.cacheNeedsAdminRegistration = false

	return s.tokenize(u), nil
}

func (s *AuthenticationService) NeedsAdminRegistration() bool {
	if s.cacheNeedsAdminRegistration != nil {
		return *s.cacheNeedsAdminRegistration
	}

	// retry at most 5 times if repo call fails
	for i := 0; i < 5; i++ {
		exists, err := s.userRepo.AdminExists()
		if err != nil {
			time.Sleep(150 * time.Millisecond)
			continue
		}

		needsAdminRegistration := !exists
		s.cacheNeedsAdminRegistration = &needsAdminRegistration
		return *s.cacheNeedsAdminRegistration
	}

	logging.LogWarning("couldn't guess if admin is registered or not. allowing admin registration")
	return true
}

func (s *AuthenticationService) tokenize(
	user models.User,
) string {
	// clear existing token before adding a new one
	for t, u := range s.tokens {
		if u == user {
			delete(s.tokens, t)
		}
	}

	token := models.UUID()
	s.tokens[token] = user

	return token
}

func (s *AuthenticationService) hash(pt string) string {
	hash := sha512.New()
	hash.Write([]byte(pt))
	bs := hash.Sum(nil)

	return hex.EncodeToString(bs)
}
