# schema

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

## Table

```
              room   user
-----------   ----   ----
GET  /fetch   CINE   CINE*
POST /pub     CINE   COU
POST /lock    must   must*
POST /unlock  must   must*
POST /users   must   COU

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

Response:

```
{
  "message": {
    "color": "#ff0000",
    "message": "text",
    "name": "nik",
    "ts": 1775064288968
  },
  "status": {
    "locked": false,
    "users": [
      "one",
      "two"
    ]
  }
}
```

Response can have users list

### GET /lock?room=main&user=mng3flk5-zt8nqndakv

### GET /unlock?room=main&user=mng3flk5-zt8nqndakv

### GET /list?room=main&user=mng3flk5-zt8nqndakv

