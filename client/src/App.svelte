<script lang="ts">
  import { cardResponse, getRandomCard } from "./helpers/card";
  import { trapFocus } from "./helpers/actions.svelte";
  import { onMount } from "svelte";
  import CardComponent from "./components/CardComponent.svelte";

  type Animal = {
    id: string;
    name: string;
    species: string;
    date_of_birth: string;
    sex: string;
    available: boolean;
    image_src: string;
  };

  type Card = {
    id: string;
    animal_id: string;
    organization_id: string;
    image_src: string;
    liked: boolean;
    animal_info: Animal;
  };

  let isLiked = $state<boolean | undefined>();
  let cardInfo = $state<Card | undefined>(undefined);
  let isLoading = $state(true);

  async function newCard() {
    cardInfo = await getRandomCard();
    isLoading = false;
  }

  onMount(newCard);

  async function likedCard() {
    isLiked = true;
    await cardResponse(isLiked, "001");
    newCard();
  }

  async function dislikedCard() {
    isLiked = false;
    await cardResponse(isLiked, "001");
    newCard();
  }
</script>

{#if isLoading}
  <div>Loading...</div>
{:else}
  <div class="content" use:trapFocus>
    <CardComponent card_info={cardInfo} {isLiked} />
    <div class="liked-buttons">
      <button onclick={() => dislikedCard()}>&lt</button>
      <button onclick={() => likedCard()}>&gt</button>
    </div>
    <button onclick={() => console.log($state.snapshot(isLiked))}
      >State of isLiked</button
    >
  </div>
{/if}

<style>
  .content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .liked-buttons {
    display: flex;
    gap: 5%;
  }
  .liked-buttons button {
    width: 47.5%;
  }
</style>
