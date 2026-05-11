package aitui

import (
	"context"
	"errors"
	"fmt"
)

// SelectorUI chooses one provider from the available provider options.
type SelectorUI interface {
	// SelectProvider returns the ID of the selected provider.
	SelectProvider(ctx context.Context, options []ProviderOption) (string, error)
}

// SelectClient selects an available provider and returns a client for it.
func SelectClient(ctx context.Context, providers []Provider, selectors ...SelectorUI) (Client, error) {
	if len(providers) == 0 {
		return nil, errors.New("aitui: no providers configured")
	}

	available, options, err := availableProviders(ctx, providers)
	if err != nil {
		return nil, err
	}
	if len(available) == 0 {
		return nil, errors.New("aitui: no available providers")
	}
	if len(available) == 1 {
		return newClient(ctx, available[0])
	}

	if len(selectors) == 0 || selectors[0] == nil {
		return nil, errors.New("aitui: selector UI required for multiple available providers")
	}

	selectedID, err := selectors[0].SelectProvider(ctx, options)
	if err != nil {
		return nil, err
	}

	for _, provider := range available {
		if provider.Option().ID == selectedID {
			return newClient(ctx, provider)
		}
	}

	return nil, fmt.Errorf("aitui: selector returned unknown provider %q", selectedID)
}

func availableProviders(ctx context.Context, providers []Provider) ([]Provider, []ProviderOption, error) {
	available := make([]Provider, 0, len(providers))
	options := make([]ProviderOption, 0, len(providers))
	for _, provider := range providers {
		option := provider.Option()
		ok, err := provider.Available(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("aitui: check provider %q availability: %w", option.ID, err)
		}
		if !ok {
			continue
		}
		available = append(available, provider)
		options = append(options, option)
	}
	return available, options, nil
}

func newClient(ctx context.Context, provider Provider) (Client, error) {
	client, err := provider.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("aitui: create provider %q client: %w", provider.Option().ID, err)
	}
	return client, nil
}
