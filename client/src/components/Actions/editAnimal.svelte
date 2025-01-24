<!-- TODO  -->
<!-- Iterate over each updatedAnimal.shots so that if we get info from the backend about a shot it will show the data -->
<script lang="ts">
	import { createAnimal, getAnimalById, updateAnimalById, GetAllShots } from '../../helpers/animals';
	import type { NewAnimal, UpdatedAnimal, Shot, Animal } from "../../helpers/animals"
	import { formatISO } from 'date-fns';
	import { onMount } from 'svelte';

	let { showEditPage = $bindable(), currentId } = $props();
	let oldAnimal = $state();
	let animal = $state<Animal>();
	let shotData = $state<Shot[]>();

	let Image_src = $state<string>();
	let updatedAnimal = $state<UpdatedAnimal>({
		// switch this to animal.name etc
		Name: '',
		Date_of_birth: '',
		Price: 0,
		Available: false,
		Shots: []
	});

	function closeForm() {
		showEditPage = false;
	}

	function addNewShot() {
		updatedAnimal.Shots.push({
			Id: 0,
			Date_given: '',
			Next_due: ''
		});
	}
	onMount(async () => { //! Not calling correct route!
		shotData = await GetAllShots();
		animal = await getAnimalById(currentId);

		updatedAnimal = {
			Name: animal?.Name || '',
			Date_of_birth: animal?.Date_of_birth ? animal.Date_of_birth.split('T')[0] : '',
			Price: animal?.Price || 0,
			Available: animal?.Available || false,
			Shots:
				animal?.Shots?.map((shot) => ({
					Id: shot.Id || shot.Id,
					Date_given: shot.Date_given ? shot.Date_given.split('T')[0] : '',
					Next_due: shot.Next_due ? shot.Next_due.split('T')[0] : ''
				})) || []
		};
	});

	async function cleanAndUpdateAnimal(e: Event) {
		e.preventDefault();

		const formattedAnimal = {
			...$state.snapshot(updatedAnimal),
			Date_of_birth: updatedAnimal.Date_of_birth
				? formatISO(new Date(updatedAnimal.Date_of_birth))
				: undefined,
			Shots: updatedAnimal?.Shots?.map((shot) => ({
				Id: Number(shot.Id),
				Date_given: shot.Date_given ? formatISO(new Date(shot.Date_given)) : '',
				Next_due: shot.Next_due ? formatISO(new Date(shot.Next_due)) : ''
			}))
		};

		await updateAnimalById(currentId, formattedAnimal);
		console.log(formattedAnimal);
		closeForm();
	}

	async function RemoveShot() {
		console.log("removing shot")
	}
</script>

<main>
	<form onsubmit={cleanAndUpdateAnimal} autocomplete="off">
		<h3>Information</h3>
		<div>
			<label for="name">Name</label>
			<input required type="text" name="name" bind:value={updatedAnimal.Name} />
		</div>
		<div>
			<label for="date_of_birth">Date of Birth</label>
			<input type="date" name="date_of_birth" bind:value={updatedAnimal.Date_of_birth} />
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
				bind:value={updatedAnimal.Price}
			/>
		</div>
		<div>
			<label for="available">Available</label>
			<input type="checkbox" name="available" bind:checked={updatedAnimal.Available} />
		</div>
		<hr />
		<h3>Images</h3>
		<div class="image-container">
			<input
				multiple={false}
				bind:value={Image_src}
				type="file"
				accept=".jpeg, .jpg, .png, .bmp, .webp, .avif, .svg"
			/>
		</div>
		<hr />
		<h3>Shots</h3>
		<div class="shots">
			{#each updatedAnimal.Shots as shot, i}
				<div class="shot-wrapper">
					<div>
						<label for="shot-name">Name</label>
						<select name="shot-name" id="sahot-name" bind:value={updatedAnimal.Shots[i].Id}>
							<option value="" disabled selected
								>{updatedAnimal.Shots[i].Name != '' ? updatedAnimal.Shots[i].Name : 'Name'}</option
							>
							{#if shotData != undefined}
								{#each shotData as shot}
									<option value={shot.Id}>{shot.Name}</option>
								{/each}
							{/if}
						</select>
					</div>
					<div>
						<label for="shot-given">Shot Given</label>
						<input type="date" name="shot-given" bind:value={updatedAnimal.Shots[i].Date_given} />
					</div>
					<div>
						<label for="shot-due">Next Due</label>
						<input type="date" name="shot-due" bind:value={updatedAnimal.Shots[i].Next_due} />
					</div>
					<button class="delete-shot" type="button" onclick={RemoveShot}>Remove</button>
				</div>
			{/each}
			<button type="button" onclick={addNewShot}>+</button>
		</div>
		<hr />
		<button type="submit">Edit</button>
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
