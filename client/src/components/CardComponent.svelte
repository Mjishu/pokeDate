<script lang="ts">
      let props = $props();
</script>

<div
      class="card"
      class:liked={props.isLiked}
      class:disliked={props.isLiked == false}
>
      <img
            src={props.card_info.animal_info.image_src}
            alt="animal"
            class="card"
      />
      <h2>{props.card_info.animal_info.name}</h2>
      <p>{props.card_info.animal_info.species}</p>
      <div>
            <p>Date of birth</p>
            <p>{props.card_info.animal_info.date_of_birth}</p>
      </div>
      <p>Sex: {props.card_info.animal_info.sex}</p>
      <h5>Available: {props.card_info.animal_info.available}</h5>
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
                  width 2s ease-in-out,
                  transform 1.25s ease-in-out,
                  opacity 1s ease-out 500ms;
            overflow-y: scroll;
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
