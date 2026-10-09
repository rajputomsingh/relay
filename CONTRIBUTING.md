# Contributing to Relay

Thank you for your interest in contributing to Relay!

## Development Setup

Prerequisites:
- Go version compatible with `go.mod`
- Git
- MongoDB connection string for database-related development

Clone your fork and create a branch:

```bash
git clone https://github.com/YOUR_USERNAME/relay.git
cd relay
git switch -c feat/your-change
go mod download
```

Configure your local environment using `.env.example`. Never commit `.env`, API keys, database credentials, or other secrets.

## Code Guidelines

- Follow standard Go conventions.
- Format Go code using `gofmt`.
- Handle errors explicitly.
- Keep pull requests focused.
- Add or update tests when changing behavior.
- Update documentation when necessary.

## Validation

Run these commands before submitting your changes:

```bash
gofmt -w .
go test ./...
go build ./...
git diff --check
```

## Commit Messages

Use clear, descriptive commit messages:

```text
feat: add event filtering
fix: handle invalid event payloads
test: add ingestion validation tests
docs: improve project documentation
refactor: simplify request validation
```

## Pull Requests

1. Explain the problem and your proposed solution.
2. Describe the changes made.
3. Include testing steps and results.
4. Link related issues when applicable.
5. Mention breaking changes or configuration updates.
6. Ensure secrets and private data are excluded.

## Bug Reports

Include reproduction steps, expected and actual behavior, and sanitized error logs. Remove credentials and personal information.

## Security

Do not publish exploitable vulnerabilities or secrets in public issues. Contact the repository maintainer privately about security concerns.

## License

By contributing, you agree that your contributions are distributed under the project's Apache License 2.0, subject to its terms. See [LICENSE](LICENSE).
