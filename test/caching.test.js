const { describe, it, before, after } = require('node:test');
const assert = require('node:assert/strict');
const http = require('http');
const path = require('path');
const fs = require('fs');

// Reuse the helper from server.test.js
function startTestServer(dir) {
  return new Promise((resolve, reject) => {
    const { handleFile } = require('../src/fileHandler');
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
        res.writeHead(403);
        res.end();
        return;
      }

      try {
        const stat = fs.statSync(safePath);
        await handleFile(safePath, urlPath, stat, res, { md, hljs, args: {}, CWD, req });
      } catch {
        res.writeHead(404);
        res.end();
      }
    });

    server.listen(0, '127.0.0.1', () => {
      resolve({ server, port: server.address().port, url: `http://127.0.0.1:${server.address().port}` });
    });
    server.on('error', reject);
  });
}

function request(url, headers = {}) {
  return new Promise((resolve, reject) => {
    http.get(url, { headers }, (res) => {
      let body = '';
      res.on('data', chunk => body += chunk);
      res.on('end', () => resolve({ status: res.statusCode, headers: res.headers, body }));
    }).on('error', reject);
  });
}

describe('HTTP Caching', () => {
  let ctx;

  before(async () => {
    ctx = await startTestServer(path.join(__dirname, 'fixtures'));
  });

  after(() => ctx.server.close());

  it('includes ETag in response', async () => {
    const res = await request(ctx.url + '/test.md');
    assert.ok(res.headers['etag'], 'Should have ETag header');
  });

  it('returns 304 for matching ETag', async () => {
    const first = await request(ctx.url + '/test.md');
    const etag = first.headers['etag'];
    assert.ok(etag);

    const second = await request(ctx.url + '/test.md', { 'if-none-match': etag });
    assert.equal(second.status, 304);
  });

  it('includes Cache-Control header', async () => {
    const res = await request(ctx.url + '/test.md');
    assert.ok(res.headers['cache-control']);
  });
});
