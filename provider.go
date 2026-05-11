package aitui

import "context"

// Provider describes an AI backend that can create chat clients.
type Provider interface {
	// Option returns the provider metadata shown to callers and selectors.
	Option() ProviderOption

	// Available reports whether the provider can be used in the current environment.
	Available(ctx context.Context) (bool, error)

	// NewClient creates a chat client for this provider.
	NewClient(ctx context.Context) (Client, error)
}

// ProviderOption is the public metadata for a provider choice.
type ProviderOption struct {
	// ID uniquely identifies the provider.
	ID string

	// Name is the human-readable provider name.
	Name string

	// Description explains when or why to choose this provider.
	Description string
}
