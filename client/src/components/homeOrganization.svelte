<script lang="ts">
      import { onMount } from "svelte";
      // import CreateAnimal from './createAnimal.svelte'
      // import AllAnimals from './allAnimals.svelte'
      import CreateAnimal from "./Actions/createAnimal.svelte"
      import AllAnimals from "../components/allAnimals.svelte"
      import Navbar from "./Navbar.svelte";

      let showNewAnimal = $state<boolean>(false);
      let searchTerm = $state<string>("Name");
</script>

<!-- put createANimal behind a button  -->
<main>
      <Navbar />
      <div class="content">
            <div class="filter-section">
                  <h1>Org Name</h1>
                  <div class="operation-bar">
                        <div class="search-bar">
                              <div class="search-toggle">
                                    <select
                                          name="search-field"
                                          id="search-field"
                                    >
                                          <option value="search-term"
                                                >{searchTerm}</option
                                          >
                                          <!-- {#each AllAnimals as animal}{/each} -->
                                    </select>
                              </div>
                              <div class="search-text">
                                    <img src="/icons/search_icon.svg" alt="" />
                                    <input
                                          type="text"
                                          placeholder={`Search ` + searchTerm}
                                    />
                              </div>
                        </div>
                        <div class="dropdown-items">
                              <select name="fields" id="fields">
                                    <option value="fields">Fields</option>
                              </select>
                              <select name="filter" id="">
                                    <option value="filter">Filter</option>
                              </select>
                              <select name="group-by" id="">
                                    <option value="group-by">Group By</option>
                              </select>
                              <select name="sort" id="">
                                    <option value="sort">Sort</option>
                              </select>
                        </div>
                        <button
                              class="add-new-button"
                              onclick={() => (showNewAnimal = !showNewAnimal)}
                              >Add New</button
                        >
                  </div>
                  <div class="show-new-animal">
                        {#if showNewAnimal}
                              <CreateAnimal bind:showNewAnimal />
                        {/if}
                  </div>
            </div>
            <AllAnimals />
      </div>
</main>

<style>
      main {
            height: 100%;
            padding: 0;
            margin: 0;
            display: grid;
            grid-template-columns: 17.25rem 1fr;
            gap: 5rem;
      }

      .filter-section {
            display: grid;
            grid-template-rows: 15.5rem 6.5rem 1fr;
      }

      .operation-bar {
            display: flex;
            align-items: center;
            justify-content: space-evenly;
            gap: 1.15rem;
            width: 95%;
            display: grid;
            gap: 3.5rem;
            grid-template-columns: 16rem 1fr 10.5rem;
      }

      .dropdown-items {
            justify-self: end;

            display: flex;
            gap: 1rem;
            align-items: center;
      }

      .search-bar {
            width: 16.5rem;
            height: 2rem;
            border: 1px solid var(--border-color);
            border-radius: 10px;
            display: grid;
            grid-template-columns: 35% 1fr;
      }

      .search-toggle {
            border-right: 1px solid var(--border-color, #000);
            display: flex;
      }

      #search-field {
            width: 5.5rem;
            font-size: 14px;
            border: none;
      }

      .search-text {
            border: none;
            display: flex;
            text-align: cetner;
            align-items: center;
            padding-left: 0.5rem;
            gap: 0.25rem;
      }
      .search-text img {
            width: 1rem;
            height: 1.2rem;
      }
      .search-text input {
            width: 80%;
            height: 85%;
            border: none;
            outline: none;
      }
      .search-text input:active,
      .search-text input:focus {
            border: none;
            outline: none;
      }

      .add-new-button {
            background: #30ca4a;
            color: #e9e9e9;
            width: 10.5rem;
            height: 2rem;
            border: none;
            border-radius: 5px;
            font-size: 2rem;
            text-align: center;
            padding: 0px;
            line-height: 0px;
      }

      select {
            height: 2rem;
            width: auto;
            background-color: transparent;
            border: 1px solid var(--border-color);
            border-radius: 10px;
            font-size: 22px;
            color: rgba(99, 99, 99, 0.767);
            text-align: center;
      }
</style>
