 
## Quick run

```shell
$ cp .env.dist .env
$ make install
$ make run
```

### Get user info
Endpoint to retrieve the user info by user id
```
curl --location 'http://localhost:8080/api/users/1' \
--header 'Content-Type: application/vnd.surfe.v1+json'
```

#### Response
```
{
    "id": 1,
    "name": "Ferdinande",
    "createdAt": "2020-07-14T05:48:54Z"
}
```
### Get action count by user
Endpoint to retrieve count of actions vy user
```
curl --location 'http://localhost:8080/api/actions/users/4' \
--header 'Content-Type: application/vnd.surfe.v1+json'
```

#### Response
```
{
    "count": 34
}
```

### Get probability 
Endpoint to retrieve probability of next action after by action name
```
curl --location 'http://localhost:8080/api/actions/probability/users/EDIT_CONTACT' \
--header 'Content-Type: application/vnd.surfe.v1+json'
```

#### Response
```
{
    "ADD_CONTACT": 0.33,
    "EDIT_CONTACT": 0.33,
    "REFER_USER": 0.02,
    "VIEW_CONTACTS": 0.32
}
```

### Get referrals by user
endpoint to get the “Referral Index” of all the Users.
### Approach used:
Since the main issue that i found was that i had to iterate a lot of records, i decided to use the Depth-First Search algorithm in order to find all the user referrals recursively, this algorithm 
allows you to avoid visiting the same user twice and optimize the memory and time of processing.
The algorithm is also taking into account if a user is referred by 2 users (even if it cannot happen).
```
curl --location 'http://localhost:8080/api/actions/referral' \
--header 'Content-Type: application/vnd.surfe.v1+json'
```

#### Response
```
{
    "1": 1,
    "10": 1,
    "104": 3,
    "106": 0,
    "11": 1,
    "110": 4,
    "112": 2,
    "113": 1,
    "114": 0,
    "115": 0,
    "116": 4,
    "117": 0,
...
```

## TEST
There's not a lot of coverage of testing, i decided to test the most important with unit test, and i ignored the integration and component test for this challenge.  
```shell
$ make test
```

#### Improvements to be done:

[Imrpovements](documentation/improvements.md)
* Improving testing.
* Integration of read models and events.
* Validation on the request schema.
* Improving logging and monitoring.