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

<main>

	<div class="messages">
		{#if MessageData !== null}
		{#each MessageData.Messages as data}
		<p>data.content</p>
		{/each}
		{/if}
	</div>
	
	<form onsubmit={sendMessage}>
		<input type="text" placeholder="send a message..." bind:value={NewMessage} />
		<button type="submit"><img src="" alt="send message" /></button>
	</form>
</main>

<style>
	main {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: space-around;
		/* background: blue; */
		height: 100vh;
		padding-bottom: 1rem;
	}
	form {
		width: 25rem;
		display: grid;
		grid-template-columns: 20rem 1fr;
	}

	form input {
		border: 1px solid #979797;
		border-radius: 15px;
		border-top-right-radius: 0;
		border-bottom-right-radius: 0;
	}

	form button {
		background: transparent;
		transition:
			ease-in-out 300ms background-color,
			ease-in 300ms color;
		background-color: #dadada;
		color: rgb(32, 32, 32);
		border: 1px solid #979797;
		border-radius: 15px;
		border-top-left-radius: 0;
		border-bottom-left-radius: 0;
	}

	form button:hover {
		background-color: rgb(32, 32, 32);
		color: #dadada;
	}
</style>
