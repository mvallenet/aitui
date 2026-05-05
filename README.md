# aitui

Package `aitui` provides provider selection and chat client interfaces for Go
command-line applications.

`aitui` helps CLI tools offer AI features without hard-coding one backend or
building their own provider picker and setup screens. API-based providers can
ask for user-supplied credentials, while CLI-backed providers can delegate to
official tools such as Codex CLI and use their supported authentication flows.
The main features are:

- It gives Go CLIs a provider selector they can embed in their own commands.
- It can guide users through the setup path for the provider they choose.
- It returns a chat client for the selected provider.
- It keeps provider-specific behavior behind small adapters.
- It supports fake providers and clients for tests.

* * *

- [Install](#install)
- [Examples](#examples)
- [Interactive Setup](#interactive-setup)
- [Credential Flow](#credential-flow)
- [Providers](#providers)
- [Authentication](#authentication)
- [Testing](#testing)

* * *

## Install

With a correctly configured Go toolchain:

```sh
go get github.com/mvallenet/aitui
```

## Examples

Let's start with a CLI command that lets the user choose an AI backend, then sends one chat request through the selected provider:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mvallenet/aitui"
	"github.com/mvallenet/aitui/providers"
)

func main() {
	ctx := context.Background()
	prompt := "Write a short release note for v0.1.0."

	client, err := aitui.SelectClient(ctx, providers.Default())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Chat(ctx, aitui.ChatRequest{
		Input: prompt,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Output)
}
```

The application owns the command, prompt, output handling, and list of providers it wants to offer. `aitui` owns the common path of checking which providers are available, asking the user to choose when needed, and returning a client for the selected backend.

Selectors can be used inside any command framework. The CLI decides when to
invoke `aitui`:

```go
func runExplain(ctx context.Context, path string) error {
	prompt := fmt.Sprintf("Explain this file: %s", path)

	client, err := aitui.SelectClient(ctx, providers.Default())
	if err != nil {
		return err
	}

	resp, err := client.Chat(ctx, aitui.ChatRequest{Input: prompt})
	if err != nil {
		return err
	}

	fmt.Println(resp.Output)
	return nil
}
```

If only one provider is available, it can be selected automatically. If multiple providers are available, the configured selector UI asks the user which backend to use.

## Interactive Setup

The selector can present provider choices in the terminal and continue into the
setup flow required by the selected provider:

```txt
Which model do you want to use?

  [ ] OpenAI API key
  [ ] OpenAI subscription
```

Choosing an API-key provider can prompt for the key using a masked input:

```txt
Enter your OpenAI API key: ****
```

Choosing a subscription-backed provider can start the official authentication flow for that provider, such as running the official CLI login command and then returning to the calling application once the provider is ready.

The calling CLI still owns the command and the surrounding user experience. `aitui` provides the reusable selector, setup screens, availability checks, and
client creation.

## Credential Flow

API-key providers and subscription-backed providers use different ownership
models.

For an API-key provider, the calling CLI can give `aitui` a credential store.
`aitui` checks that store before prompting, and writes back to that store only
when the user chooses to remember the key:

```txt
user chooses OpenAI API key
        |
        v
aitui asks the calling CLI's store for an existing key
        |
        v
if no key exists, aitui prompts for one
        |
        v
provider receives the key for the request
        |
        v
if the user wants to remember it, aitui saves it through the calling CLI's store
```

The key is not retrieved back from `aitui`. The application owns the storage implementation and passes it in:

```go
type CredentialStore interface {
	Get(ctx context.Context, providerID string) (secret string, ok bool, err error)
	Set(ctx context.Context, providerID string, secret string) error
}

client, err := aitui.SelectClient(
	ctx,
	providers.Default(),
	aitui.WithCredentialStore(appCredentials),
)
```

`appCredentials` can use a config file, the operating system keychain, a password manager, or no persistence at all. `aitui` only uses the interface while selecting and preparing the provider.

For a subscription-backed provider, `aitui` does not receive the provider token:

```txt
user chooses OpenAI subscription
        |
        v
aitui starts the official login flow
        |
        v
the official CLI stores its own auth state
        |
        v
aitui invokes the official CLI for requests
```

This keeps subscription authentication with the provider's supported tooling
while still giving the end user a single setup flow inside the host CLI.

## Providers

Providers describe an AI backend and create clients for that backend. Most applications can start with `providers.Default()`, which returns the bundled providers that `aitui` knows how to offer.

The selector asks each provider whether it is available, presents the usable options, and returns a client for the selected provider. An API provider can call a remote service directly. A CLI-backed provider can execute an official command-line tool and rely on that tool for login and account management.

Applications that need tighter control can build their own provider list:

```go
client, err := aitui.SelectClient(ctx, []aitui.Provider{
	myProvider,
	anotherProvider,
})
```

Provider adapters should keep backend-specific details narrow so downstream CLIs can depend on `aitui` interfaces instead of provider internals.

## Authentication

Applications remain responsible for deciding whether user-supplied API keys are used only for the current command or remembered in application-owned storage. API keys can also come from environment variables, explicit configuration, or a custom transport.

CLI-backed subscription providers should use the official tool's own authentication flow. For example, an OpenAI subscription provider can guide the user through Codex CLI login instead of reading browser sessions or storing subscription tokens itself.

Credential material stays with the application or the official provider tools. provider adapters receive only what they need to answer a request.

## Testing

Code that depends on `aitui` can test against fake providers and fake clients instead of making live provider calls. Provider adapters should also use fake HTTP transports or fake command runners in their own unit tests.
