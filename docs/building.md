# Building

## Building on Windows

### Requirements
The Aetherometer UI depends on the binaries inside the `resources/win`
directory to function.

For Windows builds,
**make sure to copy `xivhook.dll` into this `resources/win` directory.**
(You can obtain a copy of this file from the `resources/bin` folder in the release version of Aetherometer.)

**Requires:**
 - Golang at least go1.11 (recommended go1.13): https://golang.org/dl/
 - Node.js: https://nodejs.org/en/
 - Yarn: https://yarnpkg.com/en/

### Taskfile Instructions

> Building with the Taskfile.yml assumes that you have go1.13 installed.
  If this is not the case, you may have to remove `-trimpath` from the
  `CORE_BUILD_FLAGS)

1. Install [Task](https://taskfile.dev/#/installation), which is a simple
   Make alternative. Any installation method is fine, but installing via
   `go get` into your `$GOPATH` is easiest since you should already have
   Golang set up.
2. Run `task build`.
3. The output bundle can be found as a zip file in the `dist` directory.

For more commands, see the [Taskfile.yml](../Taskfile.yml) at the root of the
project.

### Alternative Instructions
If not using the Taskfile.yml, building is still very easy.

To build core, simply `cd` into the `core` directory and run
`go build -o ../resources/win/core.exe main.go`.

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

xivhook.dll is the only part of the system that touches the game itself. It
uses a stable and more reliable (not as lossy as using WinPcap) way of
capturing network data and does not require any memory accesses to the game
itself. Therefore, it does not need updating whenever the game executable
updates.

Eventually, it can be used to enable overlays for the game with the same
technology that Discord uses to power its ingame overlay.

Unfortunately, this is the only part of Aetherometer that is not open-source
since it can be easily modified to do a lot more ToS-breaking stuff.