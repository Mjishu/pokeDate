<script lang="ts">
      import { createAnimal } from "../helper/animals";
      import type { NewAnimal } from "../helper/animals";
      import { formatISO } from "date-fns";

      let { showNewAnimal = $bindable() } = $props();

      let newAnimal = $state<NewAnimal>({
            name: "",
            species: "",
            date_of_birth: undefined,
            sex: "",
            available: false,
            breed: "",
            price: 0,
            shots: [{ name: "", date_given: "", date_due: "" }],
      });

      async function handleCreateAnimal(e: Event) {
            e.preventDefault();

            const formattedAnimal = {
                  ...$state.snapshot(newAnimal),
                  date_of_birth: newAnimal.date_of_birth
                        ? formatISO(new Date(newAnimal.date_of_birth))
                        : undefined,
            };
            console.log(formattedAnimal);
            await createAnimal(formattedAnimal);
            showNewAnimal = false;
      }

      function closeForm() {
            showNewAnimal = false;
      }

      function addNewShot() {
            newAnimal.shots.push({ name: "", date_given: "", date_due: "" });
            console.log("add new shots was called");
      }
</script>

<main>
      <form onsubmit={handleCreateAnimal} autocomplete="off">
            <h3>Information</h3>
            <div>
                  <label for="name">Name</label>
                  <input
                        type="text"
                        name="name"
                        bind:value={newAnimal.name}
                        required
                  />
            </div>
            <div>
                  <label for="species">Species</label>
                  <input
                        type="text"
                        name="species"
                        bind:value={newAnimal.species}
                        required
                  />
            </div>
            <div>
                  <label for="date_of_birth">Date of Birth</label>
                  <input
                        type="date"
                        name="date_of_birth"
                        bind:value={newAnimal.date_of_birth}
                        required
                  />
            </div>
            <div>
                  <label for="sex">Sex</label>
                  <select
                        name="sex"
                        id="sex"
                        bind:value={newAnimal.sex}
                        required
                  >
                        <option value="" disabled selected>Sex</option>
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                        <option value="undefined">Undefined</option>
                  </select>
                  <!-- <input type="text" name="sex" bind:value={newAnimal.sex} /> -->
            </div>
            <div>
                  <label for="price">Price</label>
                  <input
                        required
                        type="number"
                        name="price"
                        bind:value={newAnimal.price}
                  />
            </div>
            <div>
                  <label for="available">Available</label>
                  <input
                        required
                        type="checkbox"
                        name="available"
                        bind:checked={newAnimal.available}
                  />
            </div>
            <div>
                  <label for="breed">Breed</label>
                  <input
                        required
                        type="text"
                        name="breed"
                        bind:value={newAnimal.breed}
                  />
            </div>
            <hr />
            <h3>Shots</h3>
            <div class="shots">
                  {#each newAnimal.shots as shot, i}
                        <div class="shot-wrapper">
                              <div>
                                    <label for="shot-name">Name</label>
                                    <input
                                          type="text"
                                          placeholder="name..."
                                          name="shot-name"
                                          bind:value={newAnimal.shots[i].name}
                                    />
                              </div>
                              <div>
                                    <label for="shot-given">Shot Given</label>
                                    <input
                                          type="date"
                                          name="shot-given"
                                          bind:value={newAnimal.shots[i]
                                                .date_given}
                                    />
                              </div>
                              <div>
                                    <label for="shot-due">Next Due</label>
                                    <input
                                          type="date"
                                          name="shot-due"
                                          bind:value={newAnimal.shots[i]
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
