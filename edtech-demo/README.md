# EdTech Notification Demo

A standalone demo web application showcasing the **Agni Notification Microservice**.

## Features

- 🔐 **Secure Authentication** - HMAC-based login with JWT tokens
- 📧 **Multi-Channel Delivery** - Email, SMS, and In-App notifications
- 📬 **Live Inbox** - Real-time student notification inbox
- 🎯 **Campaign Composer** - Send notifications across multiple channels
- 🔄 **Real-Time Updates** - WebSocket support for instant delivery

## Project Structure

```
edtech-demo/
├── src/
│   ├── lib/
│   │   ├── api.ts           # Agni microservice API integration
│   │   └── utils.ts         # Utility functions
│   ├── components/
│   │   ├── Button.tsx       # Reusable button component
│   │   ├── Card.tsx         # Reusable card component
│   │   ├── Input.tsx        # Reusable input component
│   │   ├── Textarea.tsx     # Reusable textarea component
│   │   └── NotificationIcon.tsx
│   ├── pages/
│   │   ├── AuthPanel.tsx    # Authentication page
│   │   ├── CampaignSender.tsx # Campaign composer
│   │   └── StudentInbox.tsx # Notification inbox
│   ├── App.tsx              # Main app component
│   ├── main.tsx             # Entry point
│   └── index.css            # Global styles
├── index.html               # HTML template
├── vite.config.ts           # Vite configuration
├── tsconfig.json            # TypeScript configuration
├── package.json             # Dependencies
└── README.md                # This file
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Agni backend running on `http://localhost:8080`
- Agni in-app service running on `ws://localhost:4000`

### Installation

```bash
cd edtech-demo
npm install
```

### Running the App

```bash
npm run dev
```

The app will be available at `http://localhost:3001`

### Building for Production

```bash
npm run build
npm run preview
```

## Usage

1. **Get Credentials**: Use the Agni admin panel to create an application and copy the API token and secret
2. **Sign In**: Paste credentials and a student ID in the Auth panel
3. **Send Notifications**: Use Campaign Sender to send Email, SMS, or In-App notifications
4. **View Inbox**: Monitor real-time notifications in the Student Inbox

## API Integration

The demo integrates with the following Agni endpoints:

- `POST /api/auth/login` - HMAC-based authentication
- `POST /api/notification/send` - Enqueue notifications
- `GET /api/inapp/notifications` - Fetch in-app notifications
- `GET /api/inapp/notifications/unread-count` - Get unread count
- `PUT /api/inapp/notifications/:id/read` - Mark as read
- `WS /ws` - WebSocket for real-time delivery

## Tech Stack

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Lucide React** - Icons
- **Sonner** - Toast notifications
- **date-fns** - Date formatting

## Environment Variables

Create a `.env.local` file:

```
VITE_INAPP_WS_URL=ws://localhost:4000/ws
```

## Demo Scenarios

### Scenario 1: Class Announcement
1. Sign in as a student
2. Select all channels (Email, SMS, In-App)
3. Send a class update message
4. Watch notifications arrive in real-time

### Scenario 2: Multi-Student Campaign
1. Send the same message to multiple student IDs (`student_001`, `student_002`, etc.)
2. Monitor delivery status for each channel
3. Check inbox for in-app notifications

### Scenario 3: Delivery Reliability
1. Send the same notification multiple times
2. Observe retry behavior on delivery
3. Track successful vs failed sends

## License

MIT
