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



- User swipes yes on animal -> organization can accept req -> if yes makes a new message with both user and org
- probably need a new db table called notifications that has
      - id
      - actor ref users(id)
      - notifier ref users(id)
      - entity text | what is the notif about
      - entity_type (message-request | message-reply | news | alert)
      - status (accepted | denied | unseen)
      - date_created 
      - date_seen
