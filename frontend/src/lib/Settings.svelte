<script lang="ts">
	import { Button, Modal } from "carbon-components-svelte";

	import FolderOpen16 from "carbon-icons-svelte/lib/FolderOpen16";
	import { onMount } from "svelte";
	import PluginList from "./PluginList.svelte";

	export let open = false;

	let config = {
		Plugins: {},
	};

	let appVersion = "";
	let apiVersion = "";
	let apiURL = "";
	let appDirectory = "";

	onMount(async () => {
		await window.go.app.Bindings.WaitForStartup();

		config = (await window.go.app.Bindings.GetConfig()) || config;
		if (!config.Plugins) {
			config.Plugins = {};
		}

		appVersion = (await window.go.app.Bindings.GetVersion()) || "";
		apiVersion = (await window.go.app.Bindings.GetAPIVersion()) || "";
		apiURL = (await window.go.app.Bindings.GetAPIURL()) || "";
		appDirectory = (await window.go.app.Bindings.GetAppDirectory()) || "";

		window.runtime.EventsOn("ConfigChange", async () => {
			config = (await window.go.app.Bindings.GetConfig()) || config;
			if (!config.Plugins) {
				config.Plugins = {};
			}
			console.log("Updating config", config);
		});
	});

	async function addPlugin(name: string, url: string) {
		await window.go.app.Bindings.AddPlugin(name, url);
	}

	async function deletePlugin(name: string) {
		await window.go.app.Bindings.RemovePlugin(name);
	}

	async function openAppDirectory() {
		await window.go.app.Bindings.OpenAppDirectory();
	}

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
	<PluginList plugins={config.Plugins} {addPlugin} {deletePlugin} />
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
</style>
