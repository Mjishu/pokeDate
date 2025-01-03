# Issues

## General Issues

- [x] sending update request from frontend to backend does not properly send information, works with postman though

## Frontend Issues

- [ ] when using +layout.ts and +page.svelte, i cant pass data down it just says undefined?
- [x] moving over to sveltekit for organizations the image route isnt working for trash can and searchi con, fix it, same with global styles
- [x] add part to only use sveltekit for frontend and turn off ssr

## Backend Issues

- [ ] when trying to login gtes error (probably issue with how im sending the user to check if it exists in usercontroller)
      -Error getting user from database500{"id":"","status":200,"username":""}
- [x] Getting shots from backend in searchup.go/SelectShots gives an error in the rows.Scan()
