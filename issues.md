# Issues

## General Issues

- [ ] when on sign in or login page it shows the profile picture layout in the top right, figure out how to remove this on these pages.
- [x] sending update request from frontend to backend does not properly send information, works with postman though
- [ ] when creating a USER account, it is just empty and then still logs us in, need to find out why it 
            doesnt return with an error and why it is giving us an error.

## Frontend Issues

- [ ] when creating an account on both clients, it should auto sign in the newly created user.
- [ ] where to call get tokens, so that it updates whenever i reload/switch the page? 

## Backend Issues

- [ ] when trying to login gtes error (probably issue with how im sending the user to check if it exists in usercontroller)
      -Error getting user from database500{"id":"","status":200,"username":""}
- [x] Getting shots from backend in searchup.go/SelectShots gives an error in the rows.Scan()
