<script lang="ts">
	import type { incomingUser } from '../../helpers/users';
	import { GetCurrentUser, UpdateUser } from '../../helpers/users';
	import { onMount } from 'svelte';

	let userData: incomingUser | null = $state(null);
	let updatedUserData = $state({
		Username: '',
		Email: '',
		Date_of_birth: null as Date | null
	});
	let options = $state({
		showEdit: false
	});

	onMount(async () => {
		userData = await GetCurrentUser();
		updatedUserData = {
			Username: userData?.Username ? userData?.Username : '',
			Email: userData?.Email ? userData?.Email : '',
			Date_of_birth: userData?.Date_of_birth ? userData?.Date_of_birth : null
		};
	});

	async function submitForm(e: Event) {
		e.preventDefault();
		await UpdateUser(updatedUserData);
	}
</script>

<main>
	<h1>Profile</h1>
	<p>Hello {userData?.Username}!</p>

	<button onclick={() => (options.showEdit = !options.showEdit)}>Edit</button>
	{#if options.showEdit}
		<form onsubmit={submitForm} autocomplete="off">
			<div class="inputs">
				<div class="input-holder">
					<label for="username">New Username</label>
					<input type="text" placeholder="username" bind:value={updatedUserData.Username} />
				</div>
				<div class="input-holder">
					<label for="email">Email</label>
					<input type="text" placeholder="email" bind:value={updatedUserData.Email} />
				</div>
				<div class="input-holder">
					<label for="dob">Date of Birth</label>
					<input
						type="text"
						placeholder="date of birth"
						bind:value={updatedUserData.Date_of_birth}
					/>
				</div>
			</div>

			<div class="button-holder">
				<button type="button" onclick={() => (options.showEdit = false)}>Cancel</button>
				<button type="submit">Submit</button>
			</div>
		</form>
	{/if}
</main>

<style>
	main {
		height: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 0px;
		margin: 0px;
	}

	form {
		position: absolute;
		top: 20rem;
		display: flex;
		flex-direction: column;
		gap: 2rem;
		align-items: center;
	}
	.inputs {
		display: grid;
		gap: 2rem;
		grid-template-columns: repeat(2, 1fr);
	}
	.input-holder {
		display: flex;
		flex-direction: column;
	}
	.input-holder input {
		width: 10rem;
		height: 1.5rem;
		border: 1px solid rgb(119, 119, 119);
		border-radius: 5px;
	}
	button {
		width: 10rem;
		height: 1.5rem;
	}
</style>
