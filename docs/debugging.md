# Debugging And Secret Safety

- Never print API keys, access tokens, refresh tokens, browser sessions, official-tool private auth state, or raw provider session payloads.
- Do not store secrets in `SessionStore`, logs, error messages, fixtures, snapshots, or debug output.
- Redact known secret fields as `[REDACTED]`, including `api_key`, `access_token`, `refresh_token`, `authorization`, `cookie`, `session`, and provider-specific token fields.
- Redact bearer headers as `Authorization: Bearer [REDACTED]`.
- Redact long unknown credential-like values when they appear in auth, header, cookie, or session contexts.
- Keep enough context to debug failures: provider ID, operation name, HTTP status, command exit code, timeout/cancellation state, and safe stderr summaries.
- Auth failures should identify the provider and missing or failed auth path, for example: `provider openai-api: missing API key`.
- API failures should include provider, operation, and status without response secrets, for example: `provider openai-api: chat request failed: HTTP 429`.
- Selector cancellation should be a normal user-cancel path, for example: `provider selection canceled`.
- Command-backed providers should capture stderr for diagnostics, redact it before returning errors, and avoid streaming raw command output directly to logs.
- Downstream CLIs should show user-actionable messages, keep verbose provider diagnostics behind an explicit debug flag, and run all debug output through the same redaction path.
- Downstream CLIs that persist credentials must own that storage and redaction policy; `aitui` should only receive the credential material needed for the current setup or client runtime.
