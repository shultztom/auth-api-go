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

#### DELETE - /

Headers:
```
x-auth-token
```

#### GET - /roles

Headers:
```
x-auth-token
```

#### GET - /roles/:role

Headers:
```
x-auth-token
```

#### POST - /roles

Headers:
```
x-auth-token
```

```json
{
  "role": "role-name"
}
```

#### GET - /app/verify

Headers:
```
X-API-Token
```


#### DELETE - /app/user/:username

Headers:
```
X-API-Token
```