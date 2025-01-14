# TODO

- MAKE SURE TO CALL GETTOKENS IN FUNCTIONS WHERE TOKENS ARE NEEDED (AWAIT)

## Format
- DESCRIPTION | PRIORITY LEVEL(1 being highest, 5 being lowest)

## General

## Frontend
- if backend is down, can i not go to login/signup page? | 4
-  when calling GetTokens inside funciton calls(getcurrentanimals) it doesn't seem to properly await the function? it says invalid token even though its an await? | 2

## Backend

- Compress images given from frontend | 3
- When an animal is deleted, delete its image entry from AWS - current | 1
      right now objects are getting stored as a232avasdjpeg instead of a232avasd.jpeg, remove the jpeg since it auto does that in s3 with the mimetype OR
      do something like a232avasd_jpeg so that when i want to delete the item i can strings.split(url, "_") and just get the last item? 
      tried by doing this key + ".jpeg" but that still didn't work so maybe its an issue with the query as well instead of just how the key is stored?