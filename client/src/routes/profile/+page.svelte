<script lang="ts">
	import type { incomingUser } from '../../helpers/users';
	import { GetCurrentUser, LogoutUser, UpdateUser } from '../../helpers/users';
	import { onMount } from 'svelte';
	import { formatISO } from 'date-fns';
	import Navbar from '../../components/Navbar.svelte';

	let userData: incomingUser | null = $state(null);
	let profilePicture: FileList | undefined = $state();
	let updatedUserData = $state({
		Id: '',
		Username: '',
		Email: '',
		Date_of_birth: ''
	});
	let options = $state({
		showEdit: false
	});
	let loading = $state(true);

	onMount(async () => {
		userData = await GetCurrentUser();
		updatedUserData = {
			Id: userData?.Id ? userData?.Id : '',
			Username: userData?.Username ? userData?.Username : '',
			Email: userData?.Email ? userData?.Email : '',
			Date_of_birth: userData?.Date_of_birth ? userData.Date_of_birth.split('T')[0] : ''
		};
		loading = false;
	});

	async function submitForm(e: Event) {
		e.preventDefault();
		const formattedUser = {
			...$state.snapshot(updatedUserData),
			Date_of_birth: updatedUserData.Date_of_birth
				? formatISO(new Date(updatedUserData.Date_of_birth))
				: ''
		};
		let statusCode = await UpdateUser(formattedUser);

		if (profilePicture != undefined && userData != null) {
			const formData = new FormData();
			formData.append('profile_image', profilePicture[0]);
			try {
				const response = await fetch(`/api/users/profile_pictures/${userData.Id}`, {
					method: 'POST',
					headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
					body: formData
				});
				if (!response.ok) {
					const data = await response.json();
					profilePicture = undefined;
					throw new Error(`failed to upload profile picture, error: ${data.error}`);
				}
				console.log('video uploaded');
				userData = await GetCurrentUser();
				profilePicture = undefined;
			} catch (error: any) {
				alert(`Error: ${error.message}`);
			}
		}

		if (statusCode == 200) {
			console.log('success');
			location.reload();
		}
	}
</script>

{#if loading}
	<p>Loading...</p>
{:else if userData != null}
	<div class="content">
		<Navbar />
		<main>
			<h1>Profile</h1>
			<div class="user-info">
				<p>Hello {userData.Username}!</p>
				<!-- svelte-ignore a11y_missing_attribute -->
				<img src={userData.Profile_picture} />
			</div>

			<button onclick={() => (options.showEdit = !options.showEdit)}>Edit</button>
			<button onclick={LogoutUser}>Logout</button>

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
						{#if !userData.Is_organization}
						<div class="input-holder">
							<label for="dob">Date of Birth</label>
							<input
								type="date"
								placeholder="date of birth"
								bind:value={updatedUserData.Date_of_birth}
							/>
						</div>
						{/if}
						<div class="input-holder">
							<label for="profile-picture">Profile Picture</label>
							<input
								type="file"
								class="profile-picture"
								name="profile-picture"
								multiple={false}
								accept="image/*"
								bind:files={profilePicture}
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
	</div>
{/if}

<style>
	main {
		height: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 0px;
		margin: 0px;
		gap: 2rem;
	}
	.content {
		display: grid;
		grid-template-columns: 17.25rem 1fr;
		height: 100%;
		gap: 5rem;
	}

	form {
		position: absolute;
		top: 20rem;
		display: flex;
		flex-direction: column;
		gap: 2rem;
		align-items: center;
	}
	img {
		width: 5rem;
		height: 5rem;
		border-radius: 50%;
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 2rem;
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
	.profile-picture {
		border: none !important;
		font-size: 12px;
	}
	button {
		width: 10rem;
		height: 1.5rem;
	}
</style>
