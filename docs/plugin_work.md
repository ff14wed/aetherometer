# For Plugin Developers

Aetherometer supports both GraphQL Queries and GraphQL Subscriptions,
so you can develop real-time applications that interact with data from FFXIV.

Aetherometer exposes a good amount of information that identifies your
character through its [GraphQL API](core/models/schema.graphql).

To prevent the leaking of information to practically any web page on the
internet, Aetherometer core implements two (relatively weak but probably
sufficient) safety measures to make sure only authorized plugins can access
the API:
  1. CORS validation in order to safe-guard against just random webpages from
     accessing the localhost endpoint.
  2. Validation of an API token via the `Authorization` header to permit only
     user-approved applications access (as opposed to clients that bypass
     CORS).

When adding a plugin, Aetherometer will load it as a web page in a sandboxed
iframe and inject credentials so that this plugin can access the API.

For integration with online services, opt to use a secure authorization
workflow like OAuth instead of logging the user in directly in the UI.

## API Token

In order to allow your plugin to access the API, it should parse query params
like so:
  ```typescript
  const urlParams = new URLSearchParams(window.location.search);
  let apiURL = urlParams.get('apiURL'); // URL that allows access to the API, e.g. http://localhost:8080/query
  let apiToken = urlParams.get('apiToken') || undefined;
  let streamID = undefined;
  if (urlParams.has('streamID')) {
    streamID = parseInt(urlParams.get('streamID')!);
  }
  ```

See [examples](#plugin-examples) for examples on how to integrate this
workflow into the startup of your app.

## Getting an API Token for Testing or External Applications

You may provide a `local_token` in the config file in order to allow apps
hosted on your local machine access to the API with that token.

## Playground / API Documentation

For testing of GraphQL queries against the API or reading of documentation
for the API, you may use the Playground (located at
`http://{YOUR API URL}.replace('/query', '/playground'`) as a plugin. You can find
the API URL in the Aetherometer settings.

## Plugin Examples
 - [Inspector](https://github.com/ff14wed/inspector-plugin) - TypeScript + React-based App
 - [Craftbot](https://github.com/ff14wed/craftbot-plugin) - Based on same
   stuff as above, but has a much narrower scope on the API and is simpler.
 - [Aetherometer UI](../ui/src/api/gqlClient.ts) - Okay, not really a plugin, but
   you could also read this code to get a feel for how to query the API.
 - [Triggerbot]() - Coming soon
