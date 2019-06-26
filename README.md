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

Download and extract the [latest
release](https://github.com/ff14wed/aetherometer/releases) to a local
directory on your system. Then run `aetherometer-ui.exe`.

## For Plugin Developers

Aetherometer exposes a good amount of information that identifies your
character through its [GraphQL API](core/models/schema.graphql). To prevent
the leaking of information to practically any web page on the internet,
Aetherometer core implements two (relatively weak but probably sufficient)
safety measures to make sure only authorized plugins can access the API:
  1. CORS validation in order to safe-guard against just random webpages from
     accessing the localhost endpoint.
  2. Validation of an API token via the `Authorization` header to permit only
     user-approved applications access (as opposed to clients that bypass
     CORS).

When adding a plugin, Aetherometer will load it as a web page and inject
credentials so that this plugin can access the API. That said, although
Aetherometer does basically run a full Chromium process per plugin,
for security reasons it will open all hyperlinks to external sites in an
external browser instead.

For integration with online services, opt to use a secure authorization
workflow like OAuth instead of logging the user in directly in the UI.

### API Token

In order to allow your plugin to access the API, it should expose a function
on the global `window` object called `initPlugin` that receives a single
object argument specified by the following:
  ```typescript
    interface PluginParams {
      apiURL: string; // URL that allows access to the API, e.g. http://localhost:8080/query
      apiToken: string; // API token in JWT format
      streamID: number; // The focused stream ID. Other stream IDs could be handled by other instances of the same plugin
    }
  ```

To support cases where the plugin is not launched via the Aetherometer UI,
the plugin may check for the absence of the `window.waitForInit` boolean
(loaded as part of Aetherometer's preload script) and use custom parameters
for loading the application.

See [examples](#plugin-examples) for examples on how to integrate this
workflow into the startup of your app.

### Getting an API Token for Testing or External Applications

The API token that is created for a single instance of a plugin isn't
restricted to only that plugin (and there's actually no secure way to
verify the origin of requests if all is on localhost anyways).

If your application is loaded in the web browser, you can load some plugin
that will display its provided `apiToken` (like the
[Playground](#playground)) and just copy and paste the token into
a secure place for the application to read it.

Or you could go one step further and have this plugin integrate with your
external application through another local or online web service.

### Playground

For testing of GraphQL queries against the API and documentation of the API,
you may use the Playground (located at `http://localhost:8080/playground`) as
a plugin. If the core server is not running at `8080`, then it is probably
running on some separate port. This port you can find by inspecting the
core log output (in the settings menu), or by looking at the `config.json`
used to run `core.exe`. This `config.json` can be found in the parent folder
of the `logs` directory containing core logs.

### Plugin Examples
[Inspector](https://github.com/ff14wed/inspector-plugin) - React-based App
Craftbot - Coming soon

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

xivhook.dll is the only part of the system that touches the game itself. It uses
a stable and  more reliable (not as lossy as using WinPcap) way of capturing network
data and does not require any memory accesses to the game itself. Therefore, it does not
need updating whenever the game executable updates.

Eventually, it can be used to enable overlays for the game with the same
technology that Discord uses to power its ingame overlay.

Unfortunately, this is the only part of Aetherometer that is not open-source
since it can be easily modified to do a lot more ToS-breaking stuff.
