<script lang="ts">
	import "carbon-components-svelte/css/g100.css";

	import { Header, HeaderUtilities } from "carbon-components-svelte";

	import SettingsAdjust20 from "carbon-icons-svelte/lib/SettingsAdjust20";
	import Subtract20 from "carbon-icons-svelte/lib/Subtract20";
	import Stop16 from "carbon-icons-svelte/lib/Stop16";
	import Close20 from "carbon-icons-svelte/lib/Close20";
	import HeaderButton from "./HeaderButton.svelte";
	import Settings from "./Settings.svelte";

	let openSettings = false;
</script>

<Header data-wails-drag {...$$restProps}>
	<slot />

	<HeaderUtilities>
		<HeaderButton
			aria-label="Settings"
			icon={SettingsAdjust20}
			on:click={() => (openSettings = true)}
		/>
		<HeaderButton aria-label="Minimize" icon={Subtract20} />
		<HeaderButton aria-label="Maximize" icon={Stop16} />
		<HeaderButton aria-label="Close" icon={Close20} isClose />
	</HeaderUtilities>
</Header>

<Settings bind:open={openSettings} />
