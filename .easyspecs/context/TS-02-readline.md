# Tool TS-02 — github.com/chzyer/readline

**Slug:** readline · **File:** TS-02-readline.md

## Summary

`github.com/chzyer/readline` is a Go library utilized in the CLI client to provide robust, interactive terminal input handling. It manages the input prompt, user input processing, and terminal display updates, ensuring a smooth user experience in the CLI environment.

## Role in this codebase

The application uses `readline` to create an interactive command-line interface. It serves as the primary mechanism for:
- Presenting a dynamic prompt to the user (`> ` by default, or server-driven prompts).
- Capturing user commands from the terminal.
- Managing terminal output to allow interleaving server notifications with user input.
- Handling terminal control sequences and refreshing the view when the screen content changes.

## Version and configuration

- **Version:** `v1.5.1` (per `go.mod` and `go.sum`)
- **Configuration:** Initialized in `client/client.go` with a custom `readline.Config`.

```go
rl, err := readline.NewEx(&readline.Config{
    Prompt: defaultPrompt,
})
```

## Boundaries

- **Application:** Implements the chat/file-transfer protocol, message processing, and business logic for the CLI. It uses `readline` to interface with the user.
- **Tool (`readline`):** Provides the low-level terminal UI capabilities (prompt rendering, input capture, history, control characters). The application is responsible for managing the `readline.Instance` lifecycle and handling the input string within the context of the application's protocol.

## Integration points

The `readline` instance is integrated into the `session` struct in `client/session.go`, facilitating thread-safe UI updates during concurrent server communications.

## Revision

- Initial draft: version and usage from manifests/config, usage patterns in `client/` package.
- Fix-up: Removed forbidden citation to `.easyspecs/context/architecture.md`.

## Evidence index

- `go.mod:5`: Dependency declaration for `github.com/chzyer/readline` v1.5.1.
- `client/client.go:8`: Import and usage of `readline` in `Connect()` to initialize the instance.
- `client/client.go:30-32`: Initialization of `readline.Instance` with `readline.NewEx`.
- `client/session.go:14`: Import of `readline` package.
- `client/session.go:24`: `readline.Instance` usage in `session` struct.
- `client/session.go:195`: Input reading via `s.rl.Readline()`.
- `client/session.go:185-190`: Controlling terminal input buffer via `s.rl.WriteStdin` and `s.rl.Refresh`.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

*None.*
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



