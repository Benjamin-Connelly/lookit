// lookit server entry point
const http = require('http');
const https = require('https');
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const MarkdownIt = require('markdown-it');
const hljs = require('highlight.js');
const mdHighlight = require('markdown-it-highlightjs');
const { parseArgs, printCertInstructions, ensureCerts, openBrowser, escapeHtml } = require('./utils');
const { handleFile } = require('./fileHandler');
const { registerInstance, unregisterInstance, cleanStaleInstances } = require('./instanceManager');
const { handleListCommand, handleStopCommand, handleStopAllCommand } = require('./commands');

// Parse CLI arguments
const args = parseArgs(process.argv.slice(2));

// Handle management commands (don't start server)
if (args.list) {
  handleListCommand();
  process.exit(0);
}

if (args.stopAll) {
  handleStopAllCommand();
  process.exit(0);
}

if (args.stop !== undefined) {
  handleStopCommand(args.stop);
  process.exit(0);
}

// Clean stale instances before starting server
cleanStaleInstances();
const PORT = args.port || 7777;
const HOST = args.host || '127.0.0.1';
const CWD = process.cwd();
const CERT_DIR = path.join(process.env.HOME, '.config', 'lookit');
const CERT_PATH = args.cert || path.join(CERT_DIR, 'localhost.pem');
const KEY_PATH = args.key || path.join(CERT_DIR, 'localhost-key.pem');

function generateETag(content) {
  return '"' + crypto.createHash('md5').update(content).digest('hex').slice(0, 16) + '"';
}

function setSecurityHeaders(res) {
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'SAMEORIGIN');
  res.setHeader('Referrer-Policy', 'no-referrer-when-downgrade');
}

// Initialize markdown renderer with syntax highlighting
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
}).use(mdHighlight, { hljs });

// Check for TLS certificates (auto-generate with mkcert if available)
const hasCerts = args.noHttps ? false : ensureCerts(CERT_PATH, KEY_PATH, args.quiet);
const useHttps = hasCerts && !args.noHttps;

if (args.httpsOnly && !hasCerts) {
  console.error('❌ Error: --https-only specified but certificates could not be found or generated.\n');
  printCertInstructions();
  process.exit(1);
}

// Track the port we're running on
let runningPort = null;

// WebSocket clients for live reload
const wsClients = new Set();

function handleWebSocketUpgrade(req, socket, head) {
  if (req.url !== '/__ws') {
    socket.destroy();
    return;
  }

  const key = req.headers['sec-websocket-key'];
  const accept = crypto.createHash('sha1')
    .update(key + '258EAFA5-E914-47DA-95CA-5AB9AFA7D59B')
    .digest('base64');

  socket.write(
    'HTTP/1.1 101 Switching Protocols\r\n' +
    'Upgrade: websocket\r\n' +
    'Connection: Upgrade\r\n' +
    'Sec-WebSocket-Accept: ' + accept + '\r\n' +
    '\r\n'
  );

  wsClients.add(socket);
  socket.on('close', () => wsClients.delete(socket));
  socket.on('error', () => wsClients.delete(socket));
}

function notifyClients(filePath) {
  const frame = Buffer.alloc(2 + filePath.length);
  frame[0] = 0x81; // text frame
  frame[1] = filePath.length;
  frame.write(filePath, 2);

  for (const client of wsClients) {
    try {
      client.write(frame);
    } catch {
      wsClients.delete(client);
    }
  }
}

function watchDirectory(dir) {
  try {
    const watcher = fs.watch(dir, { recursive: true }, (eventType, filename) => {
      if (filename && /\.(md|mdx)$/i.test(filename)) {
        const relPath = '/' + filename.replace(/\\/g, '/');
        notifyClients(relPath);
      }
    });
    watcher.on('error', () => {}); // Ignore watch errors
  } catch {
    // fs.watch not supported or permission denied
  }
}

// Start server with auto-increment port logic
function startServer(port, retriesLeft = 10) {
  const server = useHttps ? createHttpsServer() : createHttpServer();

  server.on('upgrade', handleWebSocketUpgrade);

  server.listen(port, HOST, () => {
    runningPort = port;
    const protocol = useHttps ? 'https' : 'http';
    const url = `${protocol}://${HOST}:${port}`;

    registerInstance(port, CWD, protocol);

    console.log(`👀 lookit - Code Browser`);
    console.log(`📂 Serving: ${CWD}`);
    console.log(`🌐 Address:  ${url}`);
    console.log(`🔒 Security: ${useHttps ? 'HTTPS (TLS)' : 'HTTP (plaintext)'}`);
    if (port === 7777) {
      console.log(`🍀 Lucky port: ${port}`);
    }
    console.log('\nPress Ctrl+C to stop.\n');

    if (args.open) {
      openBrowser(url);
    }

    // Watch for markdown file changes (for live reload)
    watchDirectory(CWD);
  });

  server.on('error', (err) => {
    if (err.code === 'EADDRINUSE' && retriesLeft > 0) {
      console.log(`⏭️  Port ${port} in use, trying ${port + 1}...`);
      startServer(port + 1, retriesLeft - 1);
    } else if (err.code === 'EADDRINUSE') {
      console.error(`❌ Error: Ports ${PORT}-${port} all in use.`);
      console.error(`   Stop existing instances: lookit --stop-all`);
      console.error(`   Or specify a higher port: lookit --port ${port + 10}`);
      process.exit(1);
    } else {
      console.error(`❌ Server error: ${err.message}`);
      process.exit(1);
    }
  });
}

// Cleanup on exit
process.on('SIGINT', () => {
  console.log('\n👋 Shutting down lookit...');
  if (runningPort) {
    unregisterInstance(runningPort);
  }
  process.exit(0);
});

process.on('SIGTERM', () => {
  if (runningPort) {
    unregisterInstance(runningPort);
  }
  process.exit(0);
});

// Cleanup on uncaught errors
process.on('uncaughtException', (err) => {
  console.error('❌ Uncaught error:', err);
  if (runningPort) {
    unregisterInstance(runningPort);
  }
  process.exit(1);
});

// Start from initial port
startServer(PORT);

function createHttpsServer() {
  const options = {
    key: fs.readFileSync(KEY_PATH),
    cert: fs.readFileSync(CERT_PATH)
  };
  return https.createServer(options, handleRequest);
}

function createHttpServer() {
  return http.createServer(handleRequest);
}

function handleRequest(req, res) {
  setSecurityHeaders(res);

  const urlPath = decodeURIComponent(req.url.split('?')[0]);
  const urlObj = new URL(req.url, `http://${req.headers.host}`);
  const sort = urlObj.searchParams.get('sort') || 'name';
  const filePath = path.join(CWD, urlPath);
  const safePath = path.normalize(filePath);

  // Prevent directory traversal
  if (!safePath.startsWith(CWD)) {
    res.writeHead(403, { 'Content-Type': 'text/plain' });
    res.end('403 Forbidden');
    return;
  }

  // API routes (virtual, not filesystem-backed)
  if (urlPath === '/__api/files') {
    return handleFileListApi(req, res);
  }

  if (urlPath === '/__api/grep') {
    return handleGrepApi(req, res);
  }

  if (urlPath === '/__search') {
    return handleSearchPage(req, res);
  }

  // Check if path exists
  if (!fs.existsSync(safePath)) {
    res.writeHead(404, { 'Content-Type': 'text/plain' });
    res.end('404 Not Found');
    return;
  }

  const stat = fs.statSync(safePath);

  if (stat.isDirectory() || stat.isFile()) {
    handleFile(safePath, urlPath, stat, res, {
      md, hljs, args, CWD, req,
      theme: args.theme || 'auto',
      sort
    });
  } else {
    res.writeHead(400, { 'Content-Type': 'text/plain' });
    res.end('400 Bad Request');
  }
}

async function handleFileListApi(req, res) {
  setSecurityHeaders(res);
  const { loadGitignore, shouldIgnoreFile } = require('./utils');
  const files = [];

  async function walk(dir, prefix) {
    const entries = await fs.promises.readdir(dir, { withFileTypes: true });
    const ig = loadGitignore(dir);
    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);
      const relPath = prefix ? prefix + '/' + entry.name : entry.name;

      if (shouldIgnoreFile(fullPath, CWD, ig, entry.isDirectory())) continue;

      if (entry.isDirectory()) {
        await walk(fullPath, relPath);
      } else {
        files.push(relPath);
      }
    }
  }

  try {
    await walk(CWD, '');
    res.writeHead(200, {
      'Content-Type': 'application/json',
      'Cache-Control': 'no-cache'
    });
    res.end(JSON.stringify(files));
  } catch (err) {
    res.writeHead(500, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: err.message }));
  }
}

async function handleGrepApi(req, res) {
  setSecurityHeaders(res);
  const url = new URL(req.url, `http://${req.headers.host}`);
  const query = url.searchParams.get('q');
  const searchPath = url.searchParams.get('path') || '.';

  if (!query) {
    res.writeHead(400, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: 'Missing q parameter' }));
    return;
  }

  const { execFile } = require('child_process');
  const { findGitRoot } = require('./gitHandler');
  const gitRoot = findGitRoot(CWD);

  const results = [];

  try {
    await new Promise((resolve, reject) => {
      if (gitRoot) {
        // Use git grep for speed
        execFile('git', ['grep', '-n', '-I', '--max-count=100', query, '--', searchPath], {
          cwd: CWD,
          timeout: 10000
        }, (err, stdout) => {
          if (err && err.code !== 1) { // code 1 = no matches
            reject(err);
            return;
          }
          if (stdout) {
            stdout.split('\n').filter(Boolean).forEach(line => {
              const match = line.match(/^([^:]+):(\d+):(.*)$/);
              if (match) {
                results.push({
                  file: match[1],
                  line: parseInt(match[2]),
                  content: match[3].substring(0, 200)
                });
              }
            });
          }
          resolve();
        });
      } else {
        // Fallback: use grep
        execFile('grep', ['-rn', '-I', '--max-count=100', query, searchPath], {
          cwd: CWD,
          timeout: 10000
        }, (err, stdout) => {
          if (err && err.code !== 1) {
            reject(err);
            return;
          }
          if (stdout) {
            stdout.split('\n').filter(Boolean).forEach(line => {
              const match = line.match(/^([^:]+):(\d+):(.*)$/);
              if (match) {
                results.push({
                  file: match[1],
                  line: parseInt(match[2]),
                  content: match[3].substring(0, 200)
                });
              }
            });
          }
          resolve();
        });
      }
    });

    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ query, results }));
  } catch (err) {
    res.writeHead(500, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: err.message }));
  }
}

async function handleSearchPage(req, res) {
  setSecurityHeaders(res);
  const { createSearchResultsTemplate } = require('./templates/search-results');
  const { escapeHtml } = require('./utils');
  const url = new URL(req.url, `http://${req.headers.host}`);
  const query = url.searchParams.get('q') || '';

  if (!query) {
    const html = createSearchResultsTemplate({
      query: '',
      results: [],
      urlPath: '/__search',
      escapeHtml
    });
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
    res.end(html);
    return;
  }

  // Reuse grep logic
  const { execFile } = require('child_process');
  const { findGitRoot } = require('./gitHandler');
  const gitRoot = findGitRoot(CWD);
  const results = [];

  try {
    await new Promise((resolve, reject) => {
      const cmd = gitRoot ? 'git' : 'grep';
      const cmdArgs = gitRoot
        ? ['grep', '-n', '-I', '--max-count=200', query, '--', '.']
        : ['-rn', '-I', '--max-count=200', query, '.'];

      execFile(cmd, cmdArgs, { cwd: CWD, timeout: 10000 }, (err, stdout) => {
        if (err && err.code !== 1) { reject(err); return; }
        if (stdout) {
          stdout.split('\n').filter(Boolean).forEach(line => {
            const match = line.match(/^([^:]+):(\d+):(.*)$/);
            if (match) {
              results.push({
                file: match[1].replace(/^\.\//, ''),
                line: parseInt(match[2]),
                content: match[3].substring(0, 200)
              });
            }
          });
        }
        resolve();
      });
    });
  } catch { /* ignore errors */ }

  const html = createSearchResultsTemplate({
    query,
    results,
    urlPath: '/__search',
    escapeHtml
  });
  res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
  res.end(html);
}

// Export for testing
module.exports = { handleRequest };