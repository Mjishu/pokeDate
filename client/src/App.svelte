<script lang="ts">
  import { cardResponse, getRandomCard } from "./helpers/card";
  import { trapFocus } from "./helpers/actions.svelte";
  import { onMount } from "svelte";
  import CardComponent from "./components/CardComponent.svelte";

  type Animal = {
    Id: string;
    Name: string;
    Species: string;
    Date_of_birth: string;
    Sex: string;
    Available: boolean;
    Price: number;
    Breed: string;
    Image_src: string;
  };

  // type Card = {
  //   id: string;
  //   animal_id: string;
  //   organization_id: string;
  //   image_src: string;
  //   liked: boolean;
  //   animal_info: Animal;
  // };

  let isLiked = $state<boolean | undefined>();
  let cardInfo = $state<Animal | undefined>(undefined);
  let isLoading = $state(true);
  let cardDone = $state(false);

  async function newCard() {
    cardInfo = await getRandomCard();
    isLoading = false;
    isLiked = undefined;
    cardDone = false;
  }

  onMount(async () => {
    await newCard();
  });

  async function likedCard() {
    isLiked = true;
    await cardResponse(isLiked, "001");
    setTimeout(async () => {
      cardDone = true;
      await newCard();
    }, 1300);
  }

  async function dislikedCard() {
    isLiked = false;
    await cardResponse(isLiked, "001");
    setTimeout(async () => {
      cardDone = true;
      await newCard();
    }, 1300);
  }
</script>

{#if isLoading}
  <div>Loading...</div>
{:else}
  <div class="content" use:trapFocus>
    <CardComponent card_info={cardInfo} {isLiked} {cardDone} />
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
