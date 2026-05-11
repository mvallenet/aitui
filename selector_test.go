package aitui

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestSelectClientEmptyProviders(t *testing.T) {
	_, err := SelectClient(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "no providers configured") {
		t.Fatalf("SelectClient() error = %v, want no providers configured", err)
	}
}

func TestSelectClientNoAvailableProviders(t *testing.T) {
	_, err := SelectClient(context.Background(), []Provider{
		fakeProvider{id: "unavailable", available: false},
	})
	if err == nil || !strings.Contains(err.Error(), "no available providers") {
		t.Fatalf("SelectClient() error = %v, want no available providers", err)
	}
}

func TestSelectClientOneAvailableProviderAutoSelects(t *testing.T) {
	want := fakeClient{text: "ok"}
	unavailable := &recordingProvider{fakeProvider: fakeProvider{id: "unavailable", available: false}}
	available := &recordingProvider{fakeProvider: fakeProvider{id: "available", available: true, client: want}}
	selector := &recordingSelector{}

	got, err := SelectClient(context.Background(), []Provider{unavailable, available}, selector)
	if err != nil {
		t.Fatalf("SelectClient() error = %v", err)
	}
	if got != want {
		t.Fatalf("SelectClient() client = %#v, want %#v", got, want)
	}
	if selector.called {
		t.Fatal("SelectClient() called selector for one available provider")
	}
	if available.newClientCalls != 1 {
		t.Fatalf("available.NewClient calls = %d, want 1", available.newClientCalls)
	}
	if unavailable.newClientCalls != 0 {
		t.Fatalf("unavailable.NewClient calls = %d, want 0", unavailable.newClientCalls)
	}
}

func TestSelectClientMultipleAvailableProvidersRequireSelector(t *testing.T) {
	_, err := SelectClient(context.Background(), []Provider{
		fakeProvider{id: "first", available: true},
		fakeProvider{id: "second", available: true},
	})
	if err == nil || !strings.Contains(err.Error(), "selector UI required") {
		t.Fatalf("SelectClient() error = %v, want selector UI required", err)
	}
}

func TestSelectClientMultipleAvailableProvidersUsesSelectorChoice(t *testing.T) {
	first := fakeClient{text: "first"}
	second := fakeClient{text: "second"}
	selector := &recordingSelector{selectedID: "second"}

	got, err := SelectClient(context.Background(), []Provider{
		fakeProvider{id: "first", name: "First", description: "first provider", available: true, client: first},
		fakeProvider{id: "second", name: "Second", description: "second provider", available: true, client: second},
	}, selector)
	if err != nil {
		t.Fatalf("SelectClient() error = %v", err)
	}
	if got != second {
		t.Fatalf("SelectClient() client = %#v, want %#v", got, second)
	}
	if !selector.called {
		t.Fatal("SelectClient() did not call selector")
	}
	if len(selector.options) != 2 {
		t.Fatalf("selector options len = %d, want 2", len(selector.options))
	}
	if selector.options[0].ID != "first" || selector.options[1].ID != "second" {
		t.Fatalf("selector option IDs = %q, %q; want first, second", selector.options[0].ID, selector.options[1].ID)
	}
}

func TestSelectClientSelectorUnknownProvider(t *testing.T) {
	_, err := SelectClient(context.Background(), []Provider{
		fakeProvider{id: "first", available: true},
		fakeProvider{id: "second", available: true},
	}, fakeSelector{selectedID: "missing"})
	if err == nil || !strings.Contains(err.Error(), "unknown provider") {
		t.Fatalf("SelectClient() error = %v, want unknown provider", err)
	}
}

func TestSelectClientReturnsSelectorError(t *testing.T) {
	wantErr := errors.New("selector canceled")

	_, err := SelectClient(context.Background(), []Provider{
		fakeProvider{id: "first", available: true},
		fakeProvider{id: "second", available: true},
	}, fakeSelector{err: wantErr})
	if !errors.Is(err, wantErr) {
		t.Fatalf("SelectClient() error = %v, want %v", err, wantErr)
	}
}

func TestSelectClientReturnsAvailabilityError(t *testing.T) {
	wantErr := errors.New("probe failed")

	_, err := SelectClient(context.Background(), []Provider{
		fakeProvider{id: "broken", availableErr: wantErr},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("SelectClient() error = %v, want %v", err, wantErr)
	}
}

func TestSelectClientReturnsNewClientError(t *testing.T) {
	wantErr := errors.New("client failed")

	_, err := SelectClient(context.Background(), []Provider{
		fakeProvider{id: "broken", available: true, newClientErr: wantErr},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("SelectClient() error = %v, want %v", err, wantErr)
	}
}

type fakeProvider struct {
	id           string
	name         string
	description  string
	available    bool
	availableErr error
	client       Client
	newClientErr error
}

func (p fakeProvider) Option() ProviderOption {
	return ProviderOption{
		ID:          p.id,
		Name:        p.name,
		Description: p.description,
	}
}

func (p fakeProvider) Available(context.Context) (bool, error) {
	return p.available, p.availableErr
}

func (p fakeProvider) NewClient(context.Context) (Client, error) {
	if p.newClientErr != nil {
		return nil, p.newClientErr
	}
	return p.client, nil
}

type recordingProvider struct {
	fakeProvider
	newClientCalls int
}

func (p *recordingProvider) NewClient(ctx context.Context) (Client, error) {
	p.newClientCalls++
	return p.fakeProvider.NewClient(ctx)
}

type fakeSelector struct {
	selectedID string
	err        error
}

func (s fakeSelector) SelectProvider(context.Context, []ProviderOption) (string, error) {
	return s.selectedID, s.err
}

type recordingSelector struct {
	selectedID string
	err        error
	called     bool
	options    []ProviderOption
}

func (s *recordingSelector) SelectProvider(_ context.Context, options []ProviderOption) (string, error) {
	s.called = true
	s.options = append([]ProviderOption(nil), options...)
	return s.selectedID, s.err
}

type fakeClient struct {
	text string
}

func (c fakeClient) Chat(context.Context, ChatRequest) (ChatResponse, error) {
	return ChatResponse{Text: c.text}, nil
}
