// Markdown template
const { createBaseTemplate, generateBreadcrumb } = require('./base.js');

/**
 * Create markdown file template with rendered content
 * @param {Object} options - Template options
 * @param {string} options.fileName - Name of the file
 * @param {string} options.html - Rendered markdown HTML
 * @param {string} options.urlPath - URL path for breadcrumb
 * @param {Function} options.escapeHtml - Function to escape HTML
 * @returns {string} Complete HTML document
 */
function createMarkdownTemplate({ fileName, html, urlPath, escapeHtml }) {
  const breadcrumb = generateBreadcrumb(urlPath, escapeHtml);

  const liveReloadScript = `
<script>
(function() {
  const currentPath = ${JSON.stringify(urlPath)};
  function connect() {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(protocol + '//' + location.host + '/__ws');
    ws.onmessage = function(event) {
      const changedPath = event.data;
      if (currentPath === changedPath || currentPath.endsWith(changedPath)) {
        location.reload();
      }
    };
    ws.onclose = function() {
      setTimeout(connect, 3000);
    };
    ws.onerror = function() {
      ws.close();
    };
  }
  connect();
})();
</script>
`;

  const content = `
    <div class="file-header">
      <div class="file-icon">📝</div>
      <div class="file-info">
        <div class="file-name">${escapeHtml(fileName)}</div>
        <div class="file-meta">
          <span class="format-badge">Markdown</span>
        </div>
      </div>
    </div>
    <article>
      <div class="markdown-body" aria-label="Document content">
        ${html}
      </div>
    </article>
    ${liveReloadScript}
  `;

  const extraStyles = `
    .file-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1.5rem;
      background: var(--bg-tertiary);
      border-radius: 8px 8px 0 0;
      border-bottom: 2px solid var(--border-primary);
    }

    .file-icon {
      font-size: 2rem;
      line-height: 1;
    }

    .file-info {
      flex: 1;
    }

    .file-name {
      font-size: 1.25rem;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 0.5rem;
    }

    .file-meta {
      display: flex;
      gap: 0.5rem;
      align-items: center;
    }

    .format-badge {
      display: inline-block;
      padding: 0.25rem 0.75rem;
      border-radius: 12px;
      font-size: 0.75rem;
      font-weight: 600;
      color: var(--bg-primary);
      background-color: var(--accent-blue);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    /* GitHub-style markdown rendering */
    .markdown-body {
      padding: 2rem;
      background: var(--bg-primary);
      border-radius: 0 0 8px 8px;
      color: var(--text-primary);
      font-size: 16px;
      line-height: 1.6;
    }

    /* Headings */
    .markdown-body h1,
    .markdown-body h2,
    .markdown-body h3,
    .markdown-body h4,
    .markdown-body h5,
    .markdown-body h6 {
      margin-top: 24px;
      margin-bottom: 16px;
      font-weight: 600;
      line-height: 1.25;
      color: var(--text-primary);
    }

    .markdown-body h1 {
      font-size: 2em;
      padding-bottom: 0.3em;
      border-bottom: 1px solid var(--border-primary);
    }

    .markdown-body h2 {
      font-size: 1.5em;
      padding-bottom: 0.3em;
      border-bottom: 1px solid var(--border-primary);
    }

    .markdown-body h3 {
      font-size: 1.25em;
    }

    .markdown-body h4 {
      font-size: 1em;
    }

    .markdown-body h5 {
      font-size: 0.875em;
    }

    .markdown-body h6 {
      font-size: 0.85em;
      color: var(--text-secondary);
    }

    /* Paragraphs */
    .markdown-body p {
      margin-top: 0;
      margin-bottom: 16px;
    }

    /* Links */
    .markdown-body a {
      color: var(--accent-blue);
      text-decoration: none;
    }

    .markdown-body a:hover {
      text-decoration: underline;
    }

    /* Lists */
    .markdown-body ul,
    .markdown-body ol {
      margin-top: 0;
      margin-bottom: 16px;
      padding-left: 2em;
    }

    .markdown-body ul ul,
    .markdown-body ul ol,
    .markdown-body ol ol,
    .markdown-body ol ul {
      margin-top: 0;
      margin-bottom: 0;
    }

    .markdown-body li {
      margin-top: 0.25em;
    }

    .markdown-body li > p {
      margin-top: 16px;
    }

    .markdown-body li + li {
      margin-top: 0.25em;
    }

    /* Code blocks */
    .markdown-body pre {
      margin-top: 0;
      margin-bottom: 16px;
      padding: 16px;
      overflow: auto;
      font-size: 85%;
      line-height: 1.45;
      background-color: var(--bg-primary);
      border-radius: 6px;
    }

    .markdown-body pre code {
      display: inline;
      padding: 0;
      margin: 0;
      overflow: visible;
      line-height: inherit;
      word-wrap: normal;
      background-color: transparent;
      border: 0;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
      color: var(--text-primary);
    }

    /* Inline code */
    .markdown-body code {
      padding: 0.2em 0.4em;
      margin: 0;
      font-size: 85%;
      color: inherit;
      background-color: rgba(175, 184, 193, 0.2);
      border: none;
      border-radius: 6px;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
    }

    /* Tables */
    .markdown-body table {
      border-spacing: 0;
      border-collapse: collapse;
      margin-top: 0;
      margin-bottom: 16px;
      width: 100%;
      overflow: auto;
    }

    .markdown-body table th {
      font-weight: 600;
      padding: 6px 13px;
      border: 1px solid var(--border-primary);
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
    }

    .markdown-body table td {
      padding: 6px 13px;
      border: 1px solid var(--border-primary);
    }

    .markdown-body table tr {
      background-color: var(--bg-primary);
      border-top: 1px solid var(--border-primary);
    }

    .markdown-body table tr:nth-child(2n) {
      background-color: var(--bg-tertiary);
    }

    /* Blockquotes */
    .markdown-body blockquote {
      margin: 0;
      margin-bottom: 16px;
      padding: 0 1em;
      color: var(--text-secondary);
      border-left: 0.25em solid var(--border-primary);
    }

    .markdown-body blockquote > :first-child {
      margin-top: 0;
    }

    .markdown-body blockquote > :last-child {
      margin-bottom: 0;
    }

    /* Horizontal rules */
    .markdown-body hr {
      height: 0.25em;
      padding: 0;
      margin: 24px 0;
      background-color: var(--border-primary);
      border: 0;
    }

    /* Images */
    .markdown-body img {
      max-width: 100%;
      height: auto;
      border-radius: 6px;
      margin: 16px 0;
    }

    /* Task lists */
    .markdown-body input[type="checkbox"] {
      margin: 0 0.2em 0.25em -1.6em;
      vertical-align: middle;
    }

    /* Strong and emphasis */
    .markdown-body strong {
      font-weight: 600;
    }

    .markdown-body em {
      font-style: italic;
    }

    /* Override highlight.js styles for code blocks */
    .markdown-body .hljs {
      background: transparent !important;
      padding: 0 !important;
    }

    /* Ensure first child has no top margin */
    .markdown-body > *:first-child {
      margin-top: 0 !important;
    }

    /* Ensure last child has no bottom margin */
    .markdown-body > *:last-child {
      margin-bottom: 0 !important;
    }
  `;

  return createBaseTemplate({
    title: fileName,
    breadcrumb,
    content,
    extraStyles
  });
}

module.exports = { createMarkdownTemplate };
