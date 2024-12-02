<script lang="ts">
  import { likedCard, dislikedCard } from "./helpers/card";
  import { trapFocus } from "./helpers/actions.svelte";

  function onKeyDown(e) {
    switch (e.keyCode) {
      case 37:
        likedCard();
        break;
      case 39:
        dislikedCard();
        break;
    }
  }
</script>

<div class="content" use:trapFocus>
  <div>This is the images!</div>
</div>

<img src="./images/dog.webp" alt="" />

<style>
  .content {
    display: flex;
    flex-direction: column;
  }
</style>
