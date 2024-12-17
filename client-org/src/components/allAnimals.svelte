<script>
      // @ts-nocheck

      import { getOrganizationAnimals } from "../helper/animals";
      import { format } from "date-fns";
      import EditAnimal from "./editAnimal.svelte";

      let animalsData = $state();
      let isLoading = $state(true);
      let currentId = $state();
      let showEditPage = $state(false);

      async function callAnimals() {
            animalsData = await getOrganizationAnimals();
            isLoading = false;
      }
      callAnimals();

      function showEdit(id) {
            showEditPage = true;
            currentId = id;
      }
</script>

<!-- get all animals -> make sure it is from X org(the one that is logged in) -->
<!--* display information in animalsData -->
{#if isLoading}
      <div>Loading...</div>
{:else}
      <div class="show-animals">
            <table>
                  <thead class="animal-headers">
                        <tr>
                              <th>Name</th>
                              <th>Species</th>
                              <th>Breed</th>
                              <th>Date of Birth</th>
                              <th>Sex</th>
                              <th>Price</th>
                              <th>Available</th>
                              <th>Actions</th>
                        </tr>
                  </thead>
                  <tbody>
                        {#each animalsData as animal}
                              <tr class="animal-data">
                                    <td>{animal.Name}</td>
                                    <td>{animal.Species}</td>
                                    <td>{animal.Breed}</td>
                                    <td
                                          >{format(
                                                animal.Date_of_birth,
                                                "yyyy/MM/dd",
                                          )}</td
                                    >
                                    <td>{animal.Sex}</td>
                                    <td>{animal.Price}</td>
                                    <td>{animal.Available}</td>
                                    <td class="action-holder">
                                          <button
                                                onclick={() =>
                                                      showEdit(animal.Id)}
                                                >edit</button
                                          >
                                          <button
                                                ><img
                                                      src="icons/trash_icon.svg"
                                                      alt=""
                                                /></button
                                          >
                                    </td>
                              </tr>
                        {/each}
                  </tbody>
            </table>
            {#if showEditPage}
                  <EditAnimal {currentId} bind:showEditPage />
            {/if}
      </div>
{/if}

<style>
      .show-animals {
            border-radius: 5px;
            border: none;
            background: hsla(0, 1%, 77%, 0.5);
            width: 95%;
            height: 58rem;
            overflow-y: scroll;
      }
      table thead tr {
            background-color: #009879; /* Header background color */
            color: #ffffff; /* Header text color */
            text-align: left; /* Align header text */
            font-weight: bold; /* Make header bold */
      }

      table th,
      table td {
            padding: 12px 15px; /* Add padding inside cells */
            border: 1px solid #dddddd; /* Add borders */
      }

      table tbody tr {
            border-bottom: 1px solid #dddddd; /* Add border between rows */
      }

      table tbody tr:nth-of-type(even) {
            background-color: #f3f3f3; /* Alternate row color */
      }

      table tbody tr:hover {
            background-color: #f1f1f1; /* Row hover effect */
      }

      table tbody tr:last-of-type {
            border-bottom: 2px solid #009879; /* Stronger border for last row */
      }
      table {
            width: 100%;
            border-collapse: collapse;
            margin: 0 0;
            font-size: 16px;
            text-align: left;
      }

      .action-holder {
            display: flex;
            align-items: center;
            text-align: center;
            margin: 0;
            gap: 0.5rem;
            justify-content: center;
      }

      .action-holder button {
            border: none;
            outline: none;
            background-color: transparent;
      }

      .action-holder button img {
            width: 1rem;
            height: 1rem;
      }
</style>
