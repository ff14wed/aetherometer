# Aetherometer

Aetherometer is a framework that parses network data for FFXIV and presents
the parsed data stream through a GraphQL API that allows plugins to access
and display this information.

Aetherometer is capable of supporting many different use cases, including
parsing combat data for trigger or DPS logging purposes, display of
player and enemy movement on the map, crafting progress, etc.

A full description of what Aetherometer currently exposes as part of the API
is available [here](core/models/schema.graphql).

Here is one example of a plugin that leverages Aetherometer to display
this information: https://github.com/ff14wed/inspector-plugin.

## Running

Download and extract the [latest release](https://github.com/ff14wed/aetherometer/releases)
to a local directory on your system. Then run `aetherometer-ui.exe`.

## For developers

### Building on Windows

Requires Golang version at least 1.11 and latest stable version of Node.js.

To build core, simply `cd` into the `core` directory and run
`go build -o ../resources/win/core.exe main.go`.

Also copy your distribution of `xivhook.dll` into this `resources/win`
directory.

Then to build the bundle including the UI, `cd` into the `ui` directory,
run `yarn install` and then `yarn run build`.

### Other Platforms

Currently, other platforms are not fully tested, but you could adapt the
steps for Windows and get pretty far in the process.

Core is written to be platform agnostic, but it currently only has an adapter
(for ingesting data) suited for running on Windows. With the right adapters,
possibly using pcap, it can be made to run on Mac OSX or Linux.

### What is xivhook.dll?

xivhook.dll is the only part of the system that touches the game itself. It
does not do any complicated memory accessing, but it uses a more reliable
(not as lossy as using WinPcap) way of capturing network data before it
reaches the game.

Eventually, it can be used to enable overlays for the game with the same
technology that Discord uses to power its ingame overlay.

Unfortunately, this is the only part of Aetherometer that is not open-source
since it can be easily modified to do a lot more ToS-breaking stuff.
