# Building

## Building on Windows

**Requires Golang version at least 1.11 and latest stable version of Node.js.**

The Aetherometer UI depends on the binaries inside the `resources/win`
directory to function.

To build core, simply `cd` into the `core` directory and run
`go build -o ../resources/win/core.exe main.go`.

Copy your distribution of `xivhook.dll` (you can just copy this one from
the release version) into this `resources/win` directory.

Then to build the bundle including the UI, `cd` into the `ui` directory,
run `yarn install` and then `yarn run build`.

The built bundle will be located as a zip file in the `dist` directory.

## Other Platforms

Currently, other platforms are not fully tested, but you could adapt the
steps for Windows and get pretty far in the process.

Core is written to be platform agnostic, but it currently only has an adapter
(for ingesting data) suited for running on Windows. With the right adapters,
possibly using pcap, it can be made to run on Mac OSX or Linux.

## What is xivhook.dll?

xivhook.dll is the only part of the system that touches the game itself. It uses
a stable and  more reliable (not as lossy as using WinPcap) way of capturing network
data and does not require any memory accesses to the game itself. Therefore, it does not
need updating whenever the game executable updates.

Eventually, it can be used to enable overlays for the game with the same
technology that Discord uses to power its ingame overlay.

Unfortunately, this is the only part of Aetherometer that is not open-source
since it can be easily modified to do a lot more ToS-breaking stuff.