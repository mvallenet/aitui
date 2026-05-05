# Testing

- Default tests must run with `go test ./...` and must not call live providers, real OAuth flows, official login commands, or real local account state.
- Selector tests should use fake providers and fake selector input to cover empty lists, unavailable providers, auto-selection, user selection, cancellation, and selector errors.
- API-key provider tests should use fake end user input and fake HTTP transports. Tests should cover missing keys, request construction, response parsing, provider errors, and secret redaction.
- OAuth or official-auth provider tests should use fake auth runners. Tests should cover setup success, setup failure, cancellation, command errors, and secret redaction.
- `SessionStore` tests should use in-memory fakes and assert that API keys, access tokens, refresh tokens, browser sessions, and official-tool private auth state are never stored.
- Client tests should use fake transports or fake command runners and cover successful single-turn chat, malformed responses, provider failures, context cancellation, and timeouts.
- Integration tests must be opt-in with an explicit environment variable, build tag, or both. They should skip clearly when credentials, tools, or accounts are missing.
- Test fixtures and golden files must use fake secrets only and must never contain real provider credentials or copied auth state.
