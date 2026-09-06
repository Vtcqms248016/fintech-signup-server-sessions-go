package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"
)

type signupInput struct{ Email, Password, Name, CaptchaToken, IP string }
type session struct{ ID, Email string }
type captchaVerifier interface{ verify(string, string) error }
type authService struct {
	captcha  captchaVerifier
	users    map[string]signupInput
	sessions map[string]session
}

func newAuthService(c captchaVerifier) *authService {
	return &authService{captcha: c, users: map[string]signupInput{}, sessions: map[string]session{}}
}

func (s *authService) signup(in signupInput) (session, error) {
	if _, err := mail.ParseAddress(in.Email); err != nil {
		return session{}, fmt.Errorf("invalid email")
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if len(in.Password) < 8 {
		return session{}, fmt.Errorf("password must have 8 characters")
	}
	if err := s.captcha.verify(in.CaptchaToken, in.IP); err != nil {
		return session{}, err
	}
	if _, exists := s.users[in.Email]; exists {
		return session{}, fmt.Errorf("account already exists")
	}
	s.users[in.Email] = in
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return session{}, err
	}
	id := hex.EncodeToString(idBytes)
	sess := session{ID: id, Email: in.Email}
	s.sessions[id] = sess
	return sess, nil
}
