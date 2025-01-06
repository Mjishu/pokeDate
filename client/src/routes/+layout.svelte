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
		alert('going to profile');
	}
</script>

{#if loading}
	<p>loading...</p>
{:else}
	<button class="pfp-holder" onclick={pfpClick}>
		<img class="pfp" src={userData?.Profile_picture} alt="profile-button" />
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
