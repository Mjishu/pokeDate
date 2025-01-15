<script lang="ts">
	import { goto } from '$app/navigation';
	import { CreateUser } from '../../helpers/users';

	let formData = {
		Username: '',
		Email: '',
		Password: '',
		C_password: '',
		Is_organization: false
	};

	async function submitForm(e: Event) {
		e.preventDefault();
		if (await CreateUser(formData)) {
			goto('/');
		}
	}
</script>

<main>
	<h2>Create Account</h2>

	<form onsubmit={submitForm}>
		<div class="input-grandparent">
			<div class="input-parent">
				<label for="username">Username</label>
				<input bind:value={formData.Username} type="text" id="username" name="username" />
			</div>
			<div class="input-parent">
				<label for="email">Email</label>
				<input bind:value={formData.Email} type="email" id="email" name="email" />
			</div>
			<div class="input-parent">
				<label for="password">Password</label>
				<input bind:value={formData.Password} type="password" id="password" name="password" />
			</div>
			<div class="input-parent">
				<label for="confirm-password">Confirm Password</label>
				<input
					bind:value={formData.C_password}
					type="password"
					id="confirm-password"
					name="confirm-password"
				/>
			</div>
			<div class="input-paren">
				<label for="is_org">Are you an Organization?</label>
				<input type="checkbox" name="is_org" id="is_org" bind:checked={formData.Is_organization}>
			</div>
		</div>

		<button class="action-button">Create</button>
	</form>
	<p>OR</p>
	<div class="org-div">
		<button onclick={() => (window.location.href = 'google.com')} class="action-button org-button">
			Create Organization
		</button>
	</div>
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

	.org-button {
		background: transparent;
		border: none;
		color: #96a0de;
		transition: ease-in-out 400ms color;
	}

	.org-button:hover {
		background: none;
		color: #747caf;
	}

	.org-div {
		display: flex;
		width: 18rem;
		border-bottom: 1px solid #96a0de;
		color: #96a0de;
		align-items: center;
		opacity: 0.8;
		transition:
			ease-in-out 400ms border,
			ease-in-out 300ms opacity;
	}

	.org-div:hover {
		opacity: 1;
		border-bottom: 1px solid #747caf;
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

	p {
		font-size: 2rem;
		font-weight: 100;
	}
</style>
