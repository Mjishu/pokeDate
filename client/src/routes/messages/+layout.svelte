<script lang="ts">
	import { onMount } from 'svelte';
	import { CurrentUserMessages } from '../../helpers/messages';
      import type {Messages, Conversation} from "../../helpers/messages"
      let {children} = $props()
 
	let AllMessages: Conversation[] | null = $state(null);

	onMount(async () => {
            AllMessages = await CurrentUserMessages()
      });

      async function OpenMessage(id: string) {
            console.log(id)
      }
</script>

{#if AllMessages !== null}
      {#each AllMessages as message}
      <button type="button" onclick={() => OpenMessage(message.Id)} aria-label="Open message">
            <h4>{message.Conversation_name}</h4>
      </button>
      {/each}
{:else}
<h3>Could not find any messges!</h3>
{/if}

<style>
</style>

{@render children()}