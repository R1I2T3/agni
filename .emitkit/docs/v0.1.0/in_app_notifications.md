# In-App Notifications

## getInAppNotifications

**GET /api/inapp/notifications**

Retrieve in-app notifications for the authenticated user

### Parameters

- **unread_only** (`query`): Filter to unread notifications only
- **limit** (`query`): Maximum number of notifications to return
- **offset** (`query`): Number of notifications to skip

### Responses

- **200**: Paginated list of notifications
- **401**: Unauthorized (missing or invalid JWT)
- **500**: Failed to fetch notifications


## getUnreadCount

**GET /api/inapp/notifications/unread-count**

Get count of unread in-app notifications for the authenticated user

### Parameters

None

### Responses

- **200**: Unread count
- **401**: Unauthorized
- **500**: Failed to count notifications


## markNotificationAsRead

**PUT /api/inapp/notifications/{id}/read**

Mark a single notification as read

### Parameters

- **id** (`path`): Notification ID

### Responses

- **200**: Notification marked as read
- **400**: Missing or invalid notification ID
- **401**: Unauthorized
- **404**: Notification not found
- **500**: Database update failed


## markAllNotificationsAsRead

**PUT /api/inapp/notifications/read-all**

Mark all in-app notifications as read for the authenticated user

### Parameters

None

### Responses

- **200**: All notifications marked as read
- **401**: Unauthorized
- **500**: Database update failed


