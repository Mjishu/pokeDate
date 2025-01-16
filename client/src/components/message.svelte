<script lang="ts">
	import { onMount } from 'svelte';
	import { GetMessage } from '../helpers/messages';
	import type { Messages,Conversation } from '../helpers/messages';

	let { id } = $props();

	let MessageData: Conversation | null = $state(null);
	let NewMessage: string = $state('');

	// let id = "grab something from message called, this will probably be props"

	onMount(async () => {
		MessageData = await GetMessage(id);
	});

	async function sendMessage(e: Event) {
		e.preventDefault();
		console.log(NewMessage);
	}
</script>

{#if MessageData !== null}
	{#each MessageData.Messages as data}
		<p>data.content</p>
	{/each}
{/if}

<form onsubmit={sendMessage}>
	<input type="text" placeholder="send a message..." bind:value={NewMessage} />
	<button type="submit"><img src="" alt="send message" /></button>
</form>

<style>
	form {
		width: 25rem;
		display: grid;
		grid-template-columns: 20rem 1fr;
	}

	form input {
		border: 1px solid #979797;
		border-radius: 10px;
	}

	form button {
		background: transparent;
		transition:
			ease-in-out 300ms background-color,
			ease-in 300ms color;
		background-color: #979797;
		color: rgb(32, 32, 32);
	}

	form button:hover {
		background-color: rgb(32, 32, 32);
		color: #979797;
	}
</style>
