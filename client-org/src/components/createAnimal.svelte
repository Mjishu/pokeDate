<script lang="ts">
      import { createAnimal } from "../helper/animals";
      import type { NewAnimal } from "../helper/animals";

      let newAnimal = $state<NewAnimal>({
            name: "",
            species: "",
            date_of_birth: undefined,
            sex: "",
            available: false,
            breed: "",
            price: 0,
      });
</script>

<main>
      <form onsubmit={(e) => createAnimal(e, newAnimal)}>
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
                        bind:value={newAnimal.name}
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
            <button>Create</button>
      </form>
</main>

<style>
</style>
