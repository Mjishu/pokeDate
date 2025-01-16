<script lang="ts">
	import { onMount } from 'svelte';
	import { CurrentUserMessages } from '../../helpers/messages';
      import type {Messages} from "../../helpers/messages"

	let AllMessages: Messages[] | null = $state(null);

	onMount(async () => {
            AllMessages = await CurrentUserMessages()
      });
</script>

{#if AllMessages !== null}
      {#each AllMessages as message}
            <h4>{message}</h4>
      {/each}
{/if}

<style>
</style>
