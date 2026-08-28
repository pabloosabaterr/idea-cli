# Terminal Notepad

This is a notepad that works on the terminal, thought to complete that need that
I have of having everything on the terminal.

Also I wanted to get into Go so... This seemed like a good idea.

How it works is very simple:

It has three main dependencies:
- [tea](https://github.com/charmbracelet/bubbletea)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [glamour](https://github.com/charmbracelet/glamour)

It will scan for `~/notes` and create the `notes` folder if it doesn't found it.
After that, at the left panel they will be shown the folders and depending on which
folder are you on, the middle pannel will show the `.md` files that are inside.

At the right side, it is a preview for the file you are hovering.

## Planed?

Commands, like vim, something like:

`:st` : will move to the starred folder the hovered file.
 
And more for creating folders, files, etc, etc.

Made with love Pablo.

