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
	/** Specify the selected tab index */
	export let selected = 0;
	export let selectedTabID = "";
	/**
	 * Specify the type of tabs
	 * @type {"default" | "container"}
	 */
	export let type = "default";
	/** Set to `true` for tabs to have an auto-width */
	export let autoWidth = false;
	/**
	 * Specify the ARIA label for the chevron icon
	 * @type {string}
	 */
	import { createEventDispatcher, afterUpdate, setContext } from "svelte";
	import { writable, derived } from "svelte/store";
	import ChevronLeftGlyph from "carbon-icons-svelte/lib/ChevronLeftGlyph";
	import ChevronRightGlyph from "carbon-icons-svelte/lib/ChevronRightGlyph";

	import { tweened } from "svelte/motion";
	import { cubicOut } from "svelte/easing";

	const dispatch = createEventDispatcher();
	const tabs = writable([]);
	const tabsById = derived(tabs, (_) =>
		_.reduce((a, c) => ({ ...a, [c.id]: c }), {})
	);
	const useAutoWidth = writable(autoWidth);
	const selectedTab = writable(undefined);
	setContext("Tabs", {
		tabs,
		selectedTab,
		useAutoWidth,
		add: (data) => {
			tabs.update((_) => [..._, { ...data, index: _.length }]);
		},
		update: (id) => {
			currentIndex = $tabsById[id].index;
			dispatch("change", currentIndex);
		},
		change: (direction) => {
			let index = currentIndex + direction;
			if (index < 0) {
				index = $tabs.length - 1;
			} else if (index >= $tabs.length) {
				index = 0;
			}
			let disabled = $tabs[index].disabled;
			while (disabled) {
				index = index + direction;
				if (index < 0) {
					index = $tabs.length - 1;
				} else if (index >= $tabs.length) {
					index = 0;
				}
				disabled = $tabs[index].disabled;
			}
			currentIndex = index;
			dispatch("change", currentIndex);
		},
	});
	afterUpdate(() => {
		checkScroll();
		selected = currentIndex;
	});

	let currentIndex = selected;
	$: currentIndex = selected;
	$: currentTab = $tabs[currentIndex] || undefined;
	$: {
		if (currentTab) {
			selectedTab.set(currentTab.id);
			selectedTabID = currentTab.id;
		}
	}
	$: useAutoWidth.set(autoWidth);

	let horizontalOverflow = true;
	let leftOverflowNavButtonDisabled = false;
	let rightOverflowNavButtonDisabled = false;

	let tablist;
	const OVERFLOW_BUTTON_OFFSET = 40;
	let overflowNavInterval;

	function checkScroll() {
		if (tablist) {
			horizontalOverflow = tablist.scrollWidth > tablist.clientWidth;
			if ($$slots.default) {
				leftOverflowNavButtonDisabled = tablist.scrollLeft <= 0;
				rightOverflowNavButtonDisabled =
					tablist.scrollLeft + tablist.clientWidth >= tablist.scrollWidth;
			}
		}
	}

	function onResize() {
		checkScroll();
	}

	const scrollTarget = tweened(0, {
		duration: 200,
		easing: cubicOut,
	});

	$: if (tablist) {
		tablist.scrollLeft = $scrollTarget;
	}

	function onWheel(e) {
		if (e.deltaY != 0) {
			scrollTarget.set(e.currentTarget.scrollLeft + e.deltaY + e.deltaX);
		}
	}

	function onOverflowClick(e, { direction, multiplier = 40 }) {
		const { scrollLeft } = tablist;
		// account for overflow button appearing and causing tablist width change
		if (direction === 1 && !scrollLeft) {
			scrollTarget.set(scrollLeft + OVERFLOW_BUTTON_OFFSET);
		}
		scrollTarget.set(scrollLeft + direction * multiplier);
		if (leftEdgeReached(direction)) {
			scrollTarget.set(0);
		}
	}
	function onOverflowMouseDown(e, { direction }) {
		// disregard mouse buttons aside from LMB
		if (e.buttons !== 1) {
			return;
		}
		overflowNavInterval = setInterval(() => {
			if (leftEdgeReached(direction) || rightEdgeReached(direction)) {
				clearInterval(overflowNavInterval);
			}
			onOverflowClick(e, { direction });
		});
	}

	function onOverflowMouseUp() {
		clearInterval(overflowNavInterval);
	}

	function leftEdgeReached(direction) {
		const { scrollLeft } = tablist;
		return direction === -1 && scrollLeft <= OVERFLOW_BUTTON_OFFSET;
	}

	function rightEdgeReached(direction) {
		const { clientWidth, scrollLeft, scrollWidth } = tablist;
		return direction === 1 && scrollLeft + clientWidth >= scrollWidth;
	}
</script>

<svelte:window on:resize={onResize} />

<div
	role="navigation"
	class:bx--tabs--scrollable={true}
	class:bx--tabs--container={type === "container"}
	class:cx--tabs--scrollable={true}
	data-wails-drag
	{...$$restProps}
>
	{#if horizontalOverflow}
		<button
			aria-hidden={true}
			aria-label="scroll left"
			class:bx--tab--overflow-nav-button={true}
			class:cx--tab--overflow-nav-button--disabled={leftOverflowNavButtonDisabled}
			on:click|stopPropagation|preventDefault={(e) =>
				onOverflowClick(e, { direction: -1 })}
			on:mousedown|stopPropagation|preventDefault={(e) =>
				onOverflowMouseDown(e, { direction: -1 })}
			on:mouseup|stopPropagation|preventDefault={onOverflowMouseUp}
			tabIndex="-1"
			type="button"
		>
			<ChevronLeftGlyph />
		</button>
		<div class:bx--tabs__overflow-indicator--left={true} />
	{/if}
	<ul
		role="tablist"
		bind:this={tablist}
		class:bx--tabs--scrollable__nav={true}
		data-wails-no-drag
		on:scroll={checkScroll}
		on:wheel|stopPropagation={onWheel}
	>
		<slot />
	</ul>
	{#if horizontalOverflow}
		<div class:bx--tabs__overflow-indicator--right={true} />
		<button
			aria-hidden={true}
			aria-label="scroll right"
			class:bx--tab--overflow-nav-button={true}
			class:cx--tab--overflow-nav-button--disabled={rightOverflowNavButtonDisabled}
			on:click|stopPropagation|preventDefault={(e) =>
				onOverflowClick(e, { direction: 1 })}
			on:mousedown|stopPropagation|preventDefault={(e) =>
				onOverflowMouseDown(e, { direction: 1 })}
			on:mouseup|stopPropagation|preventDefault={onOverflowMouseUp}
			tabIndex="-1"
			type="button"
		>
			<ChevronRightGlyph />
		</button>
	{/if}
</div>

<style>
	:global(.cx--tab--overflow-nav-button--disabled.cx--tab--overflow-nav-button--disabled
			> svg) {
		fill: #525252;
	}
	.cx--tab--overflow-nav-button--disabled:hover {
		cursor: default;
	}
	.cx--tabs--scrollable {
		min-width: 0;
		user-select: none;
	}
</style>
