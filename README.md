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

Download and extract the [latest release](https://github.com/ff14wed/aetherometer/releases) to a local directory on your system. Then run
`aetherometer-ui.exe`.

## For developers

### Building on Windows

To build core, simply run `go build -o resources/win/core.exe main.go`.

Also copy your distribution of `xivhook.dll` into this `resources/win`
directory.

Then to build the bundle including the UI, `cd` into the `ui` directory,
run `yarn install` and then `yarn run build`.

### Other Platforms

Currently, other platforms are not fully tested, but you could adapt the
steps for Windows and get pretty far in the process.

Core is written to be platform agnostic, but it only has an adapter (for
ingesting data) suited for running on Windows. With the right adapters,
possibly using pcap, it can be made to run on Mac OSX or Linux.
