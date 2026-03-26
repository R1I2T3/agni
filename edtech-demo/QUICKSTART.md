# EdTech Notification Demo - Quick Start Guide

## What is This?

This is a **completely standalone web application** that demonstrates the Agni notification microservice in action. It's built separately from the existing admin panel and showcases:

- ✅ Email notifications
- ✅ SMS notifications  
- ✅ Real-time In-App notifications via WebSocket
- ✅ Live student inbox with notification tracking

## Demo Walkthrough (3 minutes)

### Step 1: Start Services (Terminal 1 & 2)

```bash
# Terminal 1: Start main Agni server
cd agni
go run cmd/server/main.go

# Terminal 2: Start In-App WebSocket service
cd agni
go run cmd/inapp/main.go
```

Both services should print "running" messages. Keep both terminals open.

### Step 2: Start the Demo App (Terminal 3)

```bash
cd edtech-demo
npm run dev
```

Open http://localhost:3001 in your browser.

### Step 3: Authenticate

1. **Get credentials from admin panel**:
   - Go to http://localhost:3000 (existing admin panel)
   - Create an application or copy existing credentials
   - Note: `api_token` and `api_secret`

2. **Authenticate in demo app**:
   - Paste API token in "Application Token" field
   - Paste API secret in "Application Secret" field
   - Enter student ID (e.g., `student_001`)
   - Click "Sign In"

### Step 4: Send a Notification

1. **Select Channels**: Choose Email, SMS, In-App (or any combination)
2. **Fill Message**:
   - Recipient: `student_001` (or any ID)
   - Subject: "Class Announcement"
   - Message: "Your assignment is due tomorrow!"
3. **Click "Send Notifications"**

### Step 5: Watch It Work

- **Delivery Monitor** (in Campaign Sender card): Shows queue progression
- **Student Inbox** (right panel): Live in-app notifications appear in real-time
- Each notification shows: channel, status, timestamp, and read state
- Click "Mark Read" to update read status

## Project Structure

```
edtech-demo/
├── src/
│   ├── lib/
│   │   ├── api.ts              # Agni API integration (HMAC auth, send, fetch)
│   │   └── utils.ts            # Helper functions
│   ├── components/
│   │   ├── Button.tsx          # Reusable button
│   │   ├── Card.tsx            # Card container
│   │   ├── Input.tsx           # Input field
│   │   ├── Textarea.tsx        # Textarea field
│   │   └── NotificationIcon.tsx # Channel icons
│   ├── pages/
│   │   ├── AuthPanel.tsx       # Login/auth form
│   │   ├── CampaignSender.tsx  # Send notifications panel
│   │   └── StudentInbox.tsx    # View notifications panel
│   ├── App.tsx                 # Main app component
│   └── main.tsx                # Entry point
├── index.html                  # HTML template
├── vite.config.ts              # Vite build config
├── tsconfig.json               # TypeScript config
├── package.json                # Dependencies
└── README.md                   # Full documentation
```

## Key Features Explained

### 1. **Campaign Composer** (Left Panel)
- Select multiple notification channels simultaneously
- Set recipient, subject, and message
- Send button triggers real API calls to `/api/notification/send`
- Shows success/error toast for each channel

### 2. **Student Inbox** (Right Panel)
- Fetches in-app notifications via `/api/inapp/notifications`
- Shows unread count
- Mark individual notifications as read
- Auto-refreshes when new In-App notifications are sent
- Live timestamps show how long ago notifications arrived

### 3. **Authentication Panel** (Top)
- HMAC-SHA256 signature generation for secure auth
- Fetches JWT token via `/api/auth/login`
- Token stored in HTTP-only cookie (secure)
- Cookie automatically sent to protected endpoints

## Demonstration Scenarios

### Scenario 1: Multi-Channel Delivery
```
1. Select: Email + SMS + In-App
2. Recipient: student_001
3. Message: "Your grades are posted"
4. Send
→ Each channel processes independently
→ In-App arrives instantly in inbox
→ Email/SMS queued in backend (visible in logs)
```

### Scenario 2: Reliability Demo
```
1. Send same notification 3 times
2. Watch delivery timeline show progress
3. In-App inbox shows all 3 notifications
4. Each with read/unread state
5. Demonstrates queue handling and deduplication
```

### Scenario 3: Real-Time Stream
```
1. Keep inbox open
2. In another tab, send In-App notification
3. Watch inbox refresh automatically
4. Demonstrates WebSocket real-time delivery
```

## API Endpoints Used

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/auth/login` | POST | Get JWT via HMAC signature |
| `/api/notification/send` | POST | Queue notification |
| `/api/inapp/notifications` | GET | Fetch user's notifications |
| `/api/inapp/notifications/:id/read` | PUT | Mark as read |
| `/api/inapp/notifications/unread-count` | GET | Get unread count |
| `/ws` | WebSocket | Real-time in-app stream |

## Development

### Build for production
```bash
npm run build
```

### Preview production build
```bash
npm run preview
```

### Lint code
```bash
npm run lint
```

### Format code
```bash
npm run format
```

## Tech Stack

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool & dev server
- **Tailwind CSS** - Styling
- **Lucide React** - Icons
- **Sonner** - Toast notifications
- **date-fns** - Date formatting

## Environment Variables

Create `.env.local` (optional):

```
VITE_INAPP_WS_URL=ws://localhost:4000/ws
```

Default: Auto-detects from window location

## Troubleshooting

### "Authentication failed"
- ✓ Check token and secret are copied correctly (no spaces)
- ✓ Verify Agni backend is running on `http://localhost:8080`
- ✓ Check user ID format (alphanumeric)

### "WebSocket connection failed"
- ✓ Verify in-app service running on `http://localhost:4000`
- ✓ Check browser console for connection errors
- ✓ Notifications may still send (check via GET inbox endpoint)

### "No notifications in inbox"
- ✓ Make sure In-App channel is selected when sending
- ✓ Click "Refresh Inbox" button
- ✓ Check recipient ID matches user ID used for auth

### Build errors
```bash
# Clear cache and reinstall
rm -rf node_modules dist
npm install
npm run build
```

## Next Steps

1. **Try different channels**: Focus on Email, then SMS, then In-App
2. **Test reliability**: Send same message multiple times
3. **Explore real-time**: Send notifications while inbox is open
4. **Check logs**: Backend logs show queue processing and delivery details
5. **Scale**: Send to multiple recipient IDs

## Questions?

Refer to the README.md in the demo root directory for full documentation.
