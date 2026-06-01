# Feature IN-03 — TLS Configuration

**Slug:** tls-configuration  
**Output file:** IN-03-tls-configuration.md

## Summary

The `IN-03` feature, **TLS Configuration**, provides shared utilities for establishing TLS-secured connections within the application. It specifically encapsulates the generation of self-signed certificates for local development environments, facilitating secure communication between client and server components during testing and development.

## Scope

- **In-scope:**
  - Generation of in-memory RSA private keys and self-signed X.509 certificates.
  - Configuration of `crypto/tls.Config` objects for use in local network services.
  - Hardcoded support for `localhost` and loopback IP addresses (`127.0.0.1`, `::1`).
- **Out-of-scope:**
  - Production-grade certificate management or authority (CA) integration.
  - Configuration of TLS protocols, ciphers, or certificate validation policies for non-local environments.

## Functional behaviour

The primary capability is `TLSDevConfig`, which programmatically generates a self-signed certificate, encodes it into PEM format, and returns a `*tls.Config` object populated with that certificate. This simplifies the setup of secure servers and clients in development without requiring external infrastructure for certificate management.

## Technical design

The feature is implemented in the `tls_config` package. It leverages Go's standard library `crypto` packages to perform the following steps:
1. Generate an RSA key pair (`tls_config/config.go:17`).
2. Define a certificate template with `Subject: "CLI Chat"`, valid for 365 days, and restricted to `localhost` and loopback IPs (`tls_config/config.go:22-31`).
3. Sign the certificate with the generated key (`tls_config/config.go:33`).
4. Construct a `tls.Certificate` from the PEM-encoded certificate and private key (`tls_config/config.go:38-41`).
5. Return a `tls.Config` initialized with this certificate (`tls_config/config.go:46`).

## Entry points

The `TLSDevConfig()` function is exported and consumed by the server and client components:
- **Server:** Used in `cmd/server/main.go:35` to initialize TLS configuration for the server process.
- **Client:** Used in `client/file-server.go:24` to configure TLS for file serving capabilities.

## Dependencies

- Go standard library:
  - `crypto/*` (for TLS, X.509, RSA, PEM handling)
  - `net` (for IP address definitions)
  - `fmt`, `math/big`, `time` (for utility operations)

## Revision

- Initial draft: scope, behaviour, and implementation details populated from `tls_config/config.go`, `cmd/server/main.go`, and `client/file-server.go`.

## Evidence index

- `tls_config/config.go:17-46`: Implementation of `TLSDevConfig` generating self-signed certificates and `tls.Config`.
- `cmd/server/main.go:35`: Usage of `tls_config.TLSDevConfig()` in server initialization.
- `client/file-server.go:24`: Usage of `tls_config.TLSDevConfig()` in client file server setup.

<!-- easyspecs-nav:children:start -->
## EasySpecs — Child documents

- [Generate TLS dev private key](./IN-03_IC-01-ic01.md)
- [Create TLS dev certificate](./IN-03_IC-02-ic02.md)
- [Construct TLS configuration](./IN-03_IC-03-ic03.md)
<!-- easyspecs-nav:children:end -->




<!-- easyspecs-nav:parents:start -->
## EasySpecs — Parent documents

- [EasySpecs project](./project.md)
<!-- easyspecs-nav:parents:end -->



