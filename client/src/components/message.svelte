<script lang="ts">
	import { onMount } from 'svelte';
	import { GetMessage, SendMessage } from '../helpers/messages';
	import type { Messages,Conversation } from '../helpers/messages';

	let { id, user_id } = $props();
	
	type MessageWS = {
		body: string;
		id: string;
		user_id: string;
	}

	let MessageData: Conversation | null = $state(null);
	let NewMessage: MessageWS = $state({body: "", id: id, user_id: user_id});

	var socketUrl = process.env.NODE_ENV === "production" ? "pokefind-server.fly.dev" : "localhost:8080";
	var socket = new WebSocket("ws://" + socketUrl + "/connect")

	socket.onopen = () =>  {
		console.log("Status: Connected")
	}

	socket.onmessage = (e) => {
		console.log(`message: ${e.data}\n`)
	}

	socket.onclose = () => {
		console.log("Status: Disconnected")
	} 

	onMount(async () => {
		MessageData = await GetMessage(id);
	});


	async function sendMessage(e: Event) {
		e.preventDefault();
		socket.send(JSON.stringify(NewMessage))
		console.log("sending message " + id )
		// SendMessage(NewMessage, id)
	}
</script>

<main>
	<div class="messages">
		{#if MessageData !== null}
			{#each MessageData.Messages as message}
			<div class="message-holder">
				<!-- svelte-ignore a11y_missing_attribute -->
				<img src={message.From_user.Profile_picture} >
				<p>{message.message}</p>
			</div>
			{/each}
		{/if}
	</div>
	
	<form onsubmit={sendMessage}>
		<input type="text" placeholder="send a message..." bind:value={NewMessage.body} />
		<button type="submit"><img src="" alt="send message" /></button>
	</form>
</main>

<style>
	main {
		display: grid;
		grid-template-rows: 60rem 1fr;
		align-items: center;
		
		/* background: blue; */
		margin: 0;
		padding: 0;
		height: 100vh;
		padding-bottom: 1rem;
		margin-top: 5rem;
		padding-left: 5rem;
	}
	form {
		width: 25rem;
		display: grid;
		grid-template-columns: 20rem 1fr;
		margin-left: 25rem;
		padding-bottom: 10rem;
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

	.messages {
		display: flex;
		flex-direction: column;
		justify-content: start;
		border: blue 1px solid;
		height: 100%;
		width: 90svh;
	}

	.message-holder {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.message-holder img{
		width: 3rem;
		height: 3rem;
		border-radius: 50%;
	}
</style>
