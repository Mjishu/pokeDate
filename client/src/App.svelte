<script lang="ts">
  import { cardResponse, getRandomCard } from "./helpers/card";
  import { trapFocus } from "./helpers/actions.svelte";
  import { onMount } from "svelte";

  type Animal = {
    id: string;
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

  onMount(async () => {
    cardInfo = await getRandomCard();
    isLoading = false;
  });

  async function likedCard() {
    isLiked = true;
    await cardResponse(isLiked, "001");
  }

  async function dislikedCard() {
    isLiked = false;
    await cardResponse(isLiked, "001");
  }
</script>

{#if isLoading}
  <div>Loading...</div>
{:else}
  <div class="content" use:trapFocus>
    <img
      src={cardInfo?.animal_info.image_src}
      alt="animal"
      class="card"
      class:liked={isLiked}
      class:disliked={isLiked == false}
    />
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

  .card {
    margin: 0;
    padding: 0;
    border-radius: 5px;
    width: 20rem;
    transform: translateX(0);
    border: none;
    opacity: 1;
    transition:
      width 2s ease-in-out,
      transform 1.25s ease-in-out,
      opacity 1s ease-out 500ms;
  }

  .card.liked {
    transform: translate3d(25rem, 5rem, 0) rotateZ(35deg);
    opacity: 0.1;
  }

  .card.disliked {
    transform: translate3d(-25rem, 5rem, 0) rotateZ(-35deg);
    opacity: 0.1;
  }
</style>
