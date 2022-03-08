<script>
	import { selectedTabID } from "./stores/stores";

	/** Set an id for the top-level element */
	export let id = "ccs-" + Math.random().toString(36);
	export let label = id;

	$: selected = $selectedTabID === id;
</script>

<div
	role="tabpanel"
	aria-labelledby={label}
	aria-hidden={!selected}
	class:cx--hidden={!selected}
	{id}
	{...$$restProps}
>
	<slot />
</div>

<style>
	div {
		height: 100%;
		width: 100%;
		display: flex;
	}
	.cx--hidden {
		display: none;
	}
</style>
