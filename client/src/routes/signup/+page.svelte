<script lang="ts">
	import { goto } from '$app/navigation';
	import { CreateUser } from '../../helpers/users';
	import {ValidatePassword} from "../../helpers/validation"
	import type {StrongPassword} from "../../helpers/validation"

	let invalidPassword: boolean = $state(false)
	let showPassword: boolean = $state(false)
	let showCPassword: boolean = $state(false)
	
	let formData = $state({
		Username: '',
		Email: '',
		Password: '',
		C_password: '',
		Is_organization: false
	});
	let passwordChecker: StrongPassword = $state({symbol:false, uppercase:false,lowercase:false, length: formData.Password.length, isLength: false, number:false})

	async function submitForm(e: Event) {
		e.preventDefault();
		if (await CreateUser(formData)) {
			goto('/');
		}
	}

	function checkPassword(password: string) {
		for(const pass of Object.values(ValidatePassword(password))) {
			if (pass === false) {
				console.log(`${pass} is false`)
				invalidPassword = true
			} 
		}
		invalidPassword = false
	}
</script>

<!-- how to get classes to update each time passwordChcker changes -->
<div class="invalid-password-container"> 
	<li>
		<ol class="lowercase">
			<div class:valid-option={passwordChecker.lowercase} class:invalid-option={!passwordChecker.lowercase}></div> 
			Password has lowercase characters</ol>
		<ol class="uppercase">
			<div class:valid-option={passwordChecker.uppercase} class:invalid-option={!passwordChecker.uppercase}></div>
			Password has uppercase characters</ol>
		<ol class="numbers">
			<div  class:valid-option={passwordChecker.number} class:invalid-option={!passwordChecker.number}></div>
			Password has numbers</ol>
		<ol class="special">
			<div  class:valid-option={passwordChecker.symbol} class:invalid-option={!passwordChecker.symbol}></div>
			Password has special characters </ol>
		<ol class="length">
			<div  class:valid-option={passwordChecker.length >= 6} class:invalid-option={passwordChecker.length < 6}></div>
			Password is atleast 6 characters</ol>
	</li>
</div>

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
				<div class="password-box">
					<input bind:value={formData.Password} type={!showPassword ? "password" : "text"} id="password" name="password" onchange={() => checkPassword(formData.Password)}/>
					<button type="button" onclick={() => showPassword = !showPassword}>
						{#if !showPassword} 
						<img width="24" height="24" src="https://img.icons8.com/material-outlined/24/visible--v1.png" alt="show"/>
						{:else}
						<img width="24" height="24" src="https://img.icons8.com/material-outlined/24/hide.png" alt="hide"/>
						{/if}
					</button>
				</div>
			</div>
			<div class="input-parent">
				<label for="confirm-password">Confirm Password</label>
				<div class="password-box">
					<input bind:value={formData.C_password} type={!showCPassword ? "password" : "text"} id="confirm-password" name="confirm-password" onchange={() => checkPassword(formData.C_password)}/>
					<button type="button" onclick={() => showCPassword = !showCPassword}>
						{#if !showCPassword} 
						<img width="24" height="24" src="https://img.icons8.com/material-outlined/24/visible--v1.png" alt="show"/>
						{:else}
						<img width="24" height="24" src="https://img.icons8.com/material-outlined/24/hide.png" alt="hide"/>
						{/if}
					</button>
				</div>
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

	.password-box{
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		position: relative;
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

	.invalid-password-container li {
		display: grid;
		grid-template-columns: 20rem 20rem;
		grid-template-rows:5rem 5rem;
		align-items: center;

	}
	.invalid-password-container li ol {
		display: flex;
		flex-direction: row;
		align-items: center;
		text-align: center;
		gap: .5rem;
	}
	.invalid-password-container li ol div {
		width: .5rem;
		height: .5rem;
		border-radius: 50%;
		border: none;
		background: rgb(75, 75, 75);
	}
	.invalid-option {
		background-color: rgb(165, 54, 47) !important;
	}	
	.valid-option {
		background-color: rgb(86, 179, 86) !important;
	}
</style>
