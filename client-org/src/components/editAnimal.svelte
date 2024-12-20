<script lang="ts">
      import {
            createAnimal,
            getAnimalById,
            updateAnimalById,
            GetAllShots,
      } from "../helper/animals";
      import type {
            NewAnimal,
            UpdatedAnimal,
            Shot,
            Animal,
      } from "../helper/animals";
      import { formatISO } from "date-fns";
      import { onMount } from "svelte";

      let { showEditPage = $bindable(), currentId } = $props();
      let oldAnimal = $state();
      let animal = $state<Animal>();
      let shotData = $state<Shot[]>();

      let updatedAnimal = $state<UpdatedAnimal>({
            // switch this to animal.name etc
            name: "",
            date_of_birth: "",
            price: 0,
            available: false,
            shots: [],
      });

      function closeForm() {
            showEditPage = false;
      }

      function addNewShot() {
            updatedAnimal.shots.push({
                  id: "",
                  date_given: "",
                  date_due: "",
            });
      }
      onMount(async () => {
            shotData = await GetAllShots();

            animal = await getAnimalById(currentId);
            updatedAnimal = {
                  name: animal?.Name || "",
                  date_of_birth: animal?.Date_of_birth || "",
                  price: animal?.Price || 0,
                  available: animal?.Available || false,
                  shots: animal?.Shots || [], // figure out how to add animal.shots, will have to iterate over each shot in updatedAnimal and show the shots
            };
      });
</script>

<main>
      <form
            onsubmit={(e) => {
                  closeForm();
                  updateAnimalById(e, currentId, updatedAnimal);
            }}
            autocomplete="off"
      >
            <h3>Information</h3>
            <div>
                  <label for="name">Name</label>
                  <input
                        required
                        type="text"
                        name="name"
                        bind:value={updatedAnimal.name}
                  />
            </div>
            <div>
                  <label for="date_of_birth">Date of Birth</label>
                  <input
                        required
                        type="date"
                        name="date_of_birth"
                        bind:value={updatedAnimal.date_of_birth}
                  />
            </div>
            <div>
                  <label for="price">Price</label>
                  <input
                        required
                        type="number"
                        name="price"
                        step=".01"
                        min="0"
                        max="9999999"
                        bind:value={updatedAnimal.price}
                  />
            </div>
            <div>
                  <label for="available">Available</label>
                  <input
                        type="checkbox"
                        name="available"
                        bind:checked={updatedAnimal.available}
                  />
            </div>
            <hr />
            <h3>Shots</h3>
            <div class="shots">
                  {#each updatedAnimal.shots as shot, i}
                        <div class="shot-wrapper">
                              <div>
                                    <label for="shot-name">Name</label>
                                    <select
                                          name="shot-name"
                                          id="sahot-name"
                                          bind:value={updatedAnimal.shots[i].id}
                                    >
                                          <option value="" disabled selected
                                                >Name</option
                                          >
                                          {#if shotData != undefined}
                                                {#each shotData as shot}
                                                      <option value={shot.Id}
                                                            >{shot.Name}</option
                                                      >
                                                {/each}
                                          {/if}
                                    </select>
                              </div>
                              <div>
                                    <label for="shot-given">Shot Given</label>
                                    <input
                                          type="date"
                                          name="shot-given"
                                          bind:value={updatedAnimal.shots[i]
                                                .date_given}
                                    />
                              </div>
                              <div>
                                    <label for="shot-due">Next Due</label>
                                    <input
                                          type="date"
                                          name="shot-due"
                                          bind:value={updatedAnimal.shots[i]
                                                .date_due}
                                    />
                              </div>
                        </div>
                  {/each}
                  <button type="button" onclick={addNewShot}>+</button>
            </div>
            <hr />
            <button type="submit">Create</button>
            <button type="button" onclick={closeForm}>Cancel</button>
      </form>
</main>

<style>
      main {
            position: absolute;
            background-color: #5a5959;
            border-radius: 5px;
            color: white;
      }
      form {
            padding: 3rem;
      }

      .shots {
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
            align-items: center;
      }

      .shots button {
            width: 5rem;
            height: 1.5rem;
      }

      .shot-wrapper {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1rem;
      }
      .shot-wrapper div {
            display: flex;
            flex-direction: column;
      }
</style>
