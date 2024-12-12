<script lang="ts">
      let props = $props();
      console.log($state.snapshot(props.card_info));
</script>

<div
      class="card"
      class:liked={props.isLiked == true}
      class:disliked={props.isLiked == false}
      class:cardDone={props.cardDone}
>
      <img src={props.card_info.Image_src} alt="animal" class="card" />
      <h2>{props.card_info.Name}</h2>
      <p>{props.card_info.Species}</p>
      <div>
            <p>Date of birth</p>
            <p>{props.card_info.Date_of_birth}</p>
      </div>
      <p>Sex: {props.card_info.Sex}</p>
      <h5>Available: {props.card_info.Available}</h5>
</div>

<style>
      .card {
            border: 1px solid #312b2c;
            border-radius: 5px;
            /* border: none; */
            margin: 0;
            padding: 0;
            width: 30rem;
            height: 25rem;
            width: 20rem;
            transform: translateX(0);
            opacity: 1;
            transition:
                  transform 1.25s ease-in-out,
                  opacity 1s ease-out 500ms;
            overflow-y: scroll;
      }
      .cardDone {
            transition: none !important;
            transform: translate3d(0, 0, 0) rotateZ(0) !important;
            opacity: 1 !important;
      }
      .card img {
            width: auto;
            height: auto;
            border: none;
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
