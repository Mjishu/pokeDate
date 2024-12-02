<script lang="ts">
  import { trapFocus } from "./helpers/actions.svelte";

  let isLiked = $state();

  export function likedCard() {
    isLiked = true;
    console.log("You liked this card!");
  }

  export function dislikedCard() {
    isLiked = false;
    console.log("Youd disliked this card");
  }
</script>

<div class="content" use:trapFocus>
  <button onclick={() => dislikedCard()}>&lt</button>
  <img
    src="./images/dog.webp"
    alt=""
    class="card"
    class:liked={isLiked}
    class:disliked={isLiked == false}
  />
  <button onclick={() => likedCard()}>&gt</button>
  <button onclick={() => console.log($state.snapshot(isLiked))}
    >State of isLiked</button
  >
</div>

<style>
  .content {
    display: flex;
  }
  .card {
    margin: 0;
    padding: 0;
    border-radius: 5px;
    border: none;
    transition:
      width 2s,
      height 2s;
  }

  .card.liked {
    transform: translateX(1);
    border: 1px solid green;
    width: 50rem;
  }

  .card.disliked {
    transform: translateX(-1);
    border: 1px solid red;
  }
</style>
