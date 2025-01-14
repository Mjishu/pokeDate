# TODO

- MAKE SURE TO CALL GETTOKENS IN FUNCTIONS WHERE TOKENS ARE NEEDED (AWAIT)

## Format
- DESCRIPTION | PRIORITY LEVEL(1 being highest, 5 being lowest)

## General

## Frontend
- if user is NOT signed in should redirect them to login page 
- if user IS logged in, they should not be able to go to the sign in or login pages.
-  when calling GetTokens inside funciton calls(getcurrentanimals) it doesn't seem to properly await the function? it says invalid token even though its an await?

## Backend

- Compress images given from frontend | 3
- When an animal is deleted, delete its image entry from AWS - current | 1