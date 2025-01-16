<script lang="ts">
	import { onMount } from 'svelte';
	import { CurrentUserMessages } from '../../helpers/messages';
      import type {Messages, Conversation} from "../../helpers/messages"

	let AllMessages: Conversation[] | null = $state(null);

	onMount(async () => {
            AllMessages = await CurrentUserMessages()
      });
</script>

{#if AllMessages !== null}
      {#each AllMessages as message}
            <h4>{message.Conversation_name}</h4>
      {/each}
{:else}
<h3>Could not find any messges!</h3>
{/if}

<style>
</style>
