const { describe, it, before, after } = require('node:test');
const assert = require('node:assert/strict');
const http = require('http');
const path = require('path');
const fs = require('fs');

// Helper to start server on a random port
function startTestServer(dir) {
  return new Promise((resolve, reject) => {
    const { handleFile, detectFileType } = require('../src/fileHandler');
    const MarkdownIt = require('markdown-it');
    const hljs = require('highlight.js');
    const mdHighlight = require('markdown-it-highlightjs');

    const md = new MarkdownIt({ html: true, linkify: true, typographer: true })
      .use(mdHighlight, { hljs });

    const CWD = dir || path.join(__dirname, 'fixtures');

    const server = http.createServer(async (req, res) => {
      const urlPath = decodeURIComponent(req.url.split('?')[0]);
      const filePath = path.join(CWD, urlPath);
      const safePath = path.normalize(filePath);

      if (!safePath.startsWith(CWD)) {
        res.writeHead(403, { 'Content-Type': 'text/plain' });
        res.end('403 Forbidden');
        return;
      }

      try {
        const stat = fs.statSync(safePath);
        await handleFile(safePath, urlPath, stat, res, { md, hljs, args: {}, CWD, req });
      } catch (err) {
        if (err.code === 'ENOENT') {
          res.writeHead(404, { 'Content-Type': 'text/plain' });
          res.end('404 Not Found');
        } else {
          res.writeHead(500, { 'Content-Type': 'text/plain' });
          res.end('500 Internal Server Error');
        }
      }
    });

    server.listen(0, '127.0.0.1', () => {
      const port = server.address().port;
      resolve({ server, port, url: `http://127.0.0.1:${port}` });
    });

    server.on('error', reject);
  });
}

// Helper to make HTTP requests
function request(url, options = {}) {
  return new Promise((resolve, reject) => {
    const req = http.get(url, { ...options }, (res) => {
      let body = '';
      res.on('data', chunk => body += chunk);
      res.on('end', () => resolve({ status: res.statusCode, headers: res.headers, body }));
    });
    req.on('error', reject);
  });
}

// Ensure test fixtures exist
const fixturesDir = path.join(__dirname, 'fixtures');
if (!fs.existsSync(fixturesDir)) {
  fs.mkdirSync(fixturesDir, { recursive: true });
}

// Create test fixture files if needed
const testMd = path.join(fixturesDir, 'test.md');
if (!fs.existsSync(testMd)) {
  fs.writeFileSync(testMd, '# Test\n\nHello **world**\n');
}

const testJs = path.join(fixturesDir, 'test.js');
if (!fs.existsSync(testJs)) {
  fs.writeFileSync(testJs, 'console.log("hello");\n');
}

describe('Server', () => {
  let ctx;

  before(async () => {
    ctx = await startTestServer(fixturesDir);
  });

  after(() => {
    ctx.server.close();
  });

  it('serves directory listing', async () => {
    const res = await request(ctx.url + '/');
    assert.equal(res.status, 200);
    assert.ok(res.headers['content-type'].includes('text/html'));
    assert.ok(res.body.includes('test.md'));
  });

  it('serves markdown files as HTML', async () => {
    const res = await request(ctx.url + '/test.md');
    assert.equal(res.status, 200);
    assert.ok(res.headers['content-type'].includes('text/html'));
    assert.ok(res.body.includes('<strong>') || res.body.includes('<h1>'), 'Should render markdown to HTML');
  });

  it('serves code files with syntax highlighting', async () => {
    const res = await request(ctx.url + '/test.js');
    assert.equal(res.status, 200);
    assert.ok(res.headers['content-type'].includes('text/html'));
    assert.ok(res.body.includes('hljs'));
  });

  it('blocks path traversal attempts', async () => {
    const res = await request(ctx.url + '/../../../etc/passwd');
    assert.ok([403, 404].includes(res.status), 'Should not serve files outside CWD');
  });

  it('returns 404 for missing files', async () => {
    const res = await request(ctx.url + '/nonexistent.txt');
    assert.equal(res.status, 404);
  });

  it('includes theme toggle in HTML', async () => {
    const res = await request(ctx.url + '/');
    assert.equal(res.status, 200);
    assert.ok(res.body.includes('theme-toggle') || res.body.includes('cycleTheme'));
  });
});
