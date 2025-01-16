<script lang="ts">
      import { onMount } from "svelte";
      import {GetCurrentUser} from "../../helpers/users"
      import type {incomingUser} from "../../helpers/users"
      let CurrentUser: incomingUser |null  = $state(null)

      onMount(async() => {
            CurrentUser = await GetCurrentUser()
      })
</script>   

{#if CurrentUser !== null}
      {#if CurrentUser.Is_organization}
            <h4>Open a Message?</h4>
            <p>or</p>
            <p>check your inbox for message requests</p>
      {:else if !CurrentUser.Is_organization}
            <h4>Open a Message?</h4>
            <p>or</p>
            <p>Go Find Some <a href="/">animals</a></p>
      {/if}
{/if}

<style>

</style>