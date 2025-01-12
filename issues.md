# Issues

## General Issues

- [x] sending update request from frontend to backend does not properly send information, works with postman though

## Frontend Issues

- [ ] when creating an account on both clients, it should auto sign in the newly created user.
- [ ] where to call get tokens, so that it updates whenever i reload/switch the page? 

## Backend Issues

- [ ] when trying to login gtes error (probably issue with how im sending the user to check if it exists in usercontroller)
      -Error getting user from database500{"id":"","status":200,"username":""}
- [x] Getting shots from backend in searchup.go/SelectShots gives an error in the rows.Scan()
