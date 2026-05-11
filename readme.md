# fling - a terminal based application launcher

A lightweight TUI application launcher for Linux.
Inspired by the power of `fzf` and the likeness of `dmenu` and `rofi`.

![fling preview](.rsrc/fling.png)

## Features

- Fast fuzzy search through system applications
- Usage history to put your favorite apps on top
- Hardcoded configuration for identifying, and launching TUI applications
- Scans all PATH directories for executables
- Integrates seamlessly with window managers (i3wm)
- Pure Go implementation (uses `fzf` [as a library](https://junegunn.github.io/fzf/tips/using-fzf-in-your-program/))

## Installation

- automatic from source:
    ```bash
    go install github.com/nanvenomous/fling@latest
    ```
- binary download
    ```bash
    curl -L https://github.com/nanvenomous/fling/releases/latest/download/fling-linux-amd64 > fling
    chmod +x fling
    sudo cp fling /usr/local/bin/
    ```
- from source with [task](https://github.com/go-task/task)
    ```bash
    git clone https://github.com/nanvenomous/fling.git
    cd fling
    task install
    ```
- manually from source:
    ```bash
    git clone https://github.com/nanvenomous/fling.git
    cd fling
    go build -o fling
    sudo cp fling /usr/local/bin/
    ```

## Usage

### Basic usage
```bash
fling
```

# Integration

### i3wm

You should be able to integrate `fling` with most window managers and terminal emulators.

Here is an example with [alacritty](https://github.com/alacritty/alacritty) and [i3](https://i3wm.org/):

Add this to your i3 config (`~/.config/i3/config`):
```
bindsym $mod+d exec --no-startup-id alacritty --class fling -e fling
for_window [class="fling"] floating enable, resize set 800 600, move position center
```

This will:
- Bind `Mod+d` to launch fling in a floating alacritty terminal
- Set the window to 800x600 pixels and center it on screen

### wayland swaywm

```
bindsym $mod+f exec footclient --app-id=modal fling
for_window [app_id="modal"] floating enable, resize set 800 600, move position center
```

## Configuration
a config file will be generated in `~/.config/fling/` the first time you run `fling`

see the [example config](./config.example.yml)
