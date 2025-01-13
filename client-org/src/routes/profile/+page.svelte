<script lang="ts">
	import { onMount } from 'svelte';
	import { formatISO } from 'date-fns';
	import Navbar from '../../components/navbar.svelte';
	import {
		GetCurrentOrganization,
		LogoutOrganization,
		UpdateOrganization
	} from '../../helper/auth';

	type incomingOrganization = {
		Id?: string;
		Name: string;
		Email: string;
		Profile_picture?: string;
	};

	let profilePicture: FileList | undefined = $state();
	let orgData: incomingOrganization | null = $state(null);
	let updatedOrgData = $state({
		Id: '',
		Name: '',
		Email: ''
	});
	let options = $state({
		showEdit: false
	});
	let loading = $state(true);

	onMount(async () => {
		orgData = await GetCurrentOrganization();
		updatedOrgData = {
			Id: orgData?.Id ? orgData?.Id : '',
			Name: orgData?.Name ? orgData?.Name : '',
			Email: orgData?.Email ? orgData?.Email : ''
		};
		loading = false;
	});

	async function submitForm(e: Event) {
		e.preventDefault();

		let statusResponse = await UpdateOrganization(updatedOrgData);

		if (profilePicture != undefined && orgData != null) {
			const formData = new FormData();
			formData.append('profile_image', profilePicture[0]);
			try {
				const response = await fetch(`/api/users/profile_pictures/${orgData.Id}`, {
					method: 'POST',
					headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
					body: formData
				});
				if (!response.ok) {
					const data = await response.json();
					profilePicture = undefined;
					throw new Error(`failed to upload profile picture, error: ${data.error}`);
				}
				orgData = await GetCurrentOrganization();
				profilePicture = undefined;
			} catch (error: any) {
				alert(`Error: ${error.message}`);
			}
		}

		if (statusResponse == 200) {
			location.reload();
		}
	}
</script>

{#if loading}
	<p>Loading...</p>
{:else if orgData != null}
	<div class="content">
		<Navbar />
		<main>
			<h1>Profile</h1>
			<div class="user-info">
				<p>Hello {orgData.Name}!</p>
				<!-- svelte-ignore a11y_missing_attribute -->
				<img src={orgData?.Profile_picture} />
			</div>

			<button onclick={() => (options.showEdit = !options.showEdit)}>Edit</button>
			{#if options.showEdit}
				<form onsubmit={submitForm} autocomplete="off">
					<div class="inputs">
						<div class="input-holder">
							<label for="username">Organization Name</label>
							<input type="text" placeholder="username" bind:value={updatedOrgData.Name} />
						</div>
						<div class="input-holder">
							<label for="email">Email</label>
							<input type="text" placeholder="email" bind:value={updatedOrgData.Email} />
						</div>
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
			<button type="button" onclick={LogoutOrganization}>Logout</button>
		</main>
	</div>
{/if}

<style>
	main {
		height: 100%;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
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
