# envchain

> Utility for chaining and validating environment variable dependencies across service configs before deployment

---

## Installation

```bash
go install github.com/yourname/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envchain.git && cd envchain && go build ./...
```

---

## Usage

Define your environment variable dependencies in a `.envchain.yaml` file:

```yaml
chains:
  - name: database
    requires:
      - DB_HOST
      - DB_PORT
      - DB_PASSWORD
    depends_on:
      - APP_ENV
```

Then validate before deployment:

```bash
envchain validate --config .envchain.yaml
```

Example output:

```
✔ APP_ENV         resolved
✔ DB_HOST         resolved
✔ DB_PORT         resolved
✗ DB_PASSWORD     missing — required by chain: database

validation failed: 1 unresolved variable(s)
```

Run `envchain --help` to see all available commands.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)