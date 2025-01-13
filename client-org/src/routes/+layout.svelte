<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { GetCurrentOrganization } from '../helper/auth';
	import type { Organization } from '../helper/auth.ts';

	let orgData: Organization | null = $state(null);
	let loading: boolean = $state(true);
	let { children } = $props();

	onMount(async () => {
		orgData = await GetCurrentOrganization();
		loading = false;
	});

	function pfpClick() {
		goto('/profile');
	}
</script>

{#if loading}
	<p>loading...</p>
{:else}
	<button class="pfp-holder" onclick={pfpClick}>
		<!-- svelte-ignore a11y_missing_attribute -->
		<img class="pfp" src={orgData?.Profile_picture} />
	</button>
{/if}

{@render children()}

<!-- <slot {userData} /> -->

<style>
	.pfp {
		position: absolute;
		right: 1rem;
		top: 1rem;
		width: 3rem;
		height: 3rem;
		border-radius: 50%;
	}

	.pfp-holder {
		border: none;
		background-color: transparent;
	}
</style>
