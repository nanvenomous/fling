# fling - a terminal based application launcher

A lightweight TUI application launcher for Linux.
Inspired by the power of `fzf` and the likeness of `dmenu` and `rofi`.

![fling preview](.rsrc/fling.png)

## Features

- Fast fuzzy search through system applications
- Usage history to put your favorite apps on top
- Scans all PATH directories for executables
- Integrates seamlessly with i3wm
- Pure Go implementation

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

### i3wm Integration

Add this to your i3 config (`~/.config/i3/config`):

```
bindsym $mod+d exec --no-startup-id alacritty --class fling -e fling
for_window [class="fling"] floating enable, resize set 800 600, move position center
```

This will:
- Bind `Mod+d` to launch fling in a floating alacritty terminal
- Set the window to 800x600 pixels and center it on screen

