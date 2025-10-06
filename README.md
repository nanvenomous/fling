# Fling - Fast Application Launcher

A lightweight TUI application launcher for Linux that integrates with i3wm. Uses fzf for fuzzy searching and can launch both PATH executables and .desktop applications.

## Features

- Fast fuzzy search through system applications
- Supports both PATH executables and .desktop files
- Integrates seamlessly with i3wm
- Lightweight and responsive TUI interface

## Dependencies

- `fzf` - for the fuzzy finder interface
- `bash` - for the application discovery script

## Installation

```bash
make build
sudo make install
```

Or manually:
```bash
go build -o fling
sudo cp fling /usr/local/bin/
sudo cp apps.sh /usr/local/bin/fling-apps.sh
```

## Usage

### Basic usage
```bash
fling query apps.sh
```

### i3wm Integration

Add this to your i3 config (`~/.config/i3/config`):

```
bindsym $mod+d exec --no-startup-id alacritty --class fling -e fling query /usr/local/bin/fling-apps.sh
for_window [class="fling"] floating enable, resize set 800 600, move position center
```

This will:
- Bind `Mod+d` to launch fling in a floating alacritty terminal
- Set the window to 800x600 pixels and center it on screen
- Use the installed apps.sh script

### Custom Application Scripts

You can create your own application discovery scripts. The script should output one application per line. For .desktop applications, use the format: `Name|Exec|DesktopFile`.

Example:
```bash
#!/bin/bash
echo "Firefox|firefox|/usr/share/applications/firefox.desktop"
echo "vim"
echo "htop"
```

## How it works

1. `apps.sh` scans the system PATH and .desktop files to find available applications
2. `fling` feeds this list to `fzf` for interactive selection
3. Selected applications are launched with proper process detachment

## Building from Source

```bash
git clone <repository>
cd fling
go mod tidy
make build
```