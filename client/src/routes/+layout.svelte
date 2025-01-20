<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { incomingUser } from '../helpers/users.js';
	import { GetCurrentUser } from '../helpers/users.js';
	import Notifications from '../components/notifications.svelte';
	
	let { children } = $props();

	let userData: incomingUser | null = $state(null);
	let signedIn: boolean = $state(false);
	let loading: boolean = $state(true);

	let showNotifications: boolean = $state(false);

	onMount(async () => {
		userData = await GetCurrentUser();
		if (userData != null) {
			signedIn = true
			if (window.location.pathname == "/login" || window.location.pathname == "/signup") {
				goto("/")
			}
		} else if (userData == null) {
			if (window.location.pathname != "/login" && window.location.pathname != "/signup") {
				goto("/login")
			}
		}
		loading = false;
	});

	function pfpClick() {
		goto('/profile');
	}
</script>

{#if loading}
	<p>loading...</p>
{:else if signedIn}
	<button class="pfp-holder" onclick={pfpClick}>
		<!-- svelte-ignore a11y_missing_attribute -->
		<img class="pfp" src={userData?.Profile_picture} />
	</button>
{/if}

<button class="notification-button" onclick={() => showNotifications = !showNotifications}><img class="notification-bell" src="/icons/notification_bell.svg" alt="notification"></button>

{#if showNotifications}
<div class="notification">
	<Notifications/>
	<button type="button" class="close-notifications" onclick={() => showNotifications = false}>Close</button>
</div>
{/if}

{@render children()}

<!-- <slot {userData} /> -->

<style>
	.pfp {
		/* z-index: 3; */
		position: relative;
		width: 3rem;
		height: 3rem;
		border-radius: 50%;
		padding: 0;
		margin: 0;
	}

	.pfp-holder {
		position: absolute;
		top: 1rem;
		right: 1rem;
		width: 3rem;
		height: 3rem;
		padding: 0;
		margin: 0;
		border: none;
		background-color: transparent;

	}

	.notification {
		position: absolute;
		left: 50vh;
		top: 50vh;
		background: rgb(218, 160, 189);
		width: 20rem;
		height: 25rem;
		border-radius: 10px;
		display: flex;
		flex-direction: column;
		align-items: center;
		padding-bottom: 3rem;
		justify-content: end;
		overflow-y: auto;
	}

	.notification-button {
		position: absolute;
		bottom: 3rem;
		background-color: transparent;
		border: none;
		width: fit-content;
		height: fit-content;
	}
	.notification-bell {
		width: 3rem;
		height: 3rem;
	}
	.close-notifications {
		background: white;
		border: none;
		border-radius: 10px;
		height: 3rem;
		width: 5rem;
	}
</style>
