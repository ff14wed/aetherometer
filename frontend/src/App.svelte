<script lang="ts">
	import Shell from "./lib/Shell.svelte";
	import Tabs from "./lib/Tabs.svelte";
	import Tab from "./lib/Tab.svelte";
	import TabContent from "./lib/TabContent.svelte";
	import { selectedTabID } from "./lib/stores/stores";

	import { Content, ToastNotification } from "carbon-components-svelte";

	import { onMount } from "svelte";

	import type { runtime } from "../wailsjs/runtime";
	import go from "../wailsjs/go/bindings";
	import type { StreamInfo } from "../wailsjs/go/models";

	function sleep(ms) {
		return new Promise((resolve) => setTimeout(resolve, ms));
	}
	interface PluginInfo {
		PluginID: string;
		PluginURL: string;
	}

	let registeredPlugins: { [name: string]: PluginInfo } = {};

	function generatePluginList(
		streams: StreamInfo[],
		pluginObj: { [name: string]: PluginInfo }
	) {
		let plugins = [];
		for (const stream of streams) {
			for (const [name, pluginInfo] of Object.entries(pluginObj)) {
				plugins.push({
					name: `${name} - ${stream.name}`,
					id: `${pluginInfo.PluginID}-${stream.id}`,
					url: pluginInfo.PluginURL,
				});
			}
		}
		return plugins;
	}

	let activeStreams: StreamInfo[] = [];

	onMount(async () => {
		await window.go.main.App.WaitForStartup();

		// Load initial streams (though normally there aren't any)
		activeStreams = (await window.go.main.App.GetStreams()) || [];
		registeredPlugins = (await window.go.main.App.GetPlugins()) || {};

		console.log("Active streams", activeStreams);
		console.log("Registered Plugins", registeredPlugins);

		window.runtime.EventsOn("StreamChange", async () => {
			activeStreams = (await window.go.main.App.GetStreams()) || [];
			console.log("Updating active streams", activeStreams);
		});

		window.runtime.EventsOn("ConfigChange", async () => {
			registeredPlugins = (await window.go.main.App.GetPlugins()) || {};
			console.log("Updating registered Plugins", registeredPlugins);
		});
	});

	$: plugins = generatePluginList(activeStreams, registeredPlugins);
</script>

<Shell company="XIV" platformName="Aetherometer">
	<Tabs autoWidth bind:selectedTabID={$selectedTabID}>
		{#each plugins as plugin (plugin.id)}
			<Tab label={plugin.name} id={plugin.id} />
		{/each}
	</Tabs>
</Shell>
<section>
	<div class:spacer={true} />
	<div class:content={true}>
		{#if activeStreams.length === 0}
			<div class:padding={true}>
				<ToastNotification
					lowContrast
					hideCloseButton
					kind="warning-alt"
					title="No FFXIV processes detected."
					subtitle="Please launch the game and/or change zones."
					caption="If streams are still not detected, please check the application log (can be found in the settings page)."
				/>
			</div>
		{:else if plugins.length === 0}
			<div class:padding={true}>
				<ToastNotification
					lowContrast
					hideCloseButton
					kind="warning-alt"
					title="No plugins registered."
					subtitle="Please go to the settings page and add plugins."
				/>
			</div>
		{:else}
			{#each plugins as plugin (plugin.id)}
				<TabContent id={plugin.id} label={plugin.name}>
					<iframe
						sandbox="allow-same-origin allow-scripts allow-downloads"
						class:iframe={true}
						title={plugin.name}
						src={plugin.url}
					/>
				</TabContent>
			{/each}
		{/if}
	</div>
</section>

<style>
	section {
		display: flex;
		flex-flow: column;
		height: 100vh;
	}

	.spacer {
		height: 3rem;
		width: 100%;
		flex: 0 0 auto;
		clear: both;
	}

	.content {
		flex: 1;
		margin: 2px;
		user-select: none;
	}

	.padding {
		padding: 2rem;
	}

	.iframe {
		height: 100%;
		width: 100%;
	}
</style>
