# tele

SSH destination manager. Save SSH connections behind a single master password and connect by name.

All destination passwords are encrypted on disk using AES-256-GCM with a key derived from your master password via Argon2id. You remember one password — tele handles the rest.

## Install

```
go install .
```

Requires Go 1.23+.

## Usage

```
tele init          Set up your master password
tele add <name>    Save a new SSH destination
tele go <name>     Connect to a destination
tele list          List saved destinations
tele rm <name>     Remove a destination
```

### Set up

```
$ tele init
Enter master password:
Confirm master password:
Master password set successfully.
```

### Add a destination

```
$ tele add prod
Enter master password:
Host: 10.0.1.50
Port [22]: 2222
User: deploy
Password:
Destination "prod" added.
```

### Connect

```
$ tele go prod
Enter master password:
# opens SSH session to deploy@10.0.1.50:2222
```

### List destinations

```
$ tele list
  prod → deploy@10.0.1.50:2222
  staging → admin@10.0.1.51:22
```

### Remove a destination

```
$ tele rm staging
Destination "staging" removed.
```

## How it works

- `tele init` creates a master config with a random salt and an Argon2id hash of your password (for verification only).
- `tele add` encrypts the destination password with AES-256-GCM using a key derived from your master password + a per-destination random salt.
- `tele go` re-derives the key, decrypts the password, and execs into `sshpass + ssh`.
- `sshpass` is installed automatically on first use if not already on your PATH. It is compiled from source and stored in tele's config directory.

## Data storage

All data lives in your OS config directory:

| OS | Path |
|----|------|
| macOS | `~/Library/Application Support/tele/` |
| Linux | `~/.local/share/tele/` |
| Windows | `%APPDATA%\tele\` |

```
tele/
├── master.json              # salt + password hash
├── bin/
│   └── sshpass              # auto-installed binary
└── destinations/
    └── <name>.json          # host, port, user, encrypted password
```

No passwords are stored in plaintext.

## Dependencies

- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) — Argon2id key derivation
- [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) — terminal password input (no echo)
- A C compiler (Xcode CLI tools on macOS, gcc/clang on Linux) — only needed if sshpass isn't already installed

## License

MIT
