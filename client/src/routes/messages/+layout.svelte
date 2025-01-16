<script lang="ts">
	import { onMount } from 'svelte';
	import { CurrentUserMessages } from '../../helpers/messages';
      import type {Messages, Conversation} from "../../helpers/messages"
	import { goto } from '$app/navigation';
      let {children} = $props()
 
	let AllMessages: Conversation[] | null = $state(null);

	onMount(async () => {
            AllMessages = await CurrentUserMessages()
      });

      async function OpenMessage(id: string) {
            goto(`/messages/${id}`)
      }
</script>

<div class="content">
      <div class="layout-content">
            {#if AllMessages !== null}
                  {#each AllMessages as message}
                        <button class="message" type="button" onclick={() => OpenMessage(message.Id)} aria-label="Open message">
                              <h4>{message.Conversation_name}</h4>
                        </button>
                  {/each}
            {:else}
            <h3>Could not find any messges!</h3>
            {/if}
      </div>
      <div class="message-content">
            {@render children()}
      </div>
</div>
      
      <style>
            :root {
                  margin: 0;
                  padding: 0;
            }
            .content {
                  display: grid;
                  grid-template-columns: 20rem 1fr;
                  height: 100vh;
                  width: 100%;
                  padding: 0;
                  margin: 0;
                  /* background: pink; */
            }

            .layout-content {
                  border-right: 1px solid #9c9c9c;
                  overflow-y: auto;
                  padding: 1rem;
                  height: 100vh;
                  /* background: blue; */
                  /* padding-top: 0; */
            }

            .message-content{
                  width: 100%;
                  height: 100vh;
                  text-align: center;
                  padding: 0.5rem;
                  border: none;
                  background: transparent;
                  cursor: pointer;
                  /* background: green; */
            }

            .message {
                  width: 100%;
                  background: transparent;
                  border: 1px solid rgb(235, 235, 235);
                  border-radius: 5px;
                  transition: ease-in-out 300ms border;
            }
            .message:hover{
                  border: 1px solid rgb(185,185,185);
            }
      </style>

