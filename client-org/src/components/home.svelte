<script lang="ts">
      import { onMount } from "svelte";
      import CreateAnimal from "./createAnimal.svelte";

      let showNewAnimal = $state<boolean>(false);
</script>

<!-- put createANimal behind a button  -->
<main>
      <div class="navbar"></div>
      <div class="main">
            <div class="operation-bar">
                  <div class="search-bar"></div>
                  <label for="fields">Fields</label>
                  <select name="fields" id=""> </select>
                  <label for="filter">Filter</label>
                  <select name="filter" id=""> </select>
                  <label for="group-by">Group by</label>
                  <select name="group-by" id=""> </select>
                  <label for="sort">Sort</label>
                  <select name="sort" id=""> </select>
                  <button
                        class="add-new-button"
                        onclick={() => (showNewAnimal = true)}>Add New</button
                  >
            </div>
            <div class="show-new-animal">
                  {#if showNewAnimal}
                        <CreateAnimal {showNewAnimal} />
                  {/if}
            </div>
            <div class="show-animals"></div>
      </div>
</main>

<style>
      main {
            display: grid;
            grid-template-columns: repeat(2, 1f);
      }

      .show-animals {
            border-radius: 5px;
            border: none;
            background: hsla(0, 1%, 77%, 0.5);
            width: 66.25rem;
            height: 42.75rem;
      }
</style>
