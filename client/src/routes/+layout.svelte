<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { incomingUser } from '../helpers/users.js';

	let userData: incomingUser | null = $state(null);

	/**
	 * this shoudl work?? it calls gettokens on each page change so why can i not just pass the data returned from getcurUser????
	 */

	async function GetTokens(): Promise<void> {
		console.log('get tokens called');
		try {
			const refreshToken = localStorage.getItem('refresh_token');
			const bearerToken = 'Bearer ' + refreshToken;
			if (!refreshToken) {
				console.log('i dont have a refresh token. log in');
				return;
			}
			const fetchParams = {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: bearerToken
				}
			};

			const response = await fetch('/api/refresh', fetchParams);
			const data = await response.json();
			if (data.token) {
				localStorage.setItem('token', data.token);
			}
		} catch (error) {
			console.error(`error fetching tokens ${error}`);
			return;
		}
	}

	async function GetCurrentUser(): Promise<incomingUser | null> {
		try {
			const token = localStorage.getItem('token');
			const bearerToken = 'Bearer ' + token;
			const fetchParams = {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: bearerToken
				}
			};
			const response = await fetch('/api/users/current', fetchParams);
			const data = await response.json();
			if (response.status == 200) {
				return data;
			} else {
				return null;
			}
		} catch (err) {
			console.error(`error fetching curernt user data ${err}`);
			return null;
		}
	}

	onMount(async () => {
		await GetTokens();
		userData = await GetCurrentUser();
		// console.log(userData);
	});

	$effect(() => {
		GetTokens();
	});
</script>

<slot {userData} />
