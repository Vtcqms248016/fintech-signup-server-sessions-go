package main

import (
	"fmt"
	"testing"
)

type captchaStub struct{ err error }

func (c captchaStub) verify(string, string) error { return c.err }

func TestSignupDecision(t *testing.T) {
	cases := []struct {
		name       string
		in         signupInput
		captchaErr bool
		wantErr    bool
	}{
		{"valid signup creates session", signupInput{Email: "a@example.com", Password: "password8", CaptchaToken: "ok"}, false, false},
		{"short password rejected", signupInput{Email: "b@example.com", Password: "short", CaptchaToken: "ok"}, false, true},
		{"captcha rejection returned", signupInput{Email: "c@example.com", Password: "password8", CaptchaToken: "bad"}, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := newAuthService(captchaStub{})
			if tc.captchaErr {
				svc.captcha = captchaStub{err: fmt.Errorf("rejected")}
			}
			_, err := svc.signup(tc.in)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error=%v wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr && len(svc.sessions) != 1 {
				t.Fatalf("sessions=%d", len(svc.sessions))
			}
		})
	}
}
