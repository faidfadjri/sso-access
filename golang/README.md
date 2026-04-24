# SSO Access API Documentation

This document outlines the available endpoints for the SSO Access API.

## Table of Contents
- [OAuth](#oauth)
- [Services](#services)
- [Users](#users)
- [Roles](#roles)
- [Permissions](#permissions)
- [Assigned Roles](#assigned-roles)

---

## OAuth

### Forgot Password
- **URL**: `/oauth/forgot-password`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
      "email": "faidfadjri@gmail.com",
      "forgot_type": "password"
  }
  ```

### Reset Password
- **URL**: `/oauth/reset-password`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
    "token": "vFTFMBDe7ct1r4gU9o8CMyMGU0UNwS6T",
    "password": "NewSecurePassword!",
    "password_confirmation": "NewSecurePassword!"
  }
  ```

### Login
- **URL**: `/oauth/login`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
      "email_or_username" : "faidfadjri",
      "password" : "bismillah"
  }
  ```

### Refresh Token
- **URL**: `/auth/refresh-token`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

### Authorize
- **URL**: `/oauth/authorize`
- **Method**: `GET`
- **Authorization**: Bearer Token
- **Query Params**:
  - `client_id`: (Required) string
  - `redirect_uri`: (Required) string
  - `response_type`: (Required) string e.g. `code`
  - `scope`: (Required) string e.g. `read+write`
  - `state`: (Optional) string
  - `code_challenge`: (Optional) string
  - `code_challenge_method`: (Optional) string e.g. `S256`

### Token Exchange
- **URL**: `/oauth/token`
- **Method**: `POST`
- **Body** (JSON):
  ```json
  {
      "grant_type": "authorization_code",
      "code": "trkknSlNzFdAKaB18FvzqFprlx00bHIU",
      "client_id": "4597ae7856159c297907c70452d916adc6f2f6de23ab8aa7aee79f7baffa8464",
      "redirect_uri": "http://localhost:3000",
      "code_verifier": "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
  }
  ```

### Me
- **URL**: `/oauth/me`
- **Method**: `GET`
- **Authorization**: Bearer Token

### Logout
- **URL**: `/oauth/logout`
- **Method**: `POST`
- **Authorization**: Bearer Token

### Update Account
- **URL**: `/oauth/update-account`
- **Method**: `PUT`
- **Authorization**: Bearer Token
- **Body** (Form-Data):
  - `full_name`: string
  - `email`: string
  - `username`: string
  - `phone`: string
  - `photo`: file (Optional)
  - `password`: string (Optional)
  - `password_confirmation`: string (Required if password changes)

---

## Services

### Create Service Client
- **URL**: `/service/clients`
- **Method**: `POST`
- **Authorization**: Bearer Token
- **Body** (Form-Data):
  - `name`: string
  - `description`: string
  - `logo`: file
  - `redirect_url`: string

### Get Service Client
- **URL**: `/service/clients`
- **Method**: `GET`
- **Authorization**: Bearer Token
- **Query Params**:
  - `show`: number (Default 6)
  - `page`: number (Default 1)
  - `sort`: string (Default `desc`)
  - `search`: string

### Update Service Client
- **URL**: `/service/clients/:id`
- **Method**: `PUT`
- **Authorization**: Bearer Token
- **Body** (Form-Data):
  - `name`: string
  - `description`: string
  - `logo`: file
  - `redirect_url`: string

### Delete Service Client
- **URL**: `/service/clients/:id`
- **Method**: `DELETE`
- **Authorization**: Bearer Token

---

## Users

### Get User
- **URL**: `/users`
- **Method**: `GET`
- **Authorization**: Bearer Token
- **Query Params**:
  - `show`: number (Default 6, -1 for all)
  - `page`: number (Default 1)
  - `sort`: string (Default `desc`)
  - `search`: string

### Create User
- **URL**: `/users`
- **Method**: `POST`
- **Authorization**: Bearer Token
- **Body** (Form-Data):
  - `full_name`: string
  - `email`: string
  - `username`: string
  - `phone`: string
  - `password`: string
  - `admin`: boolean (`true` or `false`)
  - `photo`: file (Optional)

### Update User
- **URL**: `/users/:id`
- **Method**: `PUT`
- **Authorization**: Bearer Token
- **Body** (Form-Data):
  - `full_name`: string
  - `email`: string
  - `username`: string
  - `phone`: string
  - `admin`: boolean (`true` or `false`)
  - `password`: string (Optional)
  - `password_confirmation`: string (Required if password changes)
  - `photo`: file (Optional)

### Delete User
- **URL**: `/users/:id`
- **Method**: `DELETE`
- **Authorization**: Bearer Token

---

## Roles

### Get Roles
- **URL**: `/roles`
- **Method**: `GET`
- **Authorization**: Bearer Token
- **Query Params**:
  - `show`: number (Default 6)
  - `page`: number (Default 1)
  - `sort`: string (Default `desc`)
  - `search`: string

### Create Roles
- **URL**: `/roles`
- **Method**: `POST`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "service_id" : 1,
      "role_name": "SA Prospect"
  }
  ```

### Update Roles
- **URL**: `/roles/:id`
- **Method**: `PUT`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "service_id" : 1,
      "role_name": "SA"
  }
  ```

### Delete Roles
- **URL**: `/roles/:id`
- **Method**: `DELETE`
- **Authorization**: Bearer Token

---

## Permissions

### Get Permission
- **URL**: `/permissions`
- **Method**: `GET`
- **Authorization**: Bearer Token
- **Query Params**:
  - `show`: number (Default 6)
  - `page`: number (Default 1)
  - `sort`: string (Default `desc`)
  - `search`: string

### Create Permission
- **URL**: `/permissions`
- **Method**: `POST`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "permission_key": "task.all",
      "description": "CRUD for task"
  }
  ```

### Update Permission
- **URL**: `/permissions/:id`
- **Method**: `PUT`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "permission_key": "task.all",
      "description": "CRUD task"
  }
  ```

### Delete Permission
- **URL**: `/permissions/:id`
- **Method**: `DELETE`
- **Authorization**: Bearer Token

---

## User Access

### Get User Access
- **URL**: `/users/access`
- **Method**: `GET`
- **Authorization**: Bearer Token
- **Query Params**:
  - `show`: number (Default 6)
  - `page`: number (Default 1)
  - `sort`: string (Default `desc`)
  - `search`: string
  - `user_id`: number

### Create User Access
- **URL**: `/users/access`
- **Method**: `POST`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "user_id" : 1,
      "service_ids" : [1],
      "status": "active" 
  }
  ```

### Update User Access
- **URL**: `/users/access/:id`
- **Method**: `PUT`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "user_id" : 1,
      "service_ids" : [24,25,27],
      "status": "active" 
  }
  ```

### Delete User Access
- **URL**: `/users/access/:id`
- **Method**: `DELETE`
- **Authorization**: Bearer Token

---

## Assigned Roles

### Get List
- **URL**: `/roles/assign`
- **Method**: `GET`
- **Authorization**: Bearer Token
- **Query Params**:
  - `search`: string
  - `role_id`: number
  - `service_id`: number
  - `show`: number

### Assign Role
- **URL**: `/roles/assign`
- **Method**: `POST`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "user_ids": [2,3],
      "role_id": 2
  }
  ```

### Delete Assigned Role
- **URL**: `/roles/assign`
- **Method**: `DELETE`
- **Authorization**: Bearer Token
- **Body** (JSON):
  ```json
  {
      "user_id": 2,
      "role_id": 2,
      "service_id": 1
  }
  ```
