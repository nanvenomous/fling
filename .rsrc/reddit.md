# fling - Terminal App Launcher for Your Linux Rice

I recently created a lightweight terminal based app launcher.

Some cool features:
- built on top of `fzf` **fuzzy search**; if you enjoy fzf, you'll feel at home
- Learns your **usage patterns** and bubbles favorites to the top
- terminal interface **inherts styling** & is quite fast to launch
- can **launch TUI applications directly** from the launcher itself
- **Integrates beautifully** with i3 & other WMs
- Written in Go, so it's fast and a **single binary**

Installation is simple:
```
go install github.com/nanvenomous/fling@latest
```
(other install options in the readme)

Works great as a dmenu replacement in i3.
I've got it bound to Mod+d and it pops up centered in a floating terminal.

Been using it daily for a few weeks now and figured others might find it useful.
Also, could use help getting integrations for other windown managers and terminal emulators in the readme.

GitHub: https://github.com/nanvenomous/fling

Would love to hear what you think if you give it a try!
Hope to see `fling` in your next rice ;)
