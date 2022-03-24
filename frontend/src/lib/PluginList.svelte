<script lang="ts">
  import { Button, TextInput } from "carbon-components-svelte";
  import ComboBox from "./carbon/ComboBox.svelte";
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

  let presets = [
    {
      id: "preset-inspector",
      text: "Inspector",
      url: "https://ff14wed.github.io/inspector-plugin/",
    },
    {
      id: "preset-craftlog",
      text: "Craftlog",
      url: "https://ff14wed.github.io/craftlog-plugin/",
    },
  ];

  let comboBox;
  let selectedPresetID;

  function updatePresetURL(id: string) {
    for (let preset of presets) {
      if (preset.id === id) {
        newPluginURL = preset.url;
        return;
      }
    }
    newPluginURL = "";
  }

  $: updatePresetURL(selectedPresetID);
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
        <ComboBox
          placeholder="Enter plugin name..."
          direction="top"
          bind:this={comboBox}
          bind:value={newPluginName}
          bind:selectedId={selectedPresetID}
          items={presets}
          let:item
        >
          <div>
            <strong>{item.text}</strong>
          </div>
          <div>
            url: {item.url}
          </div>
        </ComboBox>
      </StructuredListCell>
      <StructuredListCell>
        <TextInput
          hideLabel
          labelText="Plugin URL"
          placeholder="eg. https://plugins.com/foo/"
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
            comboBox.clear({ focus: false });
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

  :global(.bx--list-box__menu-item) {
    height: auto;
  }

  :global(.bx--list-box__menu-item__option) {
    height: auto;
  }
</style>
