package service

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9-_.@]{3,25}$`)
var passwordRegex = regexp.MustCompile(`^[^\s]{5,20}$`)

// Service for managing authentication requests.
type AuthenticationService struct {
	authRepo repositories.AuthenticationRepository
	userRepo repositories.UserRepository
	tokens   *TokenBucket

	cacheNeedsAdminRegistration *bool
}

func NewAuthenticationService(
	authRepo repositories.AuthenticationRepository,
	userRepo repositories.UserRepository,
	tokenks *TokenBucket,
) *AuthenticationService {
	s := &AuthenticationService{
		authRepo: authRepo,
		userRepo: userRepo,
		tokens:   tokenks,
	}

	return s
}

func (s *AuthenticationService) Authenticate(
	username string,
	password string,
) (Token, error) {
	if err := s.validateCredentials(username, password); err != nil {
		return Token{}, err
	}

	u, err := s.authRepo.SignIn(username, s.hash(password))
	if err != nil {
		return Token{}, err
	}

	return s.tokens.New(u), nil
}

func (s *AuthenticationService) RegisterAdmin(
	username string,
	password string,
) (Token, error) {
	if !s.NeedsAdminRegistration() {
		return Token{}, fmt.Errorf("one admin account is already registered")
	}

	if err := s.validateCredentials(username, password); err != nil {
		return Token{}, err
	}

	password = s.hash(password)
	u, err := s.authRepo.RegisterAdmin(username, password)
	if err != nil {
		return Token{}, err
	}

	*s.cacheNeedsAdminRegistration = false
	return s.tokens.New(u), nil
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

func (s *AuthenticationService) hash(pt string) string {
	hash := sha512.New()
	hash.Write([]byte(pt))
	bs := hash.Sum(nil)

	return hex.EncodeToString(bs)
}

func (s *AuthenticationService) validateCredentials(
	username string,
	password string,
) error {
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username does not match pattern")
	}

	if !passwordRegex.MatchString(password) {
		return fmt.Errorf("password does not match pattern")
	}

	return nil
}
