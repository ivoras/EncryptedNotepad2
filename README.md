# Encrypted Notepad 2

This is a spiritual successor to the "Encrypted Notepad" app [previously developed on SourceForge](https://sourceforge.net/projects/enotes/), by the same author. The goal is still the same: "Encrypted Notepad 2" does only one thing, but aims to do it perfectly - a Notepad-like simple text editor where files are saved (and later loaded) encrypted with industrial strength algorithms. No ads, no network connection required, no bloat, just run it. If you know how to use the ancient Notepad app, you know how to use this app.

![Screenshot](screenshot.png)

# Current status

In development. Major features are still missing.

# Downloading Encrypted Notepad 2

TODO

# Building Encrypted Notepad 2

You will need Go 1.22+ installed. After cloning the repo, just run:

```
go build
```


# Encryption

The files are encrypted using AES-256 and stored in the PGP/OpenPGP `.asc` format, that is interoperable with any other tool using the same standard.

# F.A.Q.

## Why is the executable / package so big, compared to the old version?

Encrypted Notepad 2 is written in Go, and that means it's mostly statically compiled on all platforms. It uses the [Fyne](https://github.com/fyne-io/fyne) UI toolkit, and that means it uses almost no operating system-provided UI facilities on any platform. The flip side of that is that everything needs to be built-in into the single executable, making it bigger than expected for such a compact app.

The old version was written in Java, and that means it required a JRE to run. In that light, the new version is actually lighter-weight.

## Will you support other encryption algorithms?

No. Really, there's no need to. Either you trust AES, and in that case this is what you want, or you don't, in which case you most likely don't need this tool.

## Will you support more file formats (other than OpenPGP ASCII-armoured messages)?

Maybe - depends if there's a good use case and enough people want it.

## What is the actual cipher mode of AES-256 used in Encrypted Notepad 2?

When saving in OpenPGP's message format (the `.asc`) files, the mode is dictated by the OpenPGP spec. It is [OCFB-MDC](https://web.archive.org/web/20231230093732/https://articles.59.ca/doku.php?id=pgpfan:mdc). It is quite a bit more sophisticated than a naive ECB of even a CBC approach.
