# User Management API

This API allows you to manage users with CRUD operations.

## API Endpoints

### Create a New User

To create a new user, send a `POST` request to `/api/users` with JSON payload containing `name`, `email`, and `password` fields:

```bash
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic SU9UOjE=' \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securePassword123"
  }'
```

### Retrieve User Information

To retrieve user information by `ID`, send a `GET` request to `/api/users/{id}`:

```bash
# Replace {id} with the actual user ID you want to retrieve
curl -X GET http://localhost:8080/api/users/1 \
  -H 'Authorization: Basic SU9UOjE='
```

### Update User Information

To update user information by `ID`, send a `PUT` request to `/api/users/{id}` with JSON payload containing updated `name`, `email`, and `password` fields:

```bash
# Replace {id} with the actual user ID you want to update
curl -X PUT http://localhost:8080/api/users/1 \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic SU9UOjE=' \
  -d '{
    "name": "Updated Name",
    "email": "updated.email@example.com",
    "password": "newSecurePassword456"
  }'
```

### Delete User

To delete a user by `ID`, send a `DELETE` request to `/api/users/{id}`:

```bash
# Replace {id} with the actual user ID you want to delete
curl -X DELETE http://localhost:8080/api/users/1 \
  -H 'Authorization: Basic SU9UOjE='
```

