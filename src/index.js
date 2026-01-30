// lookit server entry point
const http = require('http');
const https = require('https');
const fs = require('fs');
const path = require('path');
const MarkdownIt = require('markdown-it');
const hljs = require('highlight.js');
const mdHighlight = require('markdown-it-highlightjs');
const { parseArgs, printCertInstructions, openBrowser, escapeHtml } = require('./utils');
const { handleFile } = require('./fileHandler');

// Parse CLI arguments
const args = parseArgs(process.argv.slice(2));
const PORT = args.port || 7777;
const HOST = args.host || '127.0.0.1';
const CWD = process.cwd();
const CERT_DIR = path.join(process.env.HOME, '.config', 'lookit');
const CERT_PATH = args.cert || path.join(CERT_DIR, 'localhost.pem');
const KEY_PATH = args.key || path.join(CERT_DIR, 'localhost-key.pem');

// Initialize markdown renderer with syntax highlighting
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
}).use(mdHighlight, { hljs });

// Check for TLS certificates
const hasCerts = fs.existsSync(CERT_PATH) && fs.existsSync(KEY_PATH);
const useHttps = hasCerts && !args.noHttps;

if (args.httpsOnly && !hasCerts) {
  console.error('❌ Error: --https-only specified but certificates not found.\n');
  printCertInstructions();
  process.exit(1);
}

if (!hasCerts && !args.quiet) {
  console.warn('⚠️  TLS certificates not found. Falling back to HTTP.\n');
  printCertInstructions();
  console.log('');
}

// Create and start server
const server = useHttps ? createHttpsServer() : createHttpServer();

server.listen(PORT, HOST, () => {
  const protocol = useHttps ? 'https' : 'http';
  const url = `${protocol}://${HOST}:${PORT}`;

  console.log(`👀 lookit - Code Browser`);
  console.log(`📂 Serving: ${CWD}`);
  console.log(`🌐 Address:  ${url}`);
  console.log(`🔒 Security: ${useHttps ? 'HTTPS (TLS)' : 'HTTP (plaintext)'}`);
  console.log('\nPress Ctrl+C to stop.\n');

  if (args.open) {
    openBrowser(url);
  }
});

server.on('error', (err) => {
  if (err.code === 'EADDRINUSE') {
    console.error(`❌ Error: Port ${PORT} is already in use.`);
    console.error(`   Try a different port with: lookit --port ${PORT + 1}`);
  } else {
    console.error(`❌ Server error: ${err.message}`);
  }
  process.exit(1);
});

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
  const urlPath = decodeURIComponent(req.url.split('?')[0]);
  const filePath = path.join(CWD, urlPath);
  const safePath = path.normalize(filePath);

  // Prevent directory traversal
  if (!safePath.startsWith(CWD)) {
    res.writeHead(403, { 'Content-Type': 'text/plain' });
    res.end('403 Forbidden');
    return;
  }

  // Check if path exists
  if (!fs.existsSync(safePath)) {
    res.writeHead(404, { 'Content-Type': 'text/plain' });
    res.end('404 Not Found');
    return;
  }

  const stat = fs.statSync(safePath);

  if (stat.isDirectory() || stat.isFile()) {
    handleFile(safePath, urlPath, res, { md, args });
  } else {
    res.writeHead(400, { 'Content-Type': 'text/plain' });
    res.end('400 Bad Request');
  }
}

// Export for testing
module.exports = { handleRequest };