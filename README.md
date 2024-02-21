# Encrypted Notepad 2

This is a spiritual successor to the "Encrypted Notepad" app [previously developed on SourceForge](https://sourceforge.net/projects/enotes/), by the same author. The goal is still the same: "Encrypted Notepad 2" does only one thing, but aims to do it perfectly - a Notepad-like simple text editor where files are saved (and later loaded) encrypted with industrial strength algorithms. No ads, no network connection required, no bloat, just run it. If you know how to use the ancient Notepad app, you know how to use this app.

# Current status

In development. Major features are still missing.

# Downloading Encrypted Notepad 2

TODO

# Building Encrypted Notepad 2

You will need Go 1.22+ installed. After cloning the repo, just run:

```
go build
```


# F.A.Q.

## Why is the executable / package so big, compared to the old version?

Encrypted Notepad 2 is written in Go, and that means it's mostly statically compiled on all platforms. It uses the [Fyne](https://github.com/fyne-io/fyne) UI toolkit, and that means it uses almost no operating system-provided UI facilities on any platform. The flip side of that is that everything needs to be built-in into the single executable, making it bigger than expected for such a compact app.

The old version was written in Java, and that means it required a JRE to run. In that light, the new version is much lighter-weight, at least on desktop platforms.
