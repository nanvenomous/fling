#!/bin/bash

# Find all executable files in PATH and .desktop files
find_executables() {
    # Get all directories in PATH
    IFS=':' read -ra PATHS <<< "$PATH"
    
    # Find executables in PATH
    for dir in "${PATHS[@]}"; do
        if [ -d "$dir" ]; then
            find "$dir" -maxdepth 1 -type f -executable -printf '%f\n' 2>/dev/null
        fi
    done
}

# Find .desktop files and extract application names
find_desktop_apps() {
    local desktop_dirs=(
        "/usr/share/applications"
        "/usr/local/share/applications"
        "$HOME/.local/share/applications"
        "/var/lib/flatpak/exports/share/applications"
        "$HOME/.local/share/flatpak/exports/share/applications"
    )
    
    for dir in "${desktop_dirs[@]}"; do
        if [ -d "$dir" ]; then
            find "$dir" -name "*.desktop" -exec grep -l "^Type=Application" {} \; 2>/dev/null | \
            while read -r file; do
                # Extract Name field from .desktop file
                name=$(grep "^Name=" "$file" | head -1 | cut -d'=' -f2-)
                exec_line=$(grep "^Exec=" "$file" | head -1 | cut -d'=' -f2-)
                if [ -n "$name" ] && [ -n "$exec_line" ]; then
                    echo "$name|$exec_line|$file"
                fi
            done
        fi
    done
}

case "$1" in
    "executables")
        find_executables | sort -u
        ;;
    "desktop")
        find_desktop_apps | sort -u
        ;;
    "all"|*)
        {
            find_executables
            find_desktop_apps | cut -d'|' -f1
        } | sort -u
        ;;
esac