# Services

| Service | Description |
| ------- | ----------- |
| CaseService | CaseService is the API to handle cases |
| EventService | EventService is the API to handle events |
| FileService | FileService is the API for handling files |
| LinkService | LinkService is a API for creating links between objects |
| ProcessService | ProcessService is the API - that handles evidence-processing |
| TestService | TestService is used for testing-purposes |

## CaseService

### Methods

| Method | Endpoint | Description | Request | Response |
| ------ | -------- | ----------- | ------- | -------- |
| Delete | /CaseService.Delete | Delete deletes the specified case | CaseDeleteRequest | CaseDeleteResponse |
| Get | /CaseService.Get | Get returns the requested case | CaseGetRequest | CaseGetResponse |
| List | /CaseService.List | List the cases for a specified user | CaseListRequest | CaseListResponse |
| New | /CaseService.New | New creates a new case | CaseNewRequest | CaseNewResponse |
| Update | /CaseService.Update | Update updates the specified case | CaseUpdateRequest | CaseUpdateResponse |

#### Delete

Delete deletes the specified case

##### Endpoint

POST `/CaseService.Delete`

##### Request

_CaseDeleteRequest is the input-object
for deleting an existing case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the case to delete | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/CaseService.Delete
```

```json
{
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_CaseDeleteResponse is the output-object
for deleting an existing case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Get

Get returns the requested case

##### Endpoint

POST `/CaseService.Get`

##### Request

_CaseGetRequest is the input-object
for getting a specified case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the case to get | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/CaseService.Get
```

```json
{
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_CaseGetResponse is the output-object
for getting a specified case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| case | Case |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "case": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "creatorID": "7a1713b0249d477d92f5e10124a59861",
        "description": "This is a case",
        "files": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                },
                "description": "This file contains evidence",
                "mime": "@file/plain",
                "name": "text-file.txt",
                "path": "/filestore/text-file.txt",
                "processed": false,
                "size": 450060
            }
        ],
        "fromDate": 1100127600,
        "investigators": [
            "sja@avian.dk",
            "jis@avian.dk"
        ],
        "name": "Case 1",
        "processes": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                }
            }
        ],
        "toDate": 1257894000
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### List

List the cases for a specified user

##### Endpoint

POST `/CaseService.List`

##### Request

_CaseListRequest is the input-object for
listing cases for a specified user_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| userID | string | UserID of the user to list cases for | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"userID":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/CaseService.List
```

```json
{
    "userID": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_CaseListResponse is the output-object for
listing cases for a specified user_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| cases | []Case |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "cases": [
        {
            "base": {
                "createdAt": 1257894000,
                "deletedAt": 0,
                "id": "7a1713b0249d477d92f5e10124a59861",
                "updatedAt": 0
            },
            "creatorID": "7a1713b0249d477d92f5e10124a59861",
            "description": "This is a case",
            "files": [
                {
                    "base": {
                        "createdAt": 1257894000,
                        "deletedAt": 0,
                        "id": "7a1713b0249d477d92f5e10124a59861",
                        "updatedAt": 0
                    },
                    "description": "This file contains evidence",
                    "mime": "@file/plain",
                    "name": "text-file.txt",
                    "path": "/filestore/text-file.txt",
                    "processed": false,
                    "size": 450060
                }
            ],
            "fromDate": 1100127600,
            "investigators": [
                "sja@avian.dk",
                "jis@avian.dk"
            ],
            "name": "Case 1",
            "processes": [
                {
                    "base": {
                        "createdAt": 1257894000,
                        "deletedAt": 0,
                        "id": "7a1713b0249d477d92f5e10124a59861",
                        "updatedAt": 0
                    }
                }
            ],
            "toDate": 1257894000
        }
    ]
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### New

New creates a new case

##### Endpoint

POST `/CaseService.New`

##### Request

_CaseNewRequest is the input-object
for creating a new case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| name | string | Name of the case | Case 1 |
| description | string | description of the case to create | This is a case |
| fromDate | int64 | FromDate is the unix-date for the start of the primary timespan for the case | 1.1001276e+09 |
| toDate | int64 | ToDate is the unix-date for the end of the primary timespan for the case | 1.257894e+09 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"description":"This is a case","fromDate":1100127600,"name":"Case 1","toDate":1257894000}' http://localhost:8080/api/CaseService.New
```

```json
{
    "description": "This is a case",
    "fromDate": 1100127600,
    "name": "Case 1",
    "toDate": 1257894000
}
```

##### Response

_CaseNewResponse is the output-object
for creating a new case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| new | Case |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "new": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "creatorID": "7a1713b0249d477d92f5e10124a59861",
        "description": "This is a case",
        "files": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                },
                "description": "This file contains evidence",
                "mime": "@file/plain",
                "name": "text-file.txt",
                "path": "/filestore/text-file.txt",
                "processed": false,
                "size": 450060
            }
        ],
        "fromDate": 1100127600,
        "investigators": [
            "sja@avian.dk",
            "jis@avian.dk"
        ],
        "name": "Case 1",
        "processes": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                }
            }
        ],
        "toDate": 1257894000
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Update

Update updates the specified case

##### Endpoint

POST `/CaseService.Update`

##### Request

_CaseUpdateRequest is the input-object
for updating an existing case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the case to update | 7a1713b0249d477d92f5e10124a59861 |
| name | string | Name of the case | Case 1 |
| description | string | description of the case to create | This is a case |
| fromDate | int64 | FromDate is the unix-date for the start of the primary timespan for the case | 1.1001276e+09 |
| toDate | int64 | ToDate is the unix-date for the end of the primary timespan for the case | 1.257894e+09 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"description":"This is a case","fromDate":1100127600,"id":"7a1713b0249d477d92f5e10124a59861","name":"Case 1","toDate":1257894000}' http://localhost:8080/api/CaseService.Update
```

```json
{
    "description": "This is a case",
    "fromDate": 1100127600,
    "id": "7a1713b0249d477d92f5e10124a59861",
    "name": "Case 1",
    "toDate": 1257894000
}
```

##### Response

_CaseUpdateResponse is the output-object
for updating an existing case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| updated | Case |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "updated": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "creatorID": "7a1713b0249d477d92f5e10124a59861",
        "description": "This is a case",
        "files": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                },
                "description": "This file contains evidence",
                "mime": "@file/plain",
                "name": "text-file.txt",
                "path": "/filestore/text-file.txt",
                "processed": false,
                "size": 450060
            }
        ],
        "fromDate": 1100127600,
        "investigators": [
            "sja@avian.dk",
            "jis@avian.dk"
        ],
        "name": "Case 1",
        "processes": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                }
            }
        ],
        "toDate": 1257894000
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

## EventService

### Methods

| Method | Endpoint | Description | Request | Response |
| ------ | -------- | ----------- | ------- | -------- |
| Create | /EventService.Create | Create creates a new event | EventCreateRequest | EventCreateResponse |
| Delete | /EventService.Delete | Delete deletes an existing event | EventDeleteRequest | EventDeleteResponse |
| Get | /EventService.Get | Get the specified event | EventGetRequest | EventGetResponse |
| List | /EventService.List | List all events | EventListRequest | EventListResponse |
| Update | /EventService.Update | Update updates an existing event | EventUpdateRequest | EventUpdateResponse |

#### Create

Create creates a new event

##### Endpoint

POST `/EventService.Create`

##### Request

_EventCreateRequest is the input-object
for creating an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID of the case to create the event for | 7a1713b0249d477d92f5e10124a59861 |
| importance | int | Set the importance of the event, defined by a number between 1 - 5. | 3 |
| description | string | Desription of the event. | This needs investigation. |
| fromDate | int64 | FromDate is the unix-timestamp of when the event started | 1.1001276e+09 |
| toDate | int64 | ToDate is the unix-timestamp of when the event finished | 1.257894e+09 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","description":"This needs investigation.","fromDate":1100127600,"importance":3,"toDate":1257894000}' http://localhost:8080/api/EventService.Create
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "description": "This needs investigation.",
    "fromDate": 1100127600,
    "importance": 3,
    "toDate": 1257894000
}
```

##### Response

_EventCreateResponse is the output-object
for creating an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| created | Event |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "created": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "description": "This needs investigation.",
        "fromDate": 1100127600,
        "importance": 3,
        "toDate": 1257894000
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Delete

Delete deletes an existing event

##### Endpoint

POST `/EventService.Delete`

##### Request

_EventDeleteRequest is the input-object
for deleting an existing event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the event to Delete | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID of the event | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/EventService.Delete
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_EventDeleteResponse is the output-object
for deleting an existing event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Get

Get the specified event

##### Endpoint

POST `/EventService.Get`

##### Request

_EventGetRequest is the input-object
for getting an existing event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the event to get | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID of the event | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/EventService.Get
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_EventGetResponse is the output-object
for deleting an existing event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| event | Event |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "event": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "description": "This needs investigation.",
        "fromDate": 1100127600,
        "importance": 3,
        "toDate": 1257894000
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### List

List all events

##### Endpoint

POST `/EventService.List`

##### Request

_EventListRequest is the input-object
for listing all existing events for a case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID to list the events for | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/EventService.List
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_EventListResponse is the output-object
for listing all existing events for a case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| events | []Event |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "events": [
        {
            "base": {
                "createdAt": 1257894000,
                "deletedAt": 0,
                "id": "7a1713b0249d477d92f5e10124a59861",
                "updatedAt": 0
            },
            "description": "This needs investigation.",
            "fromDate": 1100127600,
            "importance": 3,
            "toDate": 1257894000
        }
    ]
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Update

Update updates an existing event

##### Endpoint

POST `/EventService.Update`

##### Request

_EventUpdateRequest is the input-object
for updating an existing event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the event to update | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID of the event | 7a1713b0249d477d92f5e10124a59861 |
| importance | int | Set the importance of the event, defined by a number between 1 - 5. | 3 |
| description | string | Desription of the event. | This needs investigation. |
| fromDate | int64 | FromDate is the unix-timestamp of when the event started | 1.1001276e+09 |
| toDate | int64 | ToDate is the unix-timestamp of when the event finished | 1.257894e+09 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","description":"This needs investigation.","fromDate":1100127600,"id":"7a1713b0249d477d92f5e10124a59861","importance":3,"toDate":1257894000}' http://localhost:8080/api/EventService.Update
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "description": "This needs investigation.",
    "fromDate": 1100127600,
    "id": "7a1713b0249d477d92f5e10124a59861",
    "importance": 3,
    "toDate": 1257894000
}
```

##### Response

_EventUpdateResponse is the output-object
for updating an existing event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| updated | Event |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "updated": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "description": "This needs investigation.",
        "fromDate": 1100127600,
        "importance": 3,
        "toDate": 1257894000
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

## FileService

### Methods

| Method | Endpoint | Description | Request | Response |
| ------ | -------- | ----------- | ------- | -------- |
| Delete | /FileService.Delete | Delete deletes the specified file | FileDeleteRequest | FileDeleteResponse |
| New | /FileService.New | New uploads a file to the backend | FileNewRequest | FileNewResponse |
| Update | /FileService.Update | Update updates the information for a file | FileUpdateRequest | FileUpdateResponse |

#### Delete

Delete deletes the specified file

##### Endpoint

POST `/FileService.Delete`

##### Request

_FileDeleteRequest is the input-object
for deleting a file_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the file to delete | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID of the case where the file to delete belongs | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/FileService.Delete
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_FileDeleteResponse is the output-object
for deleting a file_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### New

New uploads a file to the backend

##### Endpoint

POST `/FileService.New`

##### Request

_FileNewRequest is the input-object
for creating a new file_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID of the case to upload the file | 7a1713b0249d477d92f5e10124a59861 |
| name | string | Name of the file | text-file.txt |
| description | string | Description of the file | This file contains evidence |
| mime | string | Mime is the mime-type of the file | @file/plain |
| data | string | Data of the file (base64 encoded) | iVBORw0KGgoAAAANSUhEUgAAA1IAAAEeCA....... |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","data":"iVBORw0KGgoAAAANSUhEUgAAA1IAAAEeCA.......","description":"This file contains evidence","mime":"@file/plain","name":"text-file.txt"}' http://localhost:8080/api/FileService.New
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "data": "iVBORw0KGgoAAAANSUhEUgAAA1IAAAEeCA.......",
    "description": "This file contains evidence",
    "mime": "@file/plain",
    "name": "text-file.txt"
}
```

##### Response

_FileNewResponse is the output-object
for creating a new file_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| new | File |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "new": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "description": "This file contains evidence",
        "mime": "@file/plain",
        "name": "text-file.txt",
        "path": "/filestore/text-file.txt",
        "processed": false,
        "size": 450060
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Update

Update updates the information for a file

##### Endpoint

POST `/FileService.Update`

##### Request

_FileUpdateRequest is the input-object
for updating a files information_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the file to update | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID of the case where the file to update belongs | 7a1713b0249d477d92f5e10124a59861 |
| description | string | Description of the file | This file contains evidence |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","description":"This file contains evidence","id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/FileService.Update
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "description": "This file contains evidence",
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_FileUpdateResponse is the output-object
for updating a files information_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| updated | File |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "updated": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "description": "This file contains evidence",
        "mime": "@file/plain",
        "name": "text-file.txt",
        "path": "/filestore/text-file.txt",
        "processed": false,
        "size": 450060
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

## LinkService

### Methods

| Method | Endpoint | Description | Request | Response |
| ------ | -------- | ----------- | ------- | -------- |
| CreateEvent | /LinkService.CreateEvent | CreateEvent creates a link for an event with multiple objects | LinkEventCreateRequest | LinkEventCreateResponse |
| DeleteEvent | /LinkService.DeleteEvent | DeleteEvent deletes all links to the specified event | LinkEventDeleteRequest | LinkEventDeleteResponse |
| GetEvent | /LinkService.GetEvent | GetEvent gets an event with its links | LinkEventCreateRequest | LinkEventCreateResponse |
| UpdateEvent | /LinkService.UpdateEvent | UpdateEvent updates links for the specified event | LinkEventUpdateRequest | LinkEventUpdateResponse |

#### CreateEvent

CreateEvent creates a link for an event
with multiple objects

##### Endpoint

POST `/LinkService.CreateEvent`

##### Request

_LinkEventCreateRequest is the input-object
for linking objects with an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID for the event | 7a1713b0249d477d92f5e10124a59861 |
| fromID | string | FromID is the ID of the event to hold the link | 7a1713b0249d477d92f5e10124a59861 |
| eventIDs | []string | EventIDs of the events to be linked | 7a1713b0249d477d92f5e10124a59861 |
| bidirectional | bool | Bidirectional means that he link also should be created for the "ToID" | true |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"bidirectional":true,"caseID":"7a1713b0249d477d92f5e10124a59861","eventIDs":["7a1713b0249d477d92f5e10124a59861"],"fromID":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/LinkService.CreateEvent
```

```json
{
    "bidirectional": true,
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "eventIDs": [
        "7a1713b0249d477d92f5e10124a59861"
    ],
    "fromID": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_LinkEventCreateResponse is the output-object
for linking objects with an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| linked | LinkEvent |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "linked": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "events": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                },
                "description": "This needs investigation.",
                "fromDate": 1100127600,
                "importance": 3,
                "toDate": 1257894000
            }
        ],
        "from": {
            "base": {
                "createdAt": 1257894000,
                "deletedAt": 0,
                "id": "7a1713b0249d477d92f5e10124a59861",
                "updatedAt": 0
            },
            "description": "This needs investigation.",
            "fromDate": 1100127600,
            "importance": 3,
            "toDate": 1257894000
        }
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### DeleteEvent

DeleteEvent deletes all links to the specified event

##### Endpoint

POST `/LinkService.DeleteEvent`

##### Request

_LinkEventDeleteRequest is the input-object
for removing a linked event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID of the case where the event belongs | 7a1713b0249d477d92f5e10124a59861 |
| eventID | string | EventID of the Event to get all links for | 7a1713b0249d477d92f5e10124a59861 |
| bidirectional | bool | Bidirectional - if bidirectional links also should get deleted | false |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"bidirectional":false,"caseID":"7a1713b0249d477d92f5e10124a59861","eventID":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/LinkService.DeleteEvent
```

```json
{
    "bidirectional": false,
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "eventID": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_LinkEventDeleteResponse is the output-object
for removing a linked event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### GetEvent

GetEvent gets an event with its links

##### Endpoint

POST `/LinkService.GetEvent`

##### Request

_LinkEventCreateRequest is the input-object
for linking objects with an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID for the event | 7a1713b0249d477d92f5e10124a59861 |
| fromID | string | FromID is the ID of the event to hold the link | 7a1713b0249d477d92f5e10124a59861 |
| eventIDs | []string | EventIDs of the events to be linked | 7a1713b0249d477d92f5e10124a59861 |
| bidirectional | bool | Bidirectional means that he link also should be created for the "ToID" | true |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"bidirectional":true,"caseID":"7a1713b0249d477d92f5e10124a59861","eventIDs":["7a1713b0249d477d92f5e10124a59861"],"fromID":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/LinkService.GetEvent
```

```json
{
    "bidirectional": true,
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "eventIDs": [
        "7a1713b0249d477d92f5e10124a59861"
    ],
    "fromID": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_LinkEventCreateResponse is the output-object
for linking objects with an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| linked | LinkEvent |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "linked": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "events": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                },
                "description": "This needs investigation.",
                "fromDate": 1100127600,
                "importance": 3,
                "toDate": 1257894000
            }
        ],
        "from": {
            "base": {
                "createdAt": 1257894000,
                "deletedAt": 0,
                "id": "7a1713b0249d477d92f5e10124a59861",
                "updatedAt": 0
            },
            "description": "This needs investigation.",
            "fromDate": 1100127600,
            "importance": 3,
            "toDate": 1257894000
        }
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### UpdateEvent

UpdateEvent updates links for the specified event

##### Endpoint

POST `/LinkService.UpdateEvent`

##### Request

_LinkEventUpdateRequest is the input-object
for updating linked objects with an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the linked event | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID for the event | 7a1713b0249d477d92f5e10124a59861 |
| fromID | string | FromID is the ID of the event to hold the link | 7a1713b0249d477d92f5e10124a59861 |
| eventAddIDs | []string | EventAddIDs of the events to be linked | 7a1713b0249d477d92f5e10124a59861 |
| eventRemoveIDs | []string | EventRemoveIDs of the events to be removed | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","eventAddIDs":["7a1713b0249d477d92f5e10124a59861"],"eventRemoveIDs":["7a1713b0249d477d92f5e10124a59861"],"fromID":"7a1713b0249d477d92f5e10124a59861","id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/LinkService.UpdateEvent
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "eventAddIDs": [
        "7a1713b0249d477d92f5e10124a59861"
    ],
    "eventRemoveIDs": [
        "7a1713b0249d477d92f5e10124a59861"
    ],
    "fromID": "7a1713b0249d477d92f5e10124a59861",
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_LinkEventUpdateResponse is the output-object
for linking objects with an event_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| updated | LinkEvent |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "updated": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        },
        "events": [
            {
                "base": {
                    "createdAt": 1257894000,
                    "deletedAt": 0,
                    "id": "7a1713b0249d477d92f5e10124a59861",
                    "updatedAt": 0
                },
                "description": "This needs investigation.",
                "fromDate": 1100127600,
                "importance": 3,
                "toDate": 1257894000
            }
        ],
        "from": {
            "base": {
                "createdAt": 1257894000,
                "deletedAt": 0,
                "id": "7a1713b0249d477d92f5e10124a59861",
                "updatedAt": 0
            },
            "description": "This needs investigation.",
            "fromDate": 1100127600,
            "importance": 3,
            "toDate": 1257894000
        }
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

## ProcessService

### Methods

| Method | Endpoint | Description | Request | Response |
| ------ | -------- | ----------- | ------- | -------- |
| Abort | /ProcessService.Abort | Abort aborts the specified processing-job | ProcessAbortRequest | ProcessAbortResponse |
| Jobs | /ProcessService.Jobs | Jobs returns the status of all processing-jobs in the specified case | ProcessJobsRequest | ProcessJobsResponse |
| Pause | /ProcessService.Pause | Pause pauses the specified processing-job | ProcessPauseRequest | ProcessPauseResponse |
| Start | /ProcessService.Start | Start starts a processing with the specified files | ProcessStartRequest | ProcessStartResponse |

#### Abort

Abort aborts the specified processing-job

##### Endpoint

POST `/ProcessService.Abort`

##### Request

_ProcessAbortRequest is the input-object
for aborting a processing-job_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the processing-job to abort | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID of the case the processing-job belongs to | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/ProcessService.Abort
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_ProcessAbortResponse is the output-object
for aborting a processing-job_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| aborted | Process |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "aborted": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        }
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Jobs

Jobs returns the status of all processing-jobs
in the specified case

##### Endpoint

POST `/ProcessService.Jobs`

##### Request

_ProcessJobsRequest is the input-object
for getting all processing-jobs for a case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID of the case to get the processing-jobs for | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/ProcessService.Jobs
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_ProcessJobsResponse is the output-object
for getting all processing-jobs for a case_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| processes | []Process |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "processes": [
        {
            "base": {
                "createdAt": 1257894000,
                "deletedAt": 0,
                "id": "7a1713b0249d477d92f5e10124a59861",
                "updatedAt": 0
            }
        }
    ]
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Pause

Pause pauses the specified processing-job

##### Endpoint

POST `/ProcessService.Pause`

##### Request

_ProcessPauseRequest is the input-object
for pausing a processing-job_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the processing-job to pause | 7a1713b0249d477d92f5e10124a59861 |
| caseID | string | CaseID of the case the processing-job belongs to | 7a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","id":"7a1713b0249d477d92f5e10124a59861"}' http://localhost:8080/api/ProcessService.Pause
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "id": "7a1713b0249d477d92f5e10124a59861"
}
```

##### Response

_ProcessPauseResponse is the output-object
for pausing a processing-job_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| paused | Process |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "paused": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        }
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### Start

Start starts a processing with the specified files

##### Endpoint

POST `/ProcessService.Start`

##### Request

_ProcessStartRequest is the input-object
for starting a processing-job_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| caseID | string | CaseID of the case to start the processing for | 7a1713b0249d477d92f5e10124a59861 |
| fileIDs | []string | FileIDs of the files to process | 7a1713b0249d477d92f5e10124a598617a1713b0249d477d92f5e10124a59861 |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"caseID":"7a1713b0249d477d92f5e10124a59861","fileIDs":["7a1713b0249d477d92f5e10124a59861","7a1713b0249d477d92f5e10124a59861"]}' http://localhost:8080/api/ProcessService.Start
```

```json
{
    "caseID": "7a1713b0249d477d92f5e10124a59861",
    "fileIDs": [
        "7a1713b0249d477d92f5e10124a59861",
        "7a1713b0249d477d92f5e10124a59861"
    ]
}
```

##### Response

_ProcessStartResponse is the output-object
for starting a processing-job_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| started | Process |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "started": {
        "base": {
            "createdAt": 1257894000,
            "deletedAt": 0,
            "id": "7a1713b0249d477d92f5e10124a59861",
            "updatedAt": 0
        }
    }
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

## TestService

### Methods

| Method | Endpoint | Description | Request | Response |
| ------ | -------- | ----------- | ------- | -------- |
| CreateUser | /TestService.CreateUser | CreateUser creates a test-user in Firebase | TestCreateUserRequest | TestCreateUserResponse |
| DeleteUser | /TestService.DeleteUser | DeleteUser deletes a test-user in Firebase | TestDeleteUserRequest | TestDeleteUserResponse |

#### CreateUser

CreateUser creates a test-user in Firebase

##### Endpoint

POST `/TestService.CreateUser`

##### Request

_TestCreateUserRequest is the input-object
for creating a test-user_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| name | string | Name of the user to create | Simon |
| id | string | ID of the user to create | aaef42k4t2 |
| email | string | Email of the user to create | sja@avian.dk |
| password | string | Password for the new user | supersecret |
| secret | string | Secret for using the test-service | supersecret |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"email":"sja@avian.dk","id":"aaef42k4t2","name":"Simon","password":"supersecret","secret":"supersecret"}' http://localhost:8080/api/TestService.CreateUser
```

```json
{
    "email": "sja@avian.dk",
    "id": "aaef42k4t2",
    "name": "Simon",
    "password": "supersecret",
    "secret": "supersecret"
}
```

##### Response

_TestCreateUserResponse is the output-object
for creating a test-user_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| token | string | Token for the created user | er324235tt.... |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "token": "er324235tt...."
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```

#### DeleteUser

DeleteUser deletes a test-user in Firebase

##### Endpoint

POST `/TestService.DeleteUser`

##### Request

_TestDeleteUserRequest is the input-object
for deleting a test-user_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| id | string | ID of the user to delete | aaef42k4t2 |
| secret | string | Secret for using the test-service | supersecret |

```sh
curl -H "Content-Type: application/json" -X POST -d '{"id":"aaef42k4t2","secret":"supersecret"}' http://localhost:8080/api/TestService.DeleteUser
```

```json
{
    "id": "aaef42k4t2",
    "secret": "supersecret"
}
```

##### Response

_TestDeleteUserResponse is the output-object
for deleting a test-user_

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```
