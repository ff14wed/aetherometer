<script lang="ts">
  import { Button, TextInput } from "carbon-components-svelte";
  import {
    StructuredList,
    StructuredListHead,
    StructuredListRow,
    StructuredListCell,
    StructuredListBody,
  } from "carbon-components-svelte";

  import Add16 from "carbon-icons-svelte/lib/Add16";
  import Delete16 from "carbon-icons-svelte/lib/Delete16";

  let newPluginName;
  let newPluginURL;

  export let plugins;
  export let addPlugin;
  export let deletePlugin;
</script>

<StructuredList condensed>
  <StructuredListHead>
    <StructuredListRow head>
      <StructuredListCell head>
        <div class:plugin={true}>Name</div>
      </StructuredListCell>
      <StructuredListCell head>
        <div class:plugin={true}>URL</div>
      </StructuredListCell>
      <StructuredListCell head />
    </StructuredListRow>
  </StructuredListHead>
  <StructuredListBody>
    {#each Object.entries(plugins) as [name, url]}
      <StructuredListRow>
        <StructuredListCell>
          <div class:plugin={true}>{name}</div>
        </StructuredListCell>
        <StructuredListCell>
          <div class:plugin={true}>{url}</div>
        </StructuredListCell>
        <StructuredListCell>
          <Button
            iconDescription="Delete Plugin"
            icon={Delete16}
            size="small"
            kind="danger"
            on:click={() => deletePlugin(name)}
          />
        </StructuredListCell>
      </StructuredListRow>
    {/each}
    <StructuredListRow>
      <StructuredListCell>
        <TextInput
          hideLabel
          labelText="Plugin Name"
          placeholder="Enter plugin name..."
          size="sm"
          bind:value={newPluginName}
        />
      </StructuredListCell>
      <StructuredListCell>
        <TextInput
          hideLabel
          labelText="Plugin URL"
          placeholder="eg. https://plugins.com/foo/"
          size="sm"
          bind:value={newPluginURL}
        />
      </StructuredListCell>
      <StructuredListCell>
        <Button
          iconDescription="Add Plugin"
          icon={Add16}
          size="small"
          on:click={() => {
            addPlugin(newPluginName, newPluginURL);
            newPluginName = "";
            newPluginURL = "";
          }}
        />
      </StructuredListCell>
    </StructuredListRow>
  </StructuredListBody>
</StructuredList>

<style>
  .plugin {
    padding: 0 1rem;
  }

  :global(.bx--structured-list) {
    margin-bottom: 2rem;
  }

  :global(.bx--structured-list-td) {
    vertical-align: middle;
  }
</style>
