# Notifications

## enqueueNotification

**POST /api/notification/send**

Enqueue a notification for delivery. The ApplicationAuth middleware reads `application_token` and `application_secret` from the same request body before passing control to the handler.


### Parameters

None

### Responses

- **200**: Notification queued
- **400**: Invalid request body
- **401**: Invalid application credentials
- **500**: Failed to enqueue notification


