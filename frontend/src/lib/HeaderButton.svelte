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
	/** Set to `true` to use the variant with the red hover */
	export let isClose = false;

	/**
	 * Specify the icon to render
	 * @type {typeof import("svelte").SvelteComponent}
	 */
	export let icon = undefined;

	/** Obtain a reference to the HTML button element */
	export let ref = null;
</script>

<button
	type="button"
	bind:this={ref}
	class:cx--header__button={true}
	class:cx--header__close={isClose}
	{...$$restProps}
	on:click
>
	<slot>
		<svelte:component this={icon} />
	</slot>
</button>

<style>
	.cx--header__button {
		display: inline-block;
		padding: 0;
		border: 0;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background: none;
		cursor: pointer;
		width: 100%;
		width: 3rem;
		height: 3rem;
		border: 0.0625rem solid transparent;
		transition: background-color 110ms, border-color 110ms;
	}
	.cx--header__button::-moz-focus-inner {
		border: 0;
	}
	.cx--header__button:hover {
		background-color: #353535;
	}
	.cx--header__close:hover {
		background-color: #b81921;
	}

	.cx--header__button--active {
		border-right: 1px solid #393939;
		border-bottom: 1px solid #161616;
		border-left: 1px solid #393939;
	}
	.cx--header__button:active {
		background-color: #393939;
	}
	.cx--header__button.bx--btn--icon-only.bx--tooltip__trigger {
		justify-content: center;
	}
	:global(.cx--header__button > svg) {
		fill: #fff;
		vertical-align: middle;
	}
</style>
