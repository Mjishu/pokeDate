<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Home from '../components/home.svelte';
	import { GetTokens } from '../helper/auth';

	onMount(async () => {
		const status = await GetTokens();
		if (status.statusCode == 400) {
			goto('/login');
		}
	});
</script>

<main>
	<Home />
</main>

<style>
	main {
		height: 100%;
		margin: 0;
		padding: 0;
	}
</style>
