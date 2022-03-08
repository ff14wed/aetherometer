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
	import { onMount } from "svelte";

	export let open = false;

	let config = {
		Plugins: {},
	};

	onMount(async () => {
		await window.go.main.App.WaitForStartup();

		config = (await window.go.main.App.GetConfig()) || config;

		window.runtime.EventsOn("ConfigChange", async () => {
			config = (await window.go.main.App.GetConfig()) || config;
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

	let newPluginName;
	let newPluginURL;
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
	<Form>
		<FormGroup>
			<StructuredList condensed>
				<StructuredListHead>
					<StructuredListRow head>
						<StructuredListCell head>
							<div class:cx--plugin={true}>Plugin Name</div>
						</StructuredListCell>
						<StructuredListCell head>
							<div class:cx--plugin={true}>Plugin URL</div>
						</StructuredListCell>
						<StructuredListCell head />
					</StructuredListRow>
				</StructuredListHead>
				<StructuredListBody>
					{#each Object.entries(config.Plugins) as [name, url]}
						<StructuredListRow>
							<StructuredListCell>
								<div class:cx--plugin={true}>
									{name}
								</div>
							</StructuredListCell>
							<StructuredListCell>
								<div class:cx--plugin={true}>
									{url}
								</div>
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
	.cx--plugin {
		padding: 0 1rem;
	}

	:global(.bx--structured-list-td) {
		vertical-align: middle;
	}
</style>
