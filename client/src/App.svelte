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
  <img
    src="./images/dog.webp"
    alt=""
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
    border: 1px solid green;
    opacity: 0.1;
  }

  .card.disliked {
    transform: translate3d(-25rem, 5rem, 0) rotateZ(-35deg);
    border: 1px solid red;
    opacity: 0.1;
  }
</style>
