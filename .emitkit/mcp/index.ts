// Auto-generated MCP Server
// Do not edit manually

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";

const server = new Server(
  {
    name: "Agni Notification Engine API",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// Register list of tools
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
            "name": "rootWelcome",
            "description": "Root welcome endpoint",
            "inputSchema": {
                  "type": "object",
                  "properties": {}
            }
      },
      {
            "name": "healthCheck",
            "description": "Health check for Redis and MySQL connectivity",
            "inputSchema": {
                  "type": "object",
                  "properties": {}
            }
      },
      {
            "name": "hmacLogin",
            "description": "HMAC-based login for client applications; returns a JWT in an HTTP-only cookie",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "body": {
                              "type": "object",
                              "description": "Request body content"
                        }
                  },
                  "required": [
                        "body"
                  ]
            }
      },
      {
            "name": "adminLogin",
            "description": "Admin login; validates username/password against env config",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "body": {
                              "type": "object",
                              "description": "Request body content"
                        }
                  },
                  "required": [
                        "body"
                  ]
            }
      },
      {
            "name": "adminLogout",
            "description": "Admin logout; clears the admin_token cookie",
            "inputSchema": {
                  "type": "object",
                  "properties": {}
            }
      },
      {
            "name": "adminDashboard",
            "description": "Admin dashboard redirect (returns welcome message and admin username)",
            "inputSchema": {
                  "type": "object",
                  "properties": {}
            }
      },
      {
            "name": "createApplication",
            "description": "Create a new application and generate API token/secret",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "body": {
                              "type": "object",
                              "description": "Request body content"
                        }
                  },
                  "required": [
                        "body"
                  ]
            }
      },
      {
            "name": "getAllApplications",
            "description": "List all registered applications",
            "inputSchema": {
                  "type": "object",
                  "properties": {}
            }
      },
      {
            "name": "regenerateToken",
            "description": "Regenerate API token and secret for an application",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "body": {
                              "type": "object",
                              "description": "Request body content"
                        }
                  },
                  "required": [
                        "body"
                  ]
            }
      },
      {
            "name": "deleteApplication",
            "description": "Delete an application by name",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "body": {
                              "type": "object",
                              "description": "Request body content"
                        }
                  },
                  "required": [
                        "body"
                  ]
            }
      },
      {
            "name": "enqueueNotification",
            "description": "Enqueue a notification for delivery. The ApplicationAuth middleware reads `application_token` and `application_secret` from the same request body before passing control to the handler.\n",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "body": {
                              "type": "object",
                              "description": "Request body content"
                        }
                  },
                  "required": [
                        "body"
                  ]
            }
      },
      {
            "name": "getInAppNotifications",
            "description": "Retrieve in-app notifications for the authenticated user",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "unread_only": {
                              "type": "string",
                              "description": "Filter to unread notifications only"
                        },
                        "limit": {
                              "type": "number",
                              "description": "Maximum number of notifications to return"
                        },
                        "offset": {
                              "type": "number",
                              "description": "Number of notifications to skip"
                        }
                  }
            }
      },
      {
            "name": "getUnreadCount",
            "description": "Get count of unread in-app notifications for the authenticated user",
            "inputSchema": {
                  "type": "object",
                  "properties": {}
            }
      },
      {
            "name": "markNotificationAsRead",
            "description": "Mark a single notification as read",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "id": {
                              "type": "string",
                              "description": "Notification ID"
                        }
                  },
                  "required": [
                        "id"
                  ]
            }
      },
      {
            "name": "markAllNotificationsAsRead",
            "description": "Mark all in-app notifications as read for the authenticated user",
            "inputSchema": {
                  "type": "object",
                  "properties": {}
            }
      },
      {
            "name": "handleWebPushSubscription",
            "description": "Register a WebPush subscription endpoint",
            "inputSchema": {
                  "type": "object",
                  "properties": {
                        "body": {
                              "type": "object",
                              "description": "Request body content"
                        }
                  },
                  "required": [
                        "body"
                  ]
            }
      }
],
  };
});

// Register tool calls
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  switch (name) {
    case "rootWelcome": {
      // Mock / dynamic implementation of rootWelcome
      console.error("Calling tool rootWelcome with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool rootWelcome",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "healthCheck": {
      // Mock / dynamic implementation of healthCheck
      console.error("Calling tool healthCheck with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool healthCheck",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "hmacLogin": {
      // Mock / dynamic implementation of hmacLogin
      console.error("Calling tool hmacLogin with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool hmacLogin",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "adminLogin": {
      // Mock / dynamic implementation of adminLogin
      console.error("Calling tool adminLogin with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool adminLogin",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "adminLogout": {
      // Mock / dynamic implementation of adminLogout
      console.error("Calling tool adminLogout with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool adminLogout",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "adminDashboard": {
      // Mock / dynamic implementation of adminDashboard
      console.error("Calling tool adminDashboard with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool adminDashboard",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "createApplication": {
      // Mock / dynamic implementation of createApplication
      console.error("Calling tool createApplication with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool createApplication",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "getAllApplications": {
      // Mock / dynamic implementation of getAllApplications
      console.error("Calling tool getAllApplications with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool getAllApplications",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "regenerateToken": {
      // Mock / dynamic implementation of regenerateToken
      console.error("Calling tool regenerateToken with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool regenerateToken",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "deleteApplication": {
      // Mock / dynamic implementation of deleteApplication
      console.error("Calling tool deleteApplication with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool deleteApplication",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "enqueueNotification": {
      // Mock / dynamic implementation of enqueueNotification
      console.error("Calling tool enqueueNotification with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool enqueueNotification",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "getInAppNotifications": {
      // Mock / dynamic implementation of getInAppNotifications
      console.error("Calling tool getInAppNotifications with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool getInAppNotifications",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "getUnreadCount": {
      // Mock / dynamic implementation of getUnreadCount
      console.error("Calling tool getUnreadCount with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool getUnreadCount",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "markNotificationAsRead": {
      // Mock / dynamic implementation of markNotificationAsRead
      console.error("Calling tool markNotificationAsRead with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool markNotificationAsRead",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "markAllNotificationsAsRead": {
      // Mock / dynamic implementation of markAllNotificationsAsRead
      console.error("Calling tool markAllNotificationsAsRead with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool markAllNotificationsAsRead",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }

    case "handleWebPushSubscription": {
      // Mock / dynamic implementation of handleWebPushSubscription
      console.error("Calling tool handleWebPushSubscription with args:", args);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              message: "Mock response from tool handleWebPushSubscription",
              received_arguments: args,
            }, null, 2),
          },
        ],
      };
    }
    default:
      throw new Error(`Unknown tool: ${name}`);
  }
});

async function run() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("MCP Server running on stdio");
}

run().catch((error) => {
  console.error("Fatal error running MCP Server:", error);
  process.exit(1);
});
