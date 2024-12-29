<script lang="ts">
	let formFields = [
		{ name: 'Username', type: 'text' },
		{ name: 'Email', type: 'email' },
		{ name: 'Password', type: 'password' },
		{ name: 'Confirm Password', type: 'password' } //might have issues with the name of this one
	];
</script>

<main>
	<h2>Create Account</h2>

	<form>
		<div class="input-grandparent">
			{#each formFields as field}
				<div class="input-parent">
					<label for={field.name.toLowerCase()}>{field.name}</label>
					<input type={field.type} id={field.name.toLowerCase()} name={field.name.toLowerCase()} />
				</div>
			{/each}
		</div>

		<button class="action-button">Create</button>
	</form>
	<p>OR</p>
	<div class="org-div">
		<button onclick={() => (window.location.href = 'google.com')} class="action-button org-button"
			>Create Organization</button
		>
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
