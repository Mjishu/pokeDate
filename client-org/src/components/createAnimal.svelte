<script lang="ts">
      import { onMount } from "svelte";
      import { createAnimal, GetAllShots } from "../helper/animals";
      import type { NewAnimal, Shot } from "../helper/animals";
      import { formatISO } from "date-fns";

      let { showNewAnimal = $bindable() } = $props();

      let shotData = $state<Shot[]>();
      let newAnimal = $state<NewAnimal>({
            Name: "",
            Species: "",
            Date_of_birth: undefined,
            Sex: "",
            Available: false,
            Breed: "",
            Price: 0,
            Shots: [],
            Image_src: "",
      });

      function cleanImageData() {
            // call a backend route that creates image and then add that image
            if (newAnimal.Image_src) {
                  const data = new FormData();
                  data.append("image_src", newAnimal.Image_src);
                  console.log(data);
                  return data;
            }
      }

      async function handleCreateAnimal(e: Event) {
            e.preventDefault();

            const formattedAnimal = {
                  ...$state.snapshot(newAnimal),
                  Date_of_birth: newAnimal.Date_of_birth
                        ? formatISO(new Date(newAnimal.Date_of_birth))
                        : undefined,
                  Image_src: "", //cleanImageData(),
                  Shots: newAnimal.Shots.map((shot) => ({
                        Id: shot.Id,
                        Date_given: formatISO(new Date(shot.Date_given)),
                        Next_due: formatISO(new Date(shot.Next_due)),
                  })),
            };
            console.log(formattedAnimal);
            await createAnimal(formattedAnimal);
            showNewAnimal = false;
      }

      function closeForm() {
            showNewAnimal = false;
      }

      function addNewShot() {
            newAnimal.Shots.push({ Id: 0, Date_given: "", Next_due: "" });
      }

      onMount(async () => {
            shotData = await GetAllShots();
      });
</script>

<main>
      <form onsubmit={handleCreateAnimal} autocomplete="off">
            <h3>Information</h3>
            <div>
                  <label for="name">Name</label>
                  <input
                        type="text"
                        name="name"
                        bind:value={newAnimal.Name}
                        required
                  />
            </div>
            <div>
                  <label for="species">Species</label>
                  <input
                        type="text"
                        name="species"
                        bind:value={newAnimal.Species}
                        required
                  />
            </div>
            <div>
                  <label for="date_of_birth">Date of Birth</label>
                  <input
                        type="date"
                        name="date_of_birth"
                        bind:value={newAnimal.Date_of_birth}
                        required
                  />
            </div>
            <div>
                  <label for="sex">Sex</label>
                  <select
                        name="sex"
                        id="sex"
                        bind:value={newAnimal.Sex}
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
                        step=".01"
                        min="0"
                        max="9999999"
                        bind:value={newAnimal.Price}
                  />
            </div>
            <div>
                  <label for="available">Available</label>
                  <input
                        type="checkbox"
                        name="available"
                        bind:checked={newAnimal.Available}
                  />
            </div>
            <div>
                  <label for="breed">Breed</label>
                  <input
                        required
                        type="text"
                        name="breed"
                        bind:value={newAnimal.Breed}
                  />
            </div>
            <hr />
            <!-- * ADD IMAGES -->
            <h3>Images</h3>
            <div class="image-container">
                  <input
                        multiple={false}
                        bind:value={newAnimal.Image_src}
                        type="file"
                        accept=".jpeg, .jpg, .png, .bmp, .webp, .avif, .svg"
                  />
            </div>
            <hr />
            <h3>Shots</h3>
            <div class="shots">
                  {#each newAnimal.Shots as shot, i}
                        <div class="shot-wrapper">
                              <div>
                                    <label for="shot-name">Name</label>
                                    <select
                                          name="shot-name"
                                          id="sahot-name"
                                          bind:value={newAnimal.Shots[i].Id}
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
                                          bind:value={newAnimal.Shots[i]
                                                .Date_given}
                                    />
                              </div>
                              <div>
                                    <label for="shot-due">Next Due</label>
                                    <input
                                          type="date"
                                          name="shot-due"
                                          bind:value={newAnimal.Shots[i]
                                                .Next_due}
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
