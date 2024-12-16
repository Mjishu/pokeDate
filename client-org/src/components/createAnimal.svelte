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
      });

      async function handleCreateAnimal(e: Event) {
            e.preventDefault();

            const formattedAnimal = {
                  ...$state.snapshot(newAnimal),
                  date_of_birth: newAnimal.date_of_birth
                        ? formatISO(new Date(newAnimal.date_of_birth))
                        : undefined,
            };
            await createAnimal(formattedAnimal);
            showNewAnimal = false;
      }

      function closeForm() {
            showNewAnimal = false;
      }
</script>

<main>
      <form onsubmit={handleCreateAnimal} autocomplete="off">
            <div>
                  <label for="name">Name</label>
                  <input type="text" name="name" bind:value={newAnimal.name} />
            </div>
            <div>
                  <label for="species">Species</label>
                  <input
                        type="text"
                        name="species"
                        bind:value={newAnimal.species}
                  />
            </div>
            <div>
                  <label for="date_of_birth">Date of Birth</label>
                  <input
                        type="date"
                        name="date_of_birth"
                        bind:value={newAnimal.date_of_birth}
                  />
            </div>
            <div>
                  <label for="sex">Sex</label>
                  <input type="text" name="sex" bind:value={newAnimal.sex} />
            </div>
            <div>
                  <label for="price">Price</label>
                  <input
                        type="number"
                        name="price"
                        bind:value={newAnimal.price}
                  />
            </div>
            <div>
                  <label for="available">Available</label>
                  <input
                        type="checkbox"
                        name="available"
                        bind:checked={newAnimal.available}
                  />
            </div>
            <div>
                  <label for="breed">Breed</label>
                  <input
                        type="text"
                        name="breed"
                        bind:value={newAnimal.breed}
                  />
            </div>
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
</style>
