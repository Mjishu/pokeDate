<script lang="ts">
	import { onMount } from "svelte";
	import { GetTokens } from "../helpers/users";
	import { CreateMessage } from "../helpers/messages";
	import { json } from "@sveltejs/kit";

      type Notification = {
            Id: string;
            Actor: string;
            Notifier: string;
            Animal_id: string;
            Entity_text: string;
            Entity_type: number;
            status: string;
            date_created: Date;
      }

      let notifications: Notification[] | null = $state(null)
      let loading: boolean = $state(true)
      
      async function GetNotifications() {
            await GetTokens()
            const fetchParams = {
                  method:"GET",
                  headers:{
                        "Content-Type" :"application/json",
                        "Authorization" : `Bearer ${localStorage.getItem("token")}`
                  }
            }
            try {
                  const response = await fetch("/api/notifications", fetchParams)
                  const data = await response.json()
                  if (!response.ok) { return null}
                  loading = false;
                  return data
            } catch(err) {
                  console.error(`error getting notifications ${err}`)
                  loading = false;
                  return null
            }
      }

      onMount(async() => {
            notifications = await GetNotifications();
      })

      async function CreateConversation(actor_id: string) { // takes in the user who SENT the message request
            const fetchParams = {
                  method: "POST",
                  headers:{
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${localStorage.getItem("token")}`
                  }, 
                  body: JSON.stringify(actor_id)
            }
            try {
                  const response = await fetch("/api/messages/create", fetchParams)
                  const data = await response.json()
                  if (!response.ok) {
                        console.error("issue creating message")
                  }
            } catch(err){
                  console.error("could not create message ")
            }
      }
</script>

{#if loading}
      <h2>Loading</h2>
{:else if notifications === null}
      <h2>No new Notifications</h2>
{:else}
      <main>
            {#each notifications as notification}
                  <div class="notification">
                        <p>{notification.Entity_text}</p>
                        <p>{notification.Actor}</p>
                        {#if notification.Entity_type == 1}
                        <button onclick={() => CreateConversation(notification.Actor)}>Start Message</button>
                        {/if}
                  </div>
            {/each}
      </main>
{/if}

<style>
      .notification {
            border: 1px solid blue;
      }
</style>