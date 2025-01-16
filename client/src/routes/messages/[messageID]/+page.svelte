<script lang="ts">
	//todo make a layout for messages and children that shows the messages on the left side and then we can update the inner/whitespace part with the messages
	// when they get selected
	import {page} from '$app/stores'
	import { onMount } from 'svelte';
	import Message from '../../../components/message.svelte';
	
	// const id = $page.params.messageID
	let id: string;

	$: id = $page.params.messageID;

	onMount(() => {
		console.log("message id is here")
		console.log(id)
	})
</script>

<h1>Messages</h1>

<Message {id} />

<style>
</style>
