<script lang="ts">
	import { goto } from '$app/navigation';
	import { loginUser } from '../../helpers/users';

	let formData = $state({
		Username: '',
		Password: ''
	});
	let showPassword:boolean = $state(false)

	let logos = [
		{ name: 'google', path: '/icons/Google.svg' },
		{ name: 'github', path: '/icons/Github.svg' },
		{ name: 'facebook', path: '/icons/Facebook.svg' }
	];

	async function formSubmit(e: Event) {
		e.preventDefault();
		if (await loginUser(formData)) {
			goto('/');
		}
	}
</script>

<main>
	<h2>Sign In</h2>

	<form onsubmit={formSubmit}>
		<div class="input-grandparent">
			<div class="input-parent">
				<label for="username">Username</label>
				<input bind:value={formData.Username} type="text" id="username" name="username" />
			</div>

			<div class="input-parent">
				<div class="password-box">
					<input bind:value={formData.Password} type={!showPassword ? "password" : "text"} id="password" name="password"/>
					<button type="button" onclick={() => showPassword = !showPassword}>
						{#if !showPassword} 
						<img width="24" height="24" src="https://img.icons8.com/material-outlined/24/visible--v1.png" alt="show"/>
						{:else}
						<img width="24" height="24" src="https://img.icons8.com/material-outlined/24/hide.png" alt="hide"/>
						{/if}
					</button>
				</div>
			</div>
		</div>

		<button class="action-button">Sign In</button>
	</form>
	<p>OR</p>
	<button onclick={() => goto('/signup')} class="action-button">Create Account</button>

	<div class="logo-parent">
		{#each logos as logo}
			<img src={logo.path} alt={`sign in with ${logo.name}`} />
		{/each}
	</div>
</main>

<style>
	main {
		height: 100%;
		display: flex;
		flex-direction: column;
		/* background: pink; */
		padding: 0;
		margin: 0;
		text-align: center;
		align-items: center;
		/* gap: 1.5rem; */
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

	.logo-parent {
		display: flex;
		gap: 3rem;
		padding-top: 1.5rem;
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

	img {
		width: 3rem;
		height: 3rem;
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

	.password-box{
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		position: relative;
		width:22rem;
	
	}

	.password-box button {
		border: none;

		background: transparent;
		width: fit-content;
		height: fit-content;
		position: absolute;
		right: 1rem;
		top: .5rem;
	}

	.password-box button img {
		width: 2rem;
		height: 2rem;
		position: relative;
		top: .25rem;
	}
</style>
