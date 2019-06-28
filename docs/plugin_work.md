# For Plugin Developers

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

## API Token

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

## Getting an API Token for Testing or External Applications

The API token that is created for a single instance of a plugin isn't
restricted to only that plugin (and there's actually no secure way to
verify the origin of requests if all is on localhost anyways).

If your application is loaded in the web browser, you can load some plugin
that will display its provided `apiToken` (like the
[Playground](#playground)) and just copy and paste the token into
a secure place for the application to read it.

Or you could go one step further and have this plugin integrate with your
external application through another local or online web service.

## Playground / API Documentation

For testing of GraphQL queries against the API or reading of documentation
for the API, you may use the Playground (located at
`http://{api_url}.replace('/query', '/playground'`) as a plugin. You can find
the API URL in the `About` section of the Aetherometer settings.

## Plugin Examples
 - [Inspector](https://github.com/ff14wed/inspector-plugin) - React-based App
 - Craftbot - Coming soon
 - Triggerbot - Coming soon