 
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
Endpoint to retrieve count of actions by user
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
Endpoint to get the “Referral Index” of all the users
### Approach used:
Since the main challenge I encountered was the need to iterate through a large number of records, I decided to use the Depth-First Search (DFS) algorithm to recursively find all user referrals. This algorithm helps avoid visiting the same user twice, optimizing both memory usage and processing time. Additionally, the algorithm accounts for scenarios where a user might be referred by multiple users, even though such cases are not expected to occur.
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
There is limited test coverage in this implementation. I prioritized writing unit tests for the most critical parts of the code, while intentionally omitting integration and component tests for the purposes of this challenge.
```shell
$ make test
```

#### Improvements to be done:

[Imrpovements documentation](documentation/improvements.md)
* Improving testing.
* Integration of read models and events.
* Validation on the request schema.
* Improving logging and monitoring.
