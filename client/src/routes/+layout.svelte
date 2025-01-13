<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { incomingUser } from '../helpers/users.js';
	import { GetCurrentUser } from '../helpers/users.js';

	let userData: incomingUser | null = $state(null);
	let loading: boolean = $state(true);
	let { children } = $props();

	onMount(async () => {
		userData = await GetCurrentUser();
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
		<img class="pfp" src={userData?.Profile_picture} />
	</button>
{/if}

{@render children()}

<!-- <slot {userData} /> -->

<style>
	.pfp {
		/* z-index: 3; */
		position: relative;
		width: 3rem;
		height: 3rem;
		border-radius: 50%;
		padding: 0;
		margin: 0;
	}

	.pfp-holder {
		position: absolute;
		top: 1rem;
		right: 1rem;
		width: 3rem;
		height: 3rem;
		padding: 0;
		margin: 0;
		border: none;
		background-color: transparent;

	}
</style>
