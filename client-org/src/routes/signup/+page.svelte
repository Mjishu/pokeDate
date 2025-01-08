<script lang="ts">
	import { goto } from '$app/navigation';
	import { CreateOrganization } from '../../helper/auth';

	let passwordValid: boolean = $state(true);

	let formData = $state({
		Name: '',
		Email: '',
		Password: '',
		C_password: ''
	});

	async function submitForm(e: Event) {
		e.preventDefault();
		if (passwordValid === true) {
			if (await CreateOrganization(formData)) {
				goto('/');
			}
		} else {
			alert('fix your form');
		}
	}

	function checkPasswords() {
		if (formData.Password !== formData.C_password) {
			if (formData.C_password.length == 0) {
				passwordValid = true;
				return;
			}
			passwordValid = false;
			return;
		}
		passwordValid = true;
	}
</script>

<main>
	<h2>Create Organization</h2>

	<form onsubmit={submitForm}>
		<div class="input-grandparent">
			<div class="input-parent">
				<label for="username">Organization Name</label>
				<input bind:value={formData.Name} type="text" id="username" name="username" />
			</div>
			<div class="input-parent">
				<label for="email">Email</label>
				<input bind:value={formData.Email} type="email" id="email" name="email" required />
			</div>
			<div class="input-parent">
				<label for="password">Password</label>
				<input
					bind:value={formData.Password}
					type="password"
					id="password"
					name="password"
					required
				/>
			</div>
			<div class="input-parent">
				<label for="confirm-password">Confirm Password</label>
				<input
					bind:value={formData.C_password}
					type="password"
					id="confirm-password"
					name="confirm-password"
					required
					class={passwordValid ? 'valid-password' : 'invalid-password'}
					onchange={checkPasswords}
				/>
				{#if !passwordValid}
					<p>passwords do not match</p>
				{/if}
			</div>
		</div>

		<button class="action-button">Create</button>
	</form>
</main>

<style>
	/* Your original styles here */
	main {
		height: 100%;
		display: flex;
		flex-direction: column;
		padding: 0;
		margin: 0;
		text-align: center;
		align-items: center;
		padding-top: 7rem;
	}

	.input-parent {
		display: flex;
		flex-direction: column;
		text-align: left;
	}
	.input-parent input {
		width: 22rem;
		height: 3.75rem;
		border-radius: 10px;
		font-family: inherit;
		border: 1px solid rgb(148, 148, 148);
		font-size: 2rem;
		font-weight: 100;
	}

	.input-grandparent {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.action-button {
		background-color: #fff;
		opacity: 0.9;
		transition:
			ease-in-out 200ms background-color,
			ease-in 200ms opacity,
			ease-in-out 450ms color;
	}

	.action-button:hover {
		background-color: rgb(68, 121, 68);
		opacity: 1;
		color: #fff;
	}

	h2 {
		font-size: 4rem;
		letter-spacing: 10%;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 2.675rem;
	}

	button {
		width: 22rem;
		height: 3.75rem;
		font-size: 2rem;
		font-weight: 100;
		background-color: #fff;
		border-radius: 5px;
		border: 1px solid rgb(185, 185, 185);
	}

	.invalid-password {
		border: 1px solid red !important;
	}
</style>
