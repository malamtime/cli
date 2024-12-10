# ShellTime CLI [![codecov](https://codecov.io/gh/malamtime/cli/graph/badge.svg?token=N09WIJHNI2)](https://codecov.io/gh/malamtime/cli)

The foundation CLI tool for shelltime.xyz - a platform for tracking DevOps work.

AnnatarHe: [![shelltime](https://api.shelltime.xyz/badge/AnnatarHe/count)](https://shelltime.xyz/users/AnnatarHe)

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/malamtime/installation/master/install.bash | bash
```

## Commands

### Authentication

```bash
shelltime auth [--token <your-token>]
```

Initializes the CLI with your shelltime.xyz authentication token. This command needs to be run before using other features.

Options:
- `--token`: Your personal access token from shelltime.xyz. if omit, you can also redirect to website to auth

Example:
```bash
shelltime auth --token abc123xyz
```

### Track

```bash
shelltime track [options]
```

Tracks your shells activities and sends them to shelltime.xyz.

Options:
- TODO: List track command options

Example:
```bash
shelltime track # TODO: Add example
```

### GC (Garbage Collection)

```bash
shelltime gc [options]
```

Performs cleanup of old tracking data and temporary files.

Options:
- TODO: List GC command options

Example:
```bash
shelltime gc # TODO: Add example
```

## Configuration

The CLI stores its configuration in `$HOME/.shelltime/config.toml`.

## Version Information

Use `shelltime --version` or `shelltime -v` to display the current version of the CLI.

## Support

For support, please contact: annatar.he+shelltime.xyz@gmail.com

## License

Copyright (c) 2024 shelltime.xyz Team
