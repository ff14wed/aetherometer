<!--
Copyright 2019 carbon-component-svelte authors
Modifications Copyright 2022 Flawed

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
<script>
	/**
	 * Specify the tab label
	 * Alternatively, use the default slot (e.g., <Tab><span>Label</span></Tab>)
	 */
	export let label = "";
	/** Specify the href attribute */
	export let href = "#";
	/** Set to `true` to disable the tab */
	export let disabled = false;
	/** Specify the tabindex */
	export let tabindex = 0;
	/** Set an id for the top-level element */
	export let id = "ccs-" + Math.random().toString(36);
	/** Obtain a reference to the anchor HTML element */
	export let ref = null;
	import { onMount, afterUpdate, getContext, tick } from "svelte";
	const { selectedTab, useAutoWidth, add, update, change } = getContext("Tabs");
	add({ id, label, disabled });
	let didMount = false;
	$: selected = $selectedTab === id;
	onMount(() => {
		tick().then(() => {
			didMount = true;
		});
	});
	afterUpdate(() => {
		if (didMount && selected && ref) {
			ref.focus();
			ref.scrollIntoView({ behavior: "smooth" });
		}
	});
</script>

<!-- svelte-ignore a11y-mouse-events-have-key-events -->
<li
	tabindex="-1"
	role="presentation"
	class:bx--tabs__nav-item={true}
	class:bx--tabs__nav-item--disabled={disabled}
	class:bx--tabs__nav-item--selected={selected}
	{...$$restProps}
	on:click|preventDefault
	on:click|preventDefault={() => {
		if (!disabled) {
			update(id);
		}
	}}
	on:mouseover
	on:mouseenter
	on:mouseleave
	on:keydown={({ key }) => {
		if (!disabled) {
			if (key === "ArrowRight") {
				change(1);
			} else if (key === "ArrowLeft") {
				change(-1);
			} else if (key === " " || key === "Enter") {
				update(id);
			}
		}
	}}
>
	<a
		bind:this={ref}
		role="tab"
		tabindex={disabled ? -1 : tabindex}
		aria-selected={selected}
		aria-disabled={disabled}
		{id}
		{href}
		class:bx--tabs__nav-link={true}
		class:cx--tabs__nav-link={true}
		style={$useAutoWidth ? "width: auto" : undefined}
	>
		<slot>{label}</slot>
	</a>
</li>

<style>
	.cx--tabs__nav-link.cx--tabs__nav-link.cx--tabs__nav-link {
		line-height: 1.5rem;
		height: 3rem;
		outline: none;
		outline-offset: 0px;
	}
	.cx--tabs__nav-link:focus {
		outline: none;
	}
</style>
