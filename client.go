package aitui

import "context"

// Client sends chat requests to a selected provider.
type Client interface {
	// Chat sends a single-turn chat request and returns the provider response.
	Chat(ctx context.Context, request ChatRequest) (ChatResponse, error)
}

// ChatRequest is a single-turn chat request.
type ChatRequest struct {
	// Input is the user prompt to send to the provider.
	Input string

	// Model optionally overrides the provider's default model.
	Model string

	// Instructions optionally provide system-level guidance for the request.
	Instructions string
}

// ChatResponse is a single-turn chat response.
type ChatResponse struct {
	// Text is the assistant text returned by the provider.
	Text string

	// Provider identifies the provider that produced the response.
	Provider string

	// Model identifies the model that produced the response when available.
	Model string
}
