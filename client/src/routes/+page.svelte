<script>
	import Home from '../components/home.svelte';
	import { GetTokens } from '../helpers/users';

	$effect(() => {
		GetTokens();
	});
</script>

<main>
	<Home />
</main>

<style>
	main {
		height: 100%;
	}
</style>
