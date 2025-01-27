<script lang="ts">
	//todo make a layout for messages and children that shows the messages on the left side and then we can update the inner/whitespace part with the messages
	// when they get selected
	import {page} from '$app/stores'
	import { onMount } from 'svelte';
	import Message from '../../../components/message.svelte';
	import { GetCurrentUser } from '../../../helpers/users';
	import type {incomingUser} from "../../../helpers/users"
	
	// const id = $page.params.messageID
	var id: string;
	$: id = $page.params.messageID;

    	let currentUser: incomingUser | null = null;

	onMount(async () => {
		currentUser = await GetCurrentUser()
	})
</script>

<h1>Messages</h1>

{#if currentUser}
	<Message {id} user_id={currentUser.Id} />
{/if}

<style>
</style>
