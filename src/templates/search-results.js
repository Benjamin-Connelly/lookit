const { createBaseTemplate, generateBreadcrumb } = require('./base.js');

/**
 * Create search results page
 * @param {Object} options
 * @param {string} options.query - Search query
 * @param {Array} options.results - Array of {file, line, content}
 * @param {string} options.urlPath - URL path
 * @param {Function} options.escapeHtml - Escape function
 * @returns {string} HTML document
 */
function createSearchResultsTemplate({ query, results, urlPath, escapeHtml }) {
  const breadcrumb = generateBreadcrumb(urlPath, escapeHtml);

  const resultGroups = {};
  results.forEach(r => {
    if (!resultGroups[r.file]) resultGroups[r.file] = [];
    resultGroups[r.file].push(r);
  });

  let resultsHtml = '';

  if (results.length === 0) {
    resultsHtml = '<div class="no-results">No results found for <strong>' + escapeHtml(query) + '</strong></div>';
  } else {
    resultsHtml += '<div class="results-summary">' + results.length + ' results in ' + Object.keys(resultGroups).length + ' files for <strong>' + escapeHtml(query) + '</strong></div>';

    for (const [file, matches] of Object.entries(resultGroups)) {
      resultsHtml += '<div class="result-group">';
      resultsHtml += '<div class="result-file"><a href="/' + escapeHtml(file) + '">' + escapeHtml(file) + '</a> <span class="match-count">(' + matches.length + ')</span></div>';
      resultsHtml += '<div class="result-matches">';

      for (const match of matches) {
        const highlighted = escapeHtml(match.content).replace(
          new RegExp(escapeHtml(query).replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi'),
          '<mark>$&</mark>'
        );
        resultsHtml += '<div class="result-line">';
        resultsHtml += '<a href="/' + escapeHtml(file) + '" class="line-number">' + match.line + '</a>';
        resultsHtml += '<code class="match-content">' + highlighted + '</code>';
        resultsHtml += '</div>';
      }

      resultsHtml += '</div></div>';
    }
  }

  const content = `
    <div class="search-header">
      <h2>Search Results</h2>
      <form class="grep-form" action="/__search" method="get">
        <input type="text" name="q" value="${escapeHtml(query)}" class="grep-input" placeholder="Search content..." autofocus>
        <button type="submit" class="btn">Search</button>
      </form>
    </div>
    <div class="search-results-body">
      ${resultsHtml}
    </div>
  `;

  const extraStyles = `
    .search-header {
      padding: 1.5rem;
      border-bottom: 1px solid var(--border-primary);
      background: var(--bg-tertiary);
    }
    .search-header h2 {
      font-size: 1.25rem;
      color: var(--text-primary);
      margin-bottom: 1rem;
    }
    .grep-form {
      display: flex;
      gap: 0.5rem;
    }
    .grep-input {
      flex: 1;
      padding: 0.5rem 0.75rem;
      background: var(--bg-primary);
      border: 1px solid var(--border-primary);
      border-radius: 6px;
      color: var(--text-primary);
      font-size: 0.875rem;
      font-family: var(--font-mono);
      outline: none;
    }
    .grep-input:focus {
      border-color: var(--accent-blue);
      box-shadow: 0 0 0 2px rgba(9, 105, 218, 0.2);
    }
    .search-results-body {
      padding: 1rem;
    }
    .no-results {
      padding: 3rem;
      text-align: center;
      color: var(--text-secondary);
    }
    .results-summary {
      padding: 0.75rem 0;
      color: var(--text-secondary);
      font-size: 0.875rem;
      border-bottom: 1px solid var(--border-primary);
      margin-bottom: 1rem;
    }
    .result-group {
      margin-bottom: 1.5rem;
      border: 1px solid var(--border-primary);
      border-radius: 8px;
      overflow: hidden;
    }
    .result-file {
      padding: 0.5rem 1rem;
      background: var(--bg-tertiary);
      border-bottom: 1px solid var(--border-primary);
      font-weight: 600;
      font-size: 0.875rem;
    }
    .result-file a {
      color: var(--accent-blue);
      text-decoration: none;
    }
    .result-file a:hover {
      text-decoration: underline;
    }
    .match-count {
      color: var(--text-tertiary);
      font-weight: 400;
    }
    .result-matches {
      background: var(--bg-primary);
    }
    .result-line {
      display: flex;
      align-items: baseline;
      padding: 0.25rem 1rem;
      border-bottom: 1px solid var(--border-secondary);
      font-family: var(--font-mono);
      font-size: 0.8125rem;
    }
    .result-line:last-child {
      border-bottom: none;
    }
    .line-number {
      min-width: 3em;
      text-align: right;
      padding-right: 1em;
      color: var(--text-tertiary);
      text-decoration: none;
      user-select: none;
      flex-shrink: 0;
    }
    .line-number:hover {
      color: var(--accent-blue);
    }
    .match-content {
      flex: 1;
      white-space: pre-wrap;
      word-break: break-all;
      color: var(--text-primary);
      background: none;
      border: none;
      padding: 0;
      font-size: inherit;
    }
    mark {
      background: rgba(255, 200, 0, 0.3);
      color: inherit;
      padding: 0.1em 0;
      border-radius: 2px;
    }
  `;

  return createBaseTemplate({
    title: 'Search: ' + query,
    breadcrumb,
    content,
    extraStyles
  });
}

module.exports = { createSearchResultsTemplate };
