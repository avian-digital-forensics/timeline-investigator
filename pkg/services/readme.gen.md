# Services

| Service | Description |
| ------- | ----------- |
| Service | Service is the main-service |

## Service

### Methods

| Method | Endpoint | Description | Request | Response |
| ------ | -------- | ----------- | ------- | -------- |
| Greet | /Service.Greet | Greet sends a polite greeting | GreetRequest | GreetResponse |

#### Greet

Greet sends a polite greeting

##### Endpoint

POST `/Service.Greet`

##### Request

_GreetRequest is the request object for GreeterService.Greet._

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| name | string | Namee of the person to greet | Simon |

```json
{
    "name": "Simon"
}
```

##### Response

_GreetResponse is the response object containing a
person&#39;s greeting._

**Fields**

| Name | Type | Description | Example |
| ---- | ---- | ----------- | ------- |
| greeting | string | Greeting is a nice message welcoming somebody. | Hello Simon |
| error | string | Error is string explaining what went wrong. Empty if everything was fine. | something went wrong |

`200 OK`

```json
{
    "greeting": "Hello Simon"
}
```

`500 Internal Server Error`

```json
{
    "error": "something went wrong"
}
```
