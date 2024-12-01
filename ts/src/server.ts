import WebSocket from 'ws';
import express from 'express';
import { createServer } from 'http';
import path from 'path';
import sqlite3 from 'sqlite3';
import { open } from 'sqlite';
import { v4 as uuidv4 } from 'uuid';
import fs from 'fs/promises';
import { existsSync } from 'fs';
import { exec } from 'child_process';
import util from 'util';

const app = express();
const server = createServer(app);
const PORT = process.env.PORT || 8020;
const wss = new WebSocket.Server({ 
  server,
  path: '/ws'
});
const channels = ["nixos-23.05", "nixos-23.11", "nixos-24.05", "nixos-unstable"];

// Enable JSON body parsing
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Serve static files from 'public' directory
app.use(express.static(path.join(__dirname, '../public')));

// Add new type definition
type ClientType = 'go' | 'browser';

// Modify clients to store client info
const clients = new Map<string, Client>();

// Initialize SQLite database
let db: any;
(async () => {
  db = await open({
    filename: 'clients.db',
    driver: sqlite3.Database
  });

  // Create clients table if it doesn't exist
  await db.exec(`
    CREATE TABLE IF NOT EXISTS clients (
      id TEXT PRIMARY KEY,
      client_type TEXT,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `);
})();

// Add new routes
app.get('/clients', async (req, res) => {
  try {
    const rows = await db.all('SELECT * FROM clients');
    res.json(rows);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch clients' });
  }
});

// Add new type definitions
type CommandMessage = {
  type: 'command';
  command: string;
  files?: Record<string, string>;
  message?: string;
};

type ResponseMessage = {
  type: 'response' | 'error';
  command: string;
  message: string;
  timestamp: string;
};

// HTTP route to send message
app.post('/send-message', (req, res) => {
  const { message } = req.body;
  
  if (!message) {
    return res.status(400).json({ error: 'Message is required' });
  }

  // Broadcast message to all WebSocket clients
  const messageData = {
    type: 'message',
    message: message,
    timestamp: new Date().toISOString()
  };

  clients.forEach(client => {
    if (client.ws.readyState === WebSocket.OPEN) {
      client.ws.send(JSON.stringify(messageData));
    }
  });

  res.json({ success: true, message: 'Message sent' });
});

// Add new type definitions
type RegisterMessage = {
  type: 'register';
  id: string;
};

// Add new type definitions at the top with other type definitions
type Message = {
  type: string;
  message?: string;
  command?: string;
  files?: Record<string, string>;
  timestamp?: string;
  id?: string;
  targetId?: string;
  excludedClients?: string[];
};

// Add new function to handle file operations
async function handleNixFiles(files: Record<string, string>) {
  const nixConfigDir = process.env.NIX_CONFIG_DIR || '/etc/nixos';
  
  // Ensure directory exists
  try {
    await fs.access(nixConfigDir);
  } catch {
    await fs.mkdir(nixConfigDir, { recursive: true });
  }
  
  // Process each file
  for (const [filename, content] of Object.entries(files)) {
    const filepath = `${nixConfigDir}/${filename}`;
    
    // Create backup if file exists
    if (existsSync(filepath)) {
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
      const backupPath = `${filepath}.backup-${timestamp}`;
      await fs.copyFile(filepath, backupPath);
    }
    
    // Write new content
    await fs.writeFile(filepath, content);
  }
}

// Modify Client type to include client type
type Client = {
  ws: WebSocket;
  id: string;
  type: ClientType;
};

// Add new route to get Go clients
app.get('/go-clients', async (req, res) => {
  try {
    const goClients = Array.from(clients.entries())
      .filter(([_, client]) => client.type === 'go')
      .map(([id, _]) => id);
    res.json(goClients);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch Go clients' });
  }
});

// Add new promisified exec
const execAsync = util.promisify(exec);

// Add new type definition for package search
type PackageSearchMessage = {
  type: 'package-search';
  query: string;
  targetId?: string;
};

app.get('/search-options/:query', async (req, res) => {
  const { query } = req.params;

  try {
    const channel = "nixos-unstable"; // or whichever channel you want to search
    const db = await open({
      filename: `options-${channel}.db`,
      driver: sqlite3.Database
    });

    const rows = await db.all(`
      SELECT 
        attribute_name,
        name,
        long_description,
        example,
        \`default\`,
        type
      FROM options 
      WHERE options MATCH ?
    `, query);

    await db.close();

    if (rows.length === 0) {
      return res.status(404).json({ error: 'No matching options found' });
    }

    const formattedOutput = rows.map(row => `
### Option: ${row.attribute_name}

Name: ${row.name}

Description: 
${row.long_description}

Default:
\`\`\`
${row.default}  
\`\`\`

Example:
\`\`\`
${row.example}
\`\`\`

Type: ${row.type}
    `).join('\n---\n');

    res.json({
      query,
      formatted: formattedOutput,
      raw: rows
    });
  } catch (error) {
    console.error('Error searching options:', error);
    res.status(500).json({ error: 'Failed to search options' });
  }
});

// Add new route for package search with formatted output
app.get('/search-package/:query', async (req, res) => {
  try {
    const { query } = req.params;
    const { stdout, stderr } = await execAsync(`nix-search ${query}`);
    
    if (stderr) {
      return res.status(500).json({ error: stderr });
    }

    // Format the output for LLM readability
    const formattedOutput = `
### Nix Package Search Results for "${query}"

\`\`\`
${stdout.trim()}
\`\`\`

The above results show available Nix packages matching the query "${query}". 
Each line typically contains: package name, version, and description.
    `.trim();
    
    res.json({ 
      query,
      formatted: formattedOutput,
      raw: stdout 
    });
  } catch (error: unknown) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
    res.status(500).json({ error: 'Failed to search packages: ' + errorMessage });
  }
});

// WebSocket connection handler
wss.on('connection', async (ws, req) => {
  console.log(`New WebSocket connection from ${req.socket.remoteAddress}`);
  
  const generatedId = uuidv4();
  ws.send(JSON.stringify({
    type: 'system',
    message: 'Connected to server, waiting for registration...',
    id: generatedId
  }));
  
  ws.once('message', async (message) => {
    try {
      const data = JSON.parse(message.toString()) as Message;
      
      if (data.type === 'register') {
        const clientId = data.id || generatedId;
        const clientType: ClientType = data.id ? 'go' : 'browser';
        
        // Store client connection
        clients.set(clientId, { ws, id: clientId, type: clientType });
        
        // Update database
        try {
          await db.run(
            'INSERT OR REPLACE INTO clients (id, client_type) VALUES (?, ?)', 
            [clientId, clientType]
          );
        } catch (error) {
          console.error('Error updating client in database:', error);
        }
        
        console.log(`${clientType} client registered with ID: ${clientId}`);
        
        // Send confirmation
        const response: any = {
          type: 'system',
          message: `Successfully registered with ID: ${clientId}`,
          id: clientId
        };
        
        ws.send(JSON.stringify(response));
        
        ws.on('message', (message) => handleMessage(message, clientId));
      } else {
        console.error('First message was not a registration');
        ws.close();
      }
    } catch (error) {
      console.error('Error processing registration:', error);
      ws.close();
    }
  });

  ws.on('close', async () => {
    // Find and remove client
    for (const [id, client] of clients.entries()) {
      if (client.ws === ws) {
        const clientType = client.type;
        clients.delete(id);
        console.log(`Client disconnected: ${id}`);
        
        // Remove from database
        try {
          await db.run('DELETE FROM clients WHERE id = ?', id);
        } catch (error) {
          console.error('Error removing client from database:', error);
        }
        
        break;
      }
    }
  });
});

// Update message handling function to handle both message types properly
function handleMessage(message: any, senderId: string) {
  try {
    const data = JSON.parse(message.toString()) as Message | PackageSearchMessage;
    console.log('Received:', data);

    // Add package search handling
    if (data.type === 'package-search' && 'query' in data) {
      execAsync(`optinix get ${data.query} --no-tui`)
        .then(({ stdout, stderr }) => {
          // Format the output for LLM readability
          const formattedOutput = `
### Nix Package Search Results for "${data.query}"

\`\`\`
${stdout.trim()}
\`\`\`

The above results show available Nix packages matching the query "${data.query}". 
Each line typically contains: package name, version, and description.
          `.trim();

          const response: ResponseMessage = {
            type: 'response',
            command: 'package-search',
            message: formattedOutput,
            timestamp: new Date().toISOString()
          };
          
          const client = clients.get(senderId);
          if (client && client.ws.readyState === WebSocket.OPEN) {
            client.ws.send(JSON.stringify(response));
          }
        })
        .catch(error => {
          const errorResponse: ResponseMessage = {
            type: 'error',
            command: 'package-search',
            message: `Error searching packages: ${error instanceof Error ? error.message : 'Unknown error'}`,
            timestamp: new Date().toISOString()
          };
          
          const client = clients.get(senderId);
          if (client && client.ws.readyState === WebSocket.OPEN) {
            client.ws.send(JSON.stringify(errorResponse));
          }
        });
      return;
    }

    // Handle set command with files
    if (data.type === 'command' && data.command === 'set' && data.files) {
      handleNixFiles(data.files)
        .then(() => {
          const response: ResponseMessage = {
            type: 'response',
            command: 'set',
            message: 'Files updated successfully',
            timestamp: new Date().toISOString()
          };
          
          // Send response back to sender
          const client = clients.get(senderId);
          if (client && client.ws.readyState === WebSocket.OPEN) {
            client.ws.send(JSON.stringify(response));
          }
        })
        .catch(error => {
          const errorResponse: ResponseMessage = {
            type: 'error',
            command: 'set',
            message: `Error updating files: ${error.message}`,
            timestamp: new Date().toISOString()
          };
          
          const client = clients.get(senderId);
          if (client && client.ws.readyState === WebSocket.OPEN) {
            client.ws.send(JSON.stringify(errorResponse));
          }
        });
      return;
    }

    // Add targetId to message types
    if (data.targetId) {
      // Send to specific client
      const targetClient = clients.get(data.targetId);
      if (targetClient && targetClient.ws.readyState === WebSocket.OPEN) {
        targetClient.ws.send(JSON.stringify(data));
      }
    } else {
      // Broadcast to all clients
      clients.forEach(client => {
        if (client.ws.readyState === WebSocket.OPEN) {
          client.ws.send(JSON.stringify(data));
        }
      });
    }

    if (data.type === 'message') {
      // Add timestamp if not present
      if (!data.timestamp) {
        data.timestamp = new Date().toISOString();
      }

      // Send to all non-excluded clients
      clients.forEach((client, id) => {
        if (
          client.ws.readyState === WebSocket.OPEN &&
          (!data.excludedClients || !data.excludedClients.includes(id))
        ) {
          client.ws.send(JSON.stringify(data));
        }
      });
    }
  } catch (error) {
    console.error('Error processing message:', error);
  }
}

server.listen(PORT, () => {
  console.log(`HTTP server running on http://localhost:${PORT}`);
  console.log(`WebSocket server running on ws://localhost:${PORT}/ws`);
}); 