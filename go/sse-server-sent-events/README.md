# IDEA 0.1

- Remove seeded status
- Universal DTOs

## Handlers actions and permissions

```
             room    user
             ----    ----
POST /enter  create  create/update
POST /pub    create  create/update
GET  /fetch  must    must
POST /lock   must    must
```

### /enter

Perform very first.

- Create/update user
- Fetch lock staus
- Fetch users list

### /pub

- Create/update user (new color, new name)

### /fetch

- User must exist
- Check it on every iteration, send death message if user not exist

### /lock

- User must exist to change lock status

# IDEA 0 (outdated)

## Users statuses

- not present
- seeded user: showed up, however hadn't sent any messagest and we do not know it's color and nik
  - appears in lists as anonymous
  - can read (according room restrictions)
  - can send (according room restrictions)
  - can get user list (according room restrictions)
  - CANNOT lock the room (according room restrictions)
- created: send at lease one message and fully sated up with nik and color
  - appears with full name and color in lists
  - can read (according room restrictions)
  - can send (according room restrictions)
  - can get user list (according room restrictions)
  - can lock the room (according room restrictions)

## Room statuses (user set statuses)

- open (not locked):
  - new users are add automatically as seeded
- locked:
  - new users can not to do anything: no reading, no writing, no getting users list

## Room expiration policy

**TODO**

## Access

```
              access
              -----------------------------------
              room    room    user    user    user
              not     exists  not     seeded  full
              exists          exists
              ----    ----    ----    ----    ----
GET  /fetch   Y(c)    Y       Y(s)    Y       Y
POST /pub     Y(c)    Y       Y(c)    Y(u)    Y
POST /lock    N       Y       N       N       Y
POST /unlock  N       Y       N       N       Y
POST /list    Y(e)    Y       Y(s)    Y       Y
----------
(c) - create
(e) - empty response
(s) - create seeded (user)
(u) - update user
```

```
just idea (legacy):

              room   user
-----------   ----   ----
GET  /fetch   CINE   CINE*
POST /pub     CINE   COU
POST /lock    must   must*
POST /unlock  must   must*
POST /list    must   COU    ?
POST /me      must   must*  ?

CINE — create if not exists
CINE* — in seeded mode if new
COU — create or update
must* — must exist AND fully setted up
```

## Methods

### POST /pub

- publish message
- legalize user: set color and nik and publish update

Request:

```json
{
  "room": "main",
  "user": "mng3flk5-zt8nqndakv",
  "color": "#ff0000",
  "name": "nik",
  "message": "text"
}
```

### GET /fetch?room=main&user=mng3flk5-zt8nqndakv

- fetch messages

Response (one SSE message):

```
{
  "message": {
    "color": "#ff0000",
    "message": "text",
    "name": "nik",
    "ts": 1775064288968
  },
  "status": { // optional
    "locked": true,
    "users": [
      {
        "name": "me",      // can be empty for just seeded user
        "color": "#ff0000" // can be empty
      }
    ]
  }
}
```

Response can have users list

### POST /lock

```
{
  "room": "main",
  "user": "mnkh2nej-h6gbfxpd7vh"
}
```

### POST /unlock

```
{
  "room": "main",
  "user": "mnkh2nej-h6gbfxpd7vh"
}
```

### POST /list (TODO +/me?)

```
{
  "room": "main",
  "user": "mnkh2nej-h6gbfxpd7vh"
}
```

### POST /me

```
{
  "room": "main",
  "user": "mnkh2nej-h6gbfxpd7vh"
}
```

