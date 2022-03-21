<script lang="ts">
	import "carbon-components-svelte/css/g100.css";

	import { Header, HeaderUtilities } from "carbon-components-svelte";

	import Renew20 from "carbon-icons-svelte/lib/Renew20";
	import SettingsAdjust20 from "carbon-icons-svelte/lib/SettingsAdjust20";
	import Subtract24 from "carbon-icons-svelte/lib/Subtract24";
	import Stop20 from "carbon-icons-svelte/lib/Stop20";
	import Copy16 from "carbon-icons-svelte/lib/Copy16";
	import Close24 from "carbon-icons-svelte/lib/Close24";
	import HeaderButton from "./HeaderButton.svelte";
	import Settings from "./Settings.svelte";

	export let refreshCurrentTab;

	function minimize() {
		window.runtime.WindowMinimise();
	}

	function toggleMaximize() {
		window.runtime.WindowToggleMaximise();
		isMaximized = !isMaximized;
	}

	function quit() {
		window.runtime.Quit();
	}

	let openSettings = false;
	let isMaximized = false;
</script>

<Header data-wails-drag {...$$restProps}>
	<slot />

	<HeaderUtilities>
		<HeaderButton
			aria-label="Refresh"
			icon={Renew20}
			on:click={refreshCurrentTab}
		/>
		<HeaderButton
			aria-label="Settings"
			icon={SettingsAdjust20}
			on:click={() => (openSettings = true)}
		/>
		<HeaderButton aria-label="Minimize" icon={Subtract24} on:click={minimize} />
		{#if isMaximized}
			<div class:cx--flip-horizontal={true}>
				<HeaderButton
					aria-label="Unmaximize"
					icon={Copy16}
					on:click={toggleMaximize}
				/>
			</div>
		{:else}
			<HeaderButton
				aria-label="Maximize"
				icon={Stop20}
				on:click={toggleMaximize}
			/>
		{/if}
		<HeaderButton aria-label="Close" icon={Close24} isClose on:click={quit} />
	</HeaderUtilities>
</Header>

<Settings bind:open={openSettings} />

<style>
	:global(.cx--flip-horizontal svg) {
		transform: scaleX(-1);
	}
</style>
