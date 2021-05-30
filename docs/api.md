## **Fupisha(v1) API Document**

Fupisha API is a REST API. Therefore,this documentation in this section assumes knowledge of REST concepts.

> This API uses POST, GET, PUT, DELETE requests to communicate and HTTP response codes to indenticate status and errors. All responses come in standard JSON. All requests and responses must include a content-type of application/json and the body must be valid JSON. We try to use verbs that match both request type (fetching vs modifying) and plurality (one vs multiple).

## **Status Codes**

The API is designed to return different status codes according to context and
action. This way, if a request results in an error, the caller is able to get
insight into what went wrong.

### Request Types

| Request type  | Description                                                                                                      |
| ------------- | ---------------------------------------------------------------------------------------------------------------- |
| `GET`         | Access one or more resources and return the result as JSON.                                                      |
| `POST`        | Return `201 Created` if the resource is successfully created and return the newly created resource as JSON.      |
| `GET` / `PUT` | Return `200 OK` if the resource is accessed or modified successfully. The (modified) result is returned as JSON. |
| `DELETE`      | Returns `204 No Content` if the resource was deleted successfully.                                               |

### Response Types

| Return values            | Description                                                                                                                                                   |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `200 OK`                 | The `GET`, `PUT` or `DELETE` request was successful, the resource(s) itself is returned as JSON.                                                              |
| `204 No Content`         | The server has successfully fulfilled the request and that there is no additional content to send in the response payload body.                               |
| `201 Created`            | The `POST` request was successful and the resource is returned as JSON.                                                                                       |
| `304 Not Modified`       | Indicates that the resource has not been modified since the last request.                                                                                     |
| `400 Bad Request`        | A required attribute of the API request is missing, e.g., the title of an issue is not given.                                                                 |
| `401 Unauthorized`       | The user is not authenticated, a valid [user token](#authentication) is necessary.                                                                            |
| `403 Forbidden`          | The request is not allowed, e.g., the user is not allowed to delete a project.                                                                                |
| `404 Not Found`          | A resource could not be accessed, e.g., an ID for a resource could not be found.                                                                              |
| `405 Method Not Allowed` | The request is not supported.                                                                                                                                 |
| `409 Conflict`           | A conflicting resource already exists, e.g., creating a project with a name that already exists.                                                              |
| `412`                    | Indicates the request was denied. May happen if the `If-Unmodified-Since` header is provided when trying to delete a resource, which was modified in between. |
| `422 Unprocessable`      | The entity could not be processed.                                                                                                                            |
| `500 Server Error`       | While handling the request something went wrong server-side.                                                                                                  |

## Authentication

Fupisha API requires authentication, or will only return public data when authentication is not provided. For those cases where it is not required, this will be mentioned in the documentation for each individual endpoint.

There are several ways to authenticate with the Fupisha API:

- API Key
- JWT

## APIKey AUTH

Used by third party applications to access fupisha.
| | |
| ----- |------ |
| Security Scheme Type | API Key |
| Header parameter name| X-API-KEY |

## JWT AUTH

Used by web clients to access fupisha.

|                       |               |
| --------------------- | ------------- |
| Security Scheme Type  | Bearer        |
| Header parameter name | Authorization |

## Generate API Key

**URL** : `/api/auth/apikey`

**Method** : POST

**Auth required** : YES

**Header required** : `Api:v1`

**Permissions required** : User is Account Owner

**Data** : `{}`

### Success Response

**Condition** : If the Account exists.

**Code** : `201 CREATED`

**Content** :

```json
{
  "apiKey": "Ggg5LYu6SpaxyYs9RAc_BK"
}
```

### Error Responses

**Condition** : If there was no Bearer Token in the request Header.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": "Unauthorized",
  "error": "invalid or expired login token"
}
```

### Or

**Condition** : If there was a malformed Bearer Token in the request Header.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": "Unauthorized",
  "error": "invalid or expired login token"
}
```

### Or

**Condition** : If there was an expired Bearer Token in the request Header.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": "Unauthorized",
  "error": "invalid or expired login token"
}
```

### Or

**Condition** : If Api Version Header is missing or invalid.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
  "status":"Bad Request",
  "error":"missing api version header"
}
```

## Login

Used to collect a Token for a registered User.

**URL** : `/api/auth/login`

**Method** : `POST`

**Auth required** : NO

**Header required** : `Api:v1`

**Data constraints**

```json
{
  "email": "[valid email address]",
  "password": "[password in plain text]"
}
```

**Data example**

```json
{
  "email": "user@fupisha.io",
  "password": "abcd1234"
}
```

### Success Response

**Code** : `200 OK`

**Content example**

```json
{
  "email": "user@fupisha.io",
  "id": "udKxcNIyTiaohWkAVPH0Jg",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjIzOTY1MTUsImlhdCI6MTYyMjM5NjE1NSwiaXNzIjoiZnVwaXNoYSIsIl91aWQiOiJiOWQyYjE3MC1kMjMyLTRlMjYtYTg4NS02OTAwNTRmMWY0MjYifQ.EJ5ZIbTn1lkaAnflYTb3Fy6aAzjGO3ePYwb7tDOO-Y0"
}
```

### Error Response

**Condition** : If 'email' and 'password' combination is wrong.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": "Unauthorized",
  "error": "invalid email or password"
}
```

### Or

**Condition** : If Api Version Header is missing or invalid.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
  "status":"Bad Request",
  "error":"missing api version header"
}
```

## Signup

Used to registered a User.

**URL** : `/api/auth/signup`

**Method** : `POST`

**Auth required** : NO

**Header required** : `Api:v1`

**Data constraints**

```json
{
  "email": "[valid email address]",
  "password": "[password in plain text]"
}
```

**Data example**

```json
{
  "email": "user@fupisha.io",
  "password": "abcd123456"
}
```

### Success Response

**Code** : `200 OK`

**Content example** : `{}`

### Error Response

**Condition** : If 'email' or 'name' or 'password' combination is invalid.

**Code** : `401 UNAUTHORIZED`

**Content** :

```json
{
  "status": "Unauthorized",
  "error": "invalid email or password"
}
```

### Or

**Condition** : If Api Version Header is missing or invalid.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
  "status":"Bad Request",
  "error":"missing api version header"
}
```