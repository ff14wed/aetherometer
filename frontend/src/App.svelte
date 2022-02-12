<script>
	import { onMount } from "svelte";

	let plugins = [];

	onMount(async () => {
		await window.go.main.App.WaitForStartup();
		plugins = Object.entries(await window.go.main.App.GetPlugins());
	});
</script>

<main>
	<div id="logo" />
	{#each plugins as [name, pluginInfo]}
		<div>{name} - {pluginInfo.PluginID} - {pluginInfo.PluginURL}</div>
	{/each}
</main>

<style>
	main {
		height: 100%;
		width: 100%;
	}

	#logo {
		width: 40%;
		height: 40%;
		padding-top: 20%;
		margin: auto;
		display: block;
		background-position: center;
		background-repeat: no-repeat;
		background-image: url("assets/images/logo-dark.svg");
	}
</style>
