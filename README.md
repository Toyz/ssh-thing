# SSH-Thing

A terminal-based SSH connection manager and log viewer with multi-server support.

## Features

- Monitor multiple SSH servers simultaneously in tabs
- Real-time log streaming
- Text colorization for common log formats
- Word-wrapping for long lines
- Keyboard navigation
- Mouse support
- Configurable key bindings

## Installation

### From Source

1. Clone the repository
   ```bash
   git clone https://github.com/toyz/ssh-thing.git
   cd ssh-thing
   ```

2. Build the application
   ```bash
   make build
   ```
   
3. Or build for your specific platform
   ```bash
   # For Linux
   make linux
   
   # For MacOS
   make darwin
   
   # For Windows
   make windows
   ```

## Configuration

### Server Configuration

Create a `servers.toml` file in the same directory as the executable (you can copy from `servers.example.toml`):

```toml
[[servers]]
name = "My Server"
host = "example.com"
user = "username"
private_key_path = "~/.ssh/id_rsa"
commands = ["tail -f /var/log/syslog"]
# Port will default to 22 if not specified

[[servers]]
name = "Docker Monitor"
host = "192.168.1.100"
user = "admin"
private_key_path = "~/.ssh/id_ed25519"
commands = ["docker logs -f --tail 10 container_name"]

[[servers]]
name = "Custom Port Server"
host = "example.org"
user = "admin"
port = 2222  # Custom SSH port
private_key_path = "~/.ssh/custom_key"
commands = ["journalctl -f"]
```

### Configuration Options

| Option | Description |
|--------|-------------|
| `name` | Display name for the server tab |
| `host` | Hostname or IP address |
| `user` | SSH username |
| `port` | SSH port (defaults to 22) |
| `private_key_path` | Path to SSH private key (supports ~ expansion) |
| `commands` | Array of commands to run after connecting |

## Keyboard Shortcuts

| Key | Description |
|-----|-------------|
| `←/h` | Previous tab |
| `→/l` | Next tab |
| `↑/k` | Scroll up |
| `↓/j` | Scroll down |
| `PgUp` | Page up |
| `PgDown` | Page down |
| `Home` | Scroll to top |
| `End/G` | Scroll to bottom |
| `tab` | Next tab |
| `shift+tab` | Previous tab |
| `c` | Toggle colorization |
| `w` | Toggle word wrap |
| `r` | Reset scroll position |
| `ctrl+l` | Clear buffer |
| `q/ctrl+c` | Quit |
| `?` | Toggle help screen |

## Customizing Key Bindings

The application will create a default keybinds.toml file in your config directory on first run.
You can edit this file to customize the key bindings.

## License

[Your License Here]
