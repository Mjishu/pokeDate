<script lang="ts">
	import Home from '../components/homeUser.svelte';
	import OrgHome from '../components/homeOrganization.svelte'
	import type {incomingUser } from "../helpers/users"
	
	import { GetCurrentUser, GetTokens } from '../helpers/users';
	import { onMount } from 'svelte';
	let currentUser: incomingUser | null = $state( null)

	// $effect(async () => {
	// 	GetTokens();
	// 	currentUser = await GetCurrentUser();
	// });
	onMount(async () => {
		GetTokens();
		currentUser = await GetCurrentUser();
	});

</script>

<main>
	{#if !currentUser?.Is_organization}
	<Home />
	{:else if currentUser?.Is_organization}
	<OrgHome/>
	{:else}
	<p>no organization found</p>
	{/if}

</main>

<style>
	main {
		height: 100%;
	}
</style>
