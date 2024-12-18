<script lang="ts">
      import { format } from "date-fns";

      let props = $props();
      console.log($state.snapshot(props.card_info));
</script>

<div
      class="card"
      class:liked={props.isLiked == true}
      class:disliked={props.isLiked == false}
      class:cardDone={props.cardDone}
>
      <img src={props.card_info.Image_src} alt="animal" class="image-1" />
      <div class="animal-info">
            <h2 class="name">{props.card_info.Name}</h2>
            <p class="species">{props.card_info.Species}</p>
            <p class="breed">{props.card_info.Breed}</p>
            <p class="date-of-birth">
                  {format(props.card_info.Date_of_birth, "yyyy/MM/dd")}
            </p>
            <p class="price">${props.card_info.Price}</p>
      </div>
      <div class="shots">
            <h3>Shots</h3>
            {#each props.card_info.shots as shot}
                  <div class="shot-holder">
                        <h5>{shot.name}</h5>
                        <p>{shot.description}</p>
                        <p>{shot.date_given}</p>
                        <p>{shot.next_due}</p>
                  </div>
            {/each}
      </div>
      <p class="sex">{props.card_info.Sex}</p>
      <button class="skip-button"
            ><img src="/icons/cancel_icon.svg" alt="skip" /></button
      >
</div>

<style>
      .card {
            border: 1px solid #312b2c;
            border-radius: 5px;
            /* border: none; */
            margin: 0;
            padding: 0;
            width: 30rem;
            height: 35rem;
            transform: translateX(0);
            opacity: 1;
            transition:
                  transform 1.25s ease-in-out,
                  opacity 1s ease-out 500ms;
            overflow-y: scroll;

            position: relative;
      }
      .skip-button {
            position: absolute;
            right: 0;
            top: 0;
            border: none;
            background-color: transparent;
      }
      .cardDone {
            transition: none !important;
            transform: translate3d(0, 0, 0) rotateZ(0) !important;
            opacity: 1 !important;
      }
      .skip-button img {
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

      .image-1 {
            width: 100% !important;
            height: 30rem;
            background: #d9d9d9;
      }

      .sex {
            position: absolute;
            background-color: hsla(var(--tag-color), 0.75);
            right: 1rem;
            top: 25.5rem;
            color: white;
            width: 6rem;
            height: 2rem;
            font-size: 24px;
            font-weight: 200;
            text-align: center;
            border-radius: 10px;
      }

      .animal-info {
            padding: 0;
            margin: 0;
            padding-left: 3rem;
            position: relative;
      }

      .shots {
            padding: 0;
            margin: 0;
            padding-left: 3rem;
            position: relative;
            font-size: 32px;
            font-weight: 800;
      }

      .price {
            position: absolute;
            top: -2rem;
            right: 2rem;
            font-size: 24px;
            font-weight: 100;
      }
      .name {
            font-weight: 600;
            font-size: 40px;
      }
      .species {
            font-weight: 400;
            font-size: 20px;
      }
      .breed {
            font-weight: 300;
            font-size: 20px;
      }

      .date-of-birth {
            font-size: 1.5rem;
      }
</style>
