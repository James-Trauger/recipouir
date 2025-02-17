# Recipe API

## Dependencies

- Docker
- A trusted certificate and private key

The database (MongoDB) and the api is ran and built in a Docker container. Also, a locally trusted certificate is needed. I recommend using [mkcert](https://github.com/FiloSottile/mkcert) to aquire this.

## Setup

After successfully adding the trusted certifiate and private key, add their paths as environment variables in the [.env](.env) file. `LOCAL_CERT` is the path to the directory containing both the certificate and private key. `CERT_NAME` and `KEY_NAME` are the names of the files for the private key and certificate.

## Running

In the root of this project, execute `docker compose up` in the terminal and wait for the docker containers to build and start. If there are no errors, go to this [example](api/example.REST) for a small tutorial on how to use the API.

## API Endpoints

### /api/signup

Add a new user and signup with a post request including a username and password. Responds with 200 OK on success.
#### Example POST request:
```json
{
    "Uname": "your-username",
    "Pass": "your-password"
}
```
### /api/login

Receive a signed jwt token using the same credentials you signed up with. Responds with 200 OK on success and the signed token.
#### Example POST request:
```json
{
    "Uname": "your-username",
    "Pass": "your-password"
}
```
#### Response:
```json
{
    "token": "signed-token"
}
```

### /api/recipe/add

Add a recipe to the database using a jwt token in the authorization header. Responds with 200 ok on success. Ingredients have a name, amount, and unit of measurement (cups, grams, ....). The amount can be a fraction or whole number.
#### Example POST request:
```json
{
    "Name": "recipe-name",
    "Ingredients": [
        {
            "Name": "ingredient name",
            "Amount": "5",
            "Unit": "cups"
        },
        ...
    ],
    "Steps": [
        "step 1",
        "step 2",
        ...
    ],
    "CreatedBy": "your-username"
}
```
### /api/user/example-username

Send a GET request to receive all the recipes associated with 'example-username'. Responds with a list of recipes as a JSON and 200 OK on success.
#### Example response:
```json
[
    {
        "Name": "first recipe",
        ...,
        "CreatedBy": "example-username"
    },
    {
        "Name": "second recipe",
        ...,
        "CreatedBy": "example-username"
    },
    ...
]
```

### /api/user/example-username/recipe-name

Send a GET request to receive the recipe with the name, 'recipe-name', that was created by 'example-username'. Responds with the single recipe as a JSON and 200 OK on success.
#### Example response:
```json
{
    "Name": "recipe-name",
    ...,
    "CreatedBy": "example-username"
},

```

### /api/user/remove

Delete a user and all their recipes from the database. Send a POST request of the login credentials and a signed jwt token in the authorization header. Responds with 200 ok on success.
#### Example POST request:
```json
{
    "Uname": "your-username",
    "Pass": "your-password"
}
```