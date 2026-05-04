# aitui

Package `aitui` provides provider selection and chat client interfaces for Go
command-line applications.

`aitui` helps CLI tools offer AI features without hard-coding one backend into
the application. API-based providers can use caller-owned configuration, while
CLI-backed providers can delegate to official tools such as Codex CLI and use
their supported authentication flows. The main features are:

- It presents a stable interface for sending chat requests from Go CLIs.
- It lets applications register multiple providers and select one at runtime.
- It keeps provider-specific behavior behind small adapters.
- It supports fake providers and clients for tests.

* * *

- [Install](#install)
- [Examples](#examples)
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

Let's start with a CLI that offers one OpenAI API provider and sends a chat
request:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mvallenet/aitui"
	"github.com/mvallenet/aitui/providers/openaiapi"
)

func main() {
	ctx := context.Background()

	client, err := aitui.SelectClient(ctx, []aitui.Provider{
		openaiapi.New(openaiapi.Config{
			APIKey: os.Getenv("OPENAI_API_KEY"),
			Model:  os.Getenv("OPENAI_MODEL"),
		}),
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Chat(ctx, aitui.ChatRequest{
		Input: "Write a short release note for v0.1.0.",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Output)
}
```

The application owns the command, prompt, output handling, and provider list.
`aitui` handles the common path of selecting a backend and calling it through
one chat interface.

More than one provider can be offered when a CLI wants the user to choose:

```go
client, err := aitui.SelectClient(ctx, []aitui.Provider{
	openaiapi.New(openaiapi.Config{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  os.Getenv("OPENAI_MODEL"),
	}),
	openaicodex.New(openaicodex.Config{}),
})
```

If only one provider is available, it can be selected automatically. If
multiple providers are available, the configured selector UI asks the user
which backend to use.

## Providers

Providers describe an AI backend and create clients for that backend. An API
provider can call a remote service directly. A CLI-backed provider can execute
an official command-line tool and rely on that tool for login and account
management.

Provider adapters should keep backend-specific details narrow so downstream
CLIs can depend on `aitui` interfaces instead of provider internals.

## Authentication

Applications remain responsible for deciding how credentials are supplied.
API keys can come from environment variables, explicit configuration, or a
custom transport. CLI-backed providers should use the official tool's own
authentication flow.

Credential material stays with the application or the official provider tools;
provider adapters receive only what they need to answer a request.

## Testing

Code that depends on `aitui` can test against fake providers and fake clients
instead of making live provider calls. Provider adapters should also use fake
HTTP transports or fake command runners in their own unit tests.
