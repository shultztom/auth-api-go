# auth-api-go

### Routes

#### POST - /login

Body: 

```json
{
    "username": "test",
    "password": "123"
}
```

#### POST - /register

Body:

```json
{
    "username": "test",
    "password": "123"
}
```

#### GET - /verify

Headers:
```
x-auth-token
```
