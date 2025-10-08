# fling - Terminal App Launcher for Your Linux Rice

I recently created a lightweight, speedy TUI app launcher.

What it does:

• Lightning-fast fuzzy search through all your executables
• Learns your usage patterns and bubbles favorites to the top
• Handles TUI apps intelligently
• Written in Go, so it's fast and a single binary
• Integrates beautifully with i3 & other WMs

Why I built it:
- I prefer terminal interfaces over GUI interfaces.
    - they inhert styling
    - are very fast to launch
- I wanted to be able to launch other TUI applications directly from the launcher itself
    - you can configure them to launch in their own terminals
- Figured `fzf` would take minimal effort to spin into an application launcher
    - if you love fzf's fuzzy search, you'll feel right at home
- Memory of recently launched apps
    - sorted by usage then alphabetically

Installation is dead simple:
```
go install github.com/nanvenomous/fling@latest
```
(other install options in the readme)

Works great as a dmenu replacement in i3.
I've got it bound to Mod+d and it pops up centered in a floating terminal.

Been using it daily for a few weeks now and figured others might find it useful.
It's not trying to reinvent the wheel, just be really good at one thing.

GitHub: https://github.com/nanvenomous/fling

Would love to hear what you think if you give it a try!
Could also use help getting integrations for other windown managers in the readme.
