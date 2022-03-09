<script lang="ts">
	import {
		Button,
		Form,
		FormGroup,
		Modal,
		TextInput,
	} from "carbon-components-svelte";
	import {
		StructuredList,
		StructuredListHead,
		StructuredListRow,
		StructuredListCell,
		StructuredListBody,
	} from "carbon-components-svelte";

	import Add16 from "carbon-icons-svelte/lib/Add16";
	import Delete16 from "carbon-icons-svelte/lib/Delete16";
	import FolderOpen16 from "carbon-icons-svelte/lib/FolderOpen16";
	import { onMount } from "svelte";

	export let open = false;

	let config = {
		Plugins: {},
	};

	let appVersion = "";
	let apiVersion = "";
	let apiURL = "";
	let appDirectory = "";

	onMount(async () => {
		await window.go.main.App.WaitForStartup();

		config = (await window.go.main.App.GetConfig()) || config;
		if (!config.Plugins) {
			config.Plugins = {};
		}

		appVersion = (await window.go.main.App.GetVersion()) || "";
		apiVersion = (await window.go.main.App.GetAPIVersion()) || "";
		apiURL = (await window.go.main.App.GetAPIURL()) || "";
		appDirectory = (await window.go.main.App.GetAppDirectory()) || "";

		window.runtime.EventsOn("ConfigChange", async () => {
			config = (await window.go.main.App.GetConfig()) || config;
			if (!config.Plugins) {
				config.Plugins = {};
			}
			console.log("Updating config", config);
		});
	});

	async function addPlugin() {
		await window.go.main.App.AddPlugin(newPluginName, newPluginURL);
		newPluginName = "";
		newPluginURL = "";
	}
	async function deletePlugin(name: string) {
		await window.go.main.App.RemovePlugin(name);
	}

	async function openAppDirectory() {
		await window.go.main.App.OpenAppDirectory();
	}

	let newPluginName;
	let newPluginURL;

	$: infoTable = {
		"App version": appVersion,
		"API version": apiVersion,
		"API URL": apiURL,
	};
</script>

<Modal
	preventCloseOnClickOutside
	passiveModal
	bind:open
	modalHeading="Settings"
	on:open
	on:close
	size="lg"
>
	<h4>Aetherometer Info</h4>
	<div class:app-info={true}>
		{#each Object.entries(infoTable) as [name, info]}
			<div style="font-weight: bold" class:info-cell={true}>{name}:</div>
			<div class:info-cell={true}>{info}</div>
		{/each}
		<div style="font-weight: bold" class:info-cell={true}>
			App/log directory:
		</div>
		<div class:info-cell={true}>
			{appDirectory}
			<Button
				iconDescription="Open Directory"
				icon={FolderOpen16}
				size="small"
				on:click={openAppDirectory}
			/>
		</div>
	</div>

	<h4>Loaded Plugins</h4>
	<br />
	<Form>
		<FormGroup>
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
					{#each Object.entries(config.Plugins) as [name, url]}
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
								on:click={addPlugin}
							/>
						</StructuredListCell>
					</StructuredListRow>
				</StructuredListBody>
			</StructuredList>
		</FormGroup>
	</Form>
	<p>Settings are automatically saved.</p>
</Modal>

<style>
	.app-info {
		margin: 2rem;
		display: grid;
		grid-column-gap: 20px;
		grid-row-gap: 5px;
		grid-template-columns: max-content auto;
		grid-auto-rows: minmax(18px, auto);
	}

	.info-cell {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.plugin {
		padding: 0 1rem;
	}

	:global(.bx--structured-list-td) {
		vertical-align: middle;
	}
</style>
