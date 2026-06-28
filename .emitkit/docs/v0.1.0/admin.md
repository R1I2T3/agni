# Admin

## adminLogin

**POST /api/admin/auth/login**

Admin login; validates username/password against env config

### Parameters

None

### Responses

- **200**: Login successful; admin JWT set in `admin_token` cookie
- **400**: Invalid request
- **401**: Invalid credentials
- **500**: Token generation failed


## adminLogout

**POST /api/admin/logout**

Admin logout; clears the admin_token cookie

### Parameters

None

### Responses

- **200**: Logged out successfully


## adminDashboard

**GET /api/admin/dashboard**

Admin dashboard redirect (returns welcome message and admin username)

### Parameters

None

### Responses

- **200**: Dashboard info
- **401**: Missing or invalid admin token
- **403**: Not an admin


## createApplication

**POST /api/admin/create-application**

Create a new application and generate API token/secret

### Parameters

None

### Responses

- **200**: Application created successfully
- **400**: Invalid request
- **401**: Missing or invalid admin token
- **403**: Not an admin
- **500**: Failed to create application or generate credentials


## getAllApplications

**GET /api/admin/applications**

List all registered applications

### Parameters

None

### Responses

- **200**: List of applications
- **401**: Missing or invalid admin token
- **403**: Not an admin
- **500**: Failed to fetch applications


## regenerateToken

**PUT /api/admin/regenerate-token**

Regenerate API token and secret for an application

### Parameters

None

### Responses

- **200**: Credentials regenerated
- **400**: Invalid request
- **401**: Missing or invalid admin token
- **403**: Not an admin
- **500**: Failed to update credentials


## deleteApplication

**PUT /api/admin/delete-application**

Delete an application by name

### Parameters

None

### Responses

- **200**: Application deleted
- **400**: Invalid request
- **401**: Missing or invalid admin token
- **403**: Not an admin
- **500**: Failed to delete application


