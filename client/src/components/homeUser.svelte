<script lang="ts">
	import { cardResponse, getRandomCard } from '../helpers/card';
	import { goto } from '$app/navigation';
	import { trapFocus } from '../helpers/actions.svelte';
	import { onMount } from 'svelte';
	import CardComponent from './CardComponent.svelte';
	import Navbar from './Navbar.svelte';

	type Animal = {
		Id: string;
		Name: string;
		Species: string;
		Date_of_birth: string;
		Sex: string;
		Available: boolean;
		Price: number;
		Breed: string;
		Image_src: string[];
	};

	let isLiked = $state<boolean | undefined>();
	let cardInfo = $state<Animal | undefined>(undefined);
	let isLoading = $state(true);
	let cardDone = $state(false);

	async function newCard() {
		let randomCardInfo = await getRandomCard();
		cardInfo = randomCardInfo.data;
		isLoading = false;
		isLiked = undefined;
		cardDone = false;
	}

	onMount(async () => {
		await newCard();
	});

	async function likedCard() {
		console.log($state.snapshot(cardInfo))
		isLiked = true;
		if (cardInfo?.Id) {
			let statusCode = await cardResponse(isLiked, cardInfo?.Id);
		} else {alert("animal does not have an id!")}
		setTimeout(async () => {
			cardDone = true;
			await newCard();
		}, 1300);
	}

	async function dislikedCard() {
		isLiked = false;
		if (cardInfo?.Id) {
			let statusCode = await cardResponse(isLiked, cardInfo?.Id);
		} else {alert("animal does not have an id!")}
		setTimeout(async () => {
			cardDone = true;
			await newCard();
		}, 1300);
	}

</script>

{#if isLoading}
	<div>Loading...</div>
{:else}
	<main>
		<Navbar />
		<div class="content" use:trapFocus>
			<CardComponent card_info={cardInfo} {isLiked} {cardDone} />
			<div class="liked-buttons">
				<button onclick={() => dislikedCard()}>&lt</button>
				<button onclick={() => likedCard()}>&gt</button>
			</div>
		</div>
	</main>
{/if}

<style>
	main {
		display: grid;
		grid-template-columns: 17.25rem 1fr;
		height: 100%;
		gap: 5rem;
	}
	.content {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: center;
		justify-content: center;
	}

	.liked-buttons {
		display: flex;
		gap: 2rem;
	}
	.liked-buttons button {
		width: 14rem;
		height: 3rem;
	}
</style>
