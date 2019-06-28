# Aetherometer

Aetherometer is a framework that parses network data for FFXIV and presents
the parsed data stream through a GraphQL API that allows plugins to access
and display this information.

**Now updated for Shadowbringers**

<img src="docs/preview.png" alt="preview" />

## Getting Started
[**Download for the latest release**](https://github.com/ff14wed/aetherometer/releases)

Download the zip file to a place with sufficient space on your system, and
extract it. Then run `aetherometer-ui.exe`.

Then try out some [plugins](#plugins-list)!

If you're a developer interested in Aetherometer, see
[here](#for-developers).

## Features

Aetherometer is capable of supporting many different use cases, including
parsing combat data for trigger or DPS logging purposes, display of
player and enemy movement on the map, crafting progress, etc.

## Plugins

Plugins are mini-applications that process data from a specific instance
of the game. They open as new "tabs" on the navigation drawer to the left,
and you can navigate between them without losing data on any other tab.

 >Techincally plugins are able to process data from mulitple
instances of the game, but Aetherometer's Stream handling abilities aim
to reduce boilerplate Stream switching code.

You can also switch between instances of games via the Stream menu on the bottom left to switch to that session's "tabs".

### Plugins List

Here are some plugins that you can try to get an idea for Aetherometer's capabilities:
- Inspector Plugin - https://github.com/ff14wed/inspector-plugin

## Installing Plugins

<img src="docs/settings.png" alt="preview" />

1. Navigate to the Settings pane of the UI.
2. Go to the `Manage Plugins` section of the Settings pane.
3. Check the top-level tree nodes corresponding to the streams to which you
   would like to add instances of plugins.
    - **If the source of data is a running instance of FFXIV, then
    Stream means Process ID**. Other sources of data are also listed here.
    - The **Default Plugins** section lists plugins you would like
      automatically started with every new instance of a Stream (game). You
      can freely add or remove plugins here to affect *future* instances of
      Streams.
4. Click `Add Plugin`.
5. You can click `Unselect All` to reset all checkboxes.

### Removing Plugins / Closing Panes

If you want to close the page for a plugin, simply make sure to check the
leaf-level tree nodes and click `Remove Plugins`.

## Miscellaneous Settings

### Automatically switch to new session when a stream is created

This setting will automatically navigate you to a new Stream's session (kind
of like a window in a browser) whenever a new instance of a Stream is
started. Old Stream sessions that are closed will still be accessible via the
Stream menu on the left navigation drawer.

i.e, if you were viewing details about one instance of a game,
starting another instance would automatically create a new "window" for this
instance with new "tabs" (plugins) and navigate you to this new "window". In
order to go back to viewing details about older instances of the game, you
would have to navigate back to the other "windows" via the Stream menu.

This setting is on by default. Toggle the switch to disable this behavior.

### Stream sessions retained

Keeping too many Stream sessions alive could eat up your RAM and hamper
performance of your machine. To clean up, Aetherometer can automatically
close old and inactive sessions (where the game is no longer running).

However, your old inactive sessions could still
include data that needs saving. This session provides a buffer to keep alive
old "tabs" so you can finish working with them before letting them be closed.

The default setting is 1. Set a negative number to disable closing any
sessions.

Changing this setting does not instantly close existing inactive sessions
until another instance of a Stream is removed (i.e. when a game process
is no longer running).

**If you want to immediately close some tabs for performance reasons, simply
just remove all plugins from those streams, taking care not to delete your
defaults**.

## For Developers

### Creating plugins
See the [docs/plugin_work.md](docs/plugin_work.md) document.
### Contributing to Aetherometer
See the [CONTRIBUTING.md](CONTRIBUTING.md) document.
