/*
Package config describes configuration for all of Aetherometer.  Aetherometer
configuration is written in the TOML format. https://toml.io/en/

Example

The following shows an example TOML config file for Aetherometer.

	api_port = 0
	disable_auth = false
	local_token = "foobar-token"

	[sources]
		data_path = "C:\\path\\to\\aetherometer\\resources\\datasheets"

		[sources.maps]
			cache = "C:\\path\\to\\aetherometer\\resources\\maps"
			api_path = ""

	[adapters]
		[adapters.hook]
			enabled = true
			dll_path = "C:\\path\\to\\aetherometer\\resources\\win\\deucalion.dll"
			ffxiv_process = "ffxiv_dx11.exe"

	[plugins]
	"Inspector" = "https://ff14wed.github.io/inspector-plugin/"
	"Craftlog" = "https://ff14wed.github.io/craftlog-plugin/"

Config File vs UI Settings

Some of this configuration can also be set in the Aetherometer UI rather than
changing it in the file. All changes are synchronized between the application
and the file, so there is no issue of out-of-sync configuration. However, some
configuration requires a restart of Aetherometer to take effect.

Also note that any comments written in the configuration file will be LOST if
you change any configuration in the UI, so it is wise to backup the config file
if you really need to.

API Port

The field `api_port` configures the port on which the GraphQL API server will
listen. For example, if this value is set to 8080, the server will be queryable
on http://localhost:8080/query

Requires a restart of Aetherometer for any changes to this field to take effect.

Disable Auth

Setting `disable_auth` to `true` allows plugins to query the API server without
an auth token. CORS validation is still enforced, so web-based plugins are still
rejected unless they originate from localhost. Intended for development purposes
only.  DISABLE AT YOUR OWN RISK.

Token for local plugins

The `local_token` field allows local plugins to specify this token in the
Authorization header for the API in order to gain access. This token may not be
used from a remote origin.

Auto Update

Setting `auto_update` to `true` allows Aetherometer to download required
resources if they are outdated or if they do not exist.

Sources Table

This table is primarily concerned with the configuration of data and map image
sources.

Data Path

The field `sources.data_path` configures the location of the CSV files
containing FFXIV data.  This directory must exist if the TOML config file is
provided.

Map Cache

The field `sources.map.cache` configures the location of the map cache.  This
directory must exist if the TOML config file is provided.

Map API Path

The field `sources.map.api_path` configures where to pull map images from if
they do not exist locally.  Defaults to "https://xivapi.com"

Adapters Table

This table lists configuration of the various ingress adapters that Aetherometer
supports. Currently, only the "hook" adapter for Windows is supported.

Hook Adapter

The table `[adapters.hook]` contains configuration for the "hook" adapter
(Windows only).  This adapter will automatically inject a hook into each FFXIV
process and read networked data to your Aetherometer instance.

Setting the field `adapters.hook.enabled` to `true` enables the adapter.

The field `adapters.hook.dll_path` specifies the hook DLL to inject into
FFXIV processes.

The field `adapters.hook.ffxiv_process` should be set to the name of the FFXIV
process into which to inject the hook. Generally it should be set to
"ffxiv_dx11.exe", but change it "ffxiv.exe" if you are using DirectX 9.

Plugins Table

This table contains a map of plugins, where the key is the display name of
the plugin and the value is the URL of the plugin.

If the plugin is a webpage-based plugin, it must be provided in this list in
order to be authorized to access the Aetherometer API.

	"My Plugin" = "https://foo.com/my/plugin"
	"Other Plugin" = "https://bar.com/other/plugin"

In the above example, "My Plugin" and "Other Plugin" will be the display names
of the two plugins added. Note that the scheme ("https://" part of the URL) is
required.

*/
package config
