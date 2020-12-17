# Services

| Service | Description |
| ------- | ----------- |
| CaseService | CaseService is the API to handle cases |
| FileService | FileService is the API for handling files |
| ProcessService | ProcessService is the API - that handles evidence-processing |

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
| file | File |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "file": {
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
| file | File |  |  |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "file": {
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
