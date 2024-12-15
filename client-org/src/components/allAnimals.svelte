<script>
      import { getOrganizationAnimals } from "../helper/animals";

      let animalsData = $state();
      let isLoading = $state(true);

      async function callAnimals() {
            animalsData = await getOrganizationAnimals();
            isLoading = false;
      }
      callAnimals();
</script>

<!-- get all animals -> make sure it is from X org(the one that is logged in) -->
{#if isLoading}
      <div>Loading...</div>
{:else}
      <div class="show-animals"></div>
{/if}

<style>
      .show-animals {
            border-radius: 5px;
            border: none;
            background: hsla(0, 1%, 77%, 0.5);
            width: 95%;
            height: 58rem;
      }
</style>
