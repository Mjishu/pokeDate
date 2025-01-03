<script lang="ts">
	import type { incomingUser } from '../../helpers/users';

	export let userData: incomingUser;
	console.log(userData);
</script>

<main>
	<h1>Profile</h1>
	<p>Hello {userData.Username}!</p>
</main>

<style>
	main {
		height: 100%;
	}
</style>
