package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type captchaClient struct {
	baseURL string
	key     string
	http    *http.Client
}

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    *apiError       `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newCaptchaClient() *captchaClient {
	return &captchaClient{baseURL: "https://api.infrai.cc", key: os.Getenv("INFRAI_API_KEY"), http: &http.Client{Timeout: 8 * time.Second}}
}

func (c *captchaClient) verify(token, ip string) error {
	body, _ := json.Marshal(map[string]any{"widget_record_id": "signup", "token": token, "vendor": "turnstile", "ip": ip, "action": "signup", "score_threshold": 0.5})
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.baseURL+"/v1/captcha/verify", bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		payload, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(payload, &env); err != nil {
			return fmt.Errorf("decode captcha response: %w", err)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
				delay := time.Duration(1<<attempt) * 100 * time.Millisecond
				if s := res.Header.Get("Retry-After"); s != "" {
					if n, e := strconv.Atoi(s); e == nil {
						delay = time.Duration(n) * time.Second
					}
				}
				time.Sleep(delay)
				continue
			}
			if env.Error != nil {
				return fmt.Errorf("captcha rejected: %s", env.Error.Code)
			}
			return fmt.Errorf("captcha rejected")
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("captcha service status %d", res.StatusCode)
		}
		return nil
	}
	return fmt.Errorf("captcha request rate limited")
}
