<script lang='ts'>
      import {onMount} from "svelte"
      import { GetCurrentUser } from "../../helpers/users";
      import type { incomingUser } from '../../helpers/users';
      import Navbar from "../../components/Navbar.svelte";
	import { goto } from "$app/navigation";
      
      let userData: incomingUser | null= $state(null)
      let loading:boolean = $state(true)
      let showResetOverlay: boolean = $state(false)
      
      onMount(async () => {
		userData = await GetCurrentUser();
		loading = false;
	});

      async function resetCards() { //todo put this in settings page instead
		const fetchParams = {
			method:"POST",
			headers: {
				"Content-Type": "application/json",
				"Authorization" : `Bearer ${localStorage.getItem("token")}`
			}
		}
		try {
			const response = await fetch("/api/users/progress/reset", fetchParams)
			const data = await response.json()
			if (!response.ok){
				alert("error resetting cards")
			}
                  goto("/")
		} catch(error) {
			console.error(`error trying to reset cards: ${error}`)
		}
	}
</script>

<div class="content">
      <Navbar/>
      <div class="main">
            {#if loading}
                  <h5>loading</h5>
            {:else if userData === null}
                  <p>Could not find data</p>
            {:else}
                  <h2>Settings</h2>
                  {#if !userData.Is_organization}
                        <button class="reset-card" onclick={() => showResetOverlay = !showResetOverlay}>Reset Cards</button>

                        {#if showResetOverlay}
                        <p>This will remove the progress for every animal you skipped</p>
                        <button class="reset-card" onclick={() => showResetOverlay = false}>Cancel</button>
                        <button class="reset-card" onclick={resetCards}>Reset</button>
                        {/if}
                  {/if}
            {/if}

      </div>
</div>
<style>
      .content {
		display: grid;
		grid-template-columns: 17.25rem 1fr;
		height: 100%;
		gap: 5rem;
	}
      .reset-card{
            width: 8rem;
            height: 3rem;
      }
      
</style>