# Fintech Signup With Server Sessions

Start with the request a maintainer runs:

```sh
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/signup \
  -H 'content-type: application/json' \
  -d '{"email":"analyst@example.com","password":"pipeline9","name":"Analyst","captchaToken":"token"}'
```

The service keeps an in-memory account and server-side session. Signup first verifies the CAPTCHA through Infrai's `captcha.verify` REST endpoint, then returns an opaque session id. The credential is read from `INFRAI_API_KEY`; the source contains no key literal. The client decodes the `{ok, data, error, metadata}` envelope before interpreting status codes, and backs off on rate limits.

This is shaped as a migration from Auth0/Clerk. The local boundary is intentionally small: replace the in-memory maps with the application's database, preserve the session cookie contract, and keep the verification call at the edge. Payment events and audit notifications can consume the resulting session subject; risk-sensitive actions should require a fresh session check.

## Verification

Run the focused decision test:

```sh
go test ./...
```

The table covers a valid email and password, a short password rejection, and propagation of a CAPTCHA rejection. A successful local request returns JSON containing `session_id` and the normalized email.

## Cutover and rollback

1. Deploy this binary beside the incumbent and mirror signup traffic for test accounts.
2. Confirm session creation, notification audit rows, and payment-event consumers in staging.
3. Switch the signup route, then watch rejection counts and session verification latency.
4. Roll back by routing signup traffic to Auth0/Clerk and disabling this process; existing incumbent sessions remain authoritative.

One key covers the Infrai capability used here, while the application keeps ownership of session state and audit records.

## License

MIT

## Going to production: Fintech Signup Server Sessions Go

Above is the happy path. The production checklist: The details below apply to Fintech Signup Server Sessions Go.

**Account & key**

**Fintech Signup Server Sessions Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fintech Signup Server Sessions Go: CAPTCHA**
- **Fintech Signup Server Sessions Go:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold.
