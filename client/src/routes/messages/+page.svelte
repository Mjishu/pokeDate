<script lang="ts">
	//todo make a layout for messages and children that shows the messages on the left side and then we can update the inner/whitespace part with the messages
	// when they get selected
	import Message from '../../components/message.svelte';

	let id =
		'get id from somemwhere. Needs to be the message selected id, so get it from layout? im not sure';
</script>

<h1>Messages</h1>

<Message {id} />

<style>
</style>
