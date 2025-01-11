## General

- [ ] CHANGE ALL CONTEXT FOR SQL TO R.CONTEXT AND UPDATE POOL TO JUST HAVE 1?
      - example func GetUser(context.Context, id any) -> GetUser(r.Header, id)
            I like to pass the request context r.Context() into all functions in handlers, propagating it. Any queries that are running will be cancelled if the context is cancelled. 
            If you have some background process in a go routine (within a handler), you do not want to use the request context, because it'll be cancelled when the handler exits/finishes 

- [ ] Work with cloudinary to upload images from frontend to backend and set the url in the animal_images
      Instead of updating animals from 1 call make 2 calls, one to create the animal and one to create animal images
      the one for images should send the animal id and the image.

## Frontend
- [ ] where to call GetTokens so that each time the page is rerendered(or switch from / to /profile) it refreshes token? or something else
- [x] Call /cards to get a card on page load, and set the cards image instead of hard coding it

## Backend

- [x] set up a route that returns a random animal? (route could be GET /cards instead of POST)
- [x] users route should check for post and then  
- [ ] organizations should have a pfp aswell/ main picture
- [ ] LOGIN GIVES INVALID MEMORY ADDRESS
- [ ] Compress images given from frontend