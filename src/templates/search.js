/**
 * Get search overlay HTML/CSS/JS
 * Returns HTML string to inject into pages for Ctrl+K file search
 */
function getSearchOverlayHtml() {
  return `
<div id="search-overlay" class="search-overlay" style="display:none">
  <div class="search-backdrop" onclick="closeSearch()"></div>
  <div class="search-modal">
    <div class="search-input-wrapper">
      <svg class="search-icon" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M11.5 7a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Zm-.82 4.74a6 6 0 1 1 1.06-1.06l3.04 3.04a.75.75 0 1 1-1.06 1.06l-3.04-3.04Z"/>
      </svg>
      <input id="search-input" type="text" class="search-field" placeholder="Search files..." autocomplete="off" spellcheck="false">
      <kbd class="search-kbd">Esc</kbd>
    </div>
    <div id="search-results" class="search-results"></div>
    <div class="search-footer">
      <span><kbd>&#x2191;&#x2193;</kbd> navigate</span>
      <span><kbd>&#x21B5;</kbd> open</span>
      <span><kbd>esc</kbd> close</span>
    </div>
  </div>
</div>

<style>
.search-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 20vh;
}
.search-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
}
.search-modal {
  position: relative;
  width: 560px;
  max-width: 90vw;
  max-height: 60vh;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.search-input-wrapper {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-primary);
}
.search-icon {
  color: var(--text-tertiary);
  flex-shrink: 0;
}
.search-field {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  font-size: 1rem;
  color: var(--text-primary);
  font-family: var(--font-sans);
}
.search-field::placeholder {
  color: var(--text-tertiary);
}
.search-kbd {
  padding: 0.125rem 0.375rem;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  font-size: 0.7rem;
  color: var(--text-tertiary);
  font-family: var(--font-mono);
}
.search-results {
  overflow-y: auto;
  max-height: 50vh;
}
.search-result {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  cursor: pointer;
  color: var(--text-primary);
  text-decoration: none;
  font-size: 0.875rem;
}
.search-result:hover,
.search-result.active {
  background: var(--bg-hover);
}
.search-result-icon {
  flex-shrink: 0;
  font-size: 1rem;
}
.search-result-path {
  color: var(--text-tertiary);
  font-size: 0.75rem;
  margin-left: 0.25rem;
}
.search-empty {
  padding: 2rem;
  text-align: center;
  color: var(--text-tertiary);
  font-size: 0.875rem;
}
.search-footer {
  display: flex;
  gap: 1rem;
  padding: 0.5rem 1rem;
  border-top: 1px solid var(--border-primary);
  font-size: 0.7rem;
  color: var(--text-tertiary);
}
.search-footer kbd {
  padding: 0.1rem 0.3rem;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: 3px;
  font-family: var(--font-mono);
  font-size: 0.65rem;
}
</style>

<script>
(function() {
  var files = null;
  var activeIndex = -1;

  window.openSearch = function() {
    var overlay = document.getElementById('search-overlay');
    overlay.style.display = 'flex';
    var input = document.getElementById('search-input');
    input.value = '';
    input.focus();
    activeIndex = -1;

    if (!files) {
      fetch('/__api/files')
        .then(function(r) { return r.json(); })
        .then(function(data) { files = data; })
        .catch(function() { files = []; });
    }

    renderResults([]);
  };

  window.closeSearch = function() {
    document.getElementById('search-overlay').style.display = 'none';
  };

  function fuzzyMatch(query, str) {
    query = query.toLowerCase();
    str = str.toLowerCase();
    if (str.includes(query)) return { match: true, score: str.indexOf(query) };
    var qi = 0;
    for (var si = 0; si < str.length && qi < query.length; si++) {
      if (str[si] === query[qi]) qi++;
    }
    if (qi === query.length) return { match: true, score: 1000 };
    return { match: false, score: Infinity };
  }

  function getIcon(path) {
    var ext = path.split('.').pop().toLowerCase();
    var icons = {
      js: '\\uD83D\\uDCBB', ts: '\\uD83D\\uDCBB', py: '\\uD83D\\uDCBB', go: '\\uD83D\\uDCBB', rs: '\\uD83D\\uDCBB', rb: '\\uD83D\\uDCBB',
      md: '\\uD83D\\uDCDD', mdx: '\\uD83D\\uDCDD',
      json: '\\uD83D\\uDCCB', yaml: '\\uD83D\\uDCCB', yml: '\\uD83D\\uDCCB', toml: '\\uD83D\\uDCCB',
      html: '\\uD83C\\uDF10', css: '\\uD83C\\uDFA8',
      png: '\\uD83D\\uDDBC\\uFE0F', jpg: '\\uD83D\\uDDBC\\uFE0F', gif: '\\uD83D\\uDDBC\\uFE0F', svg: '\\uD83D\\uDDBC\\uFE0F',
      sh: '\\u26A1', bash: '\\u26A1'
    };
    return icons[ext] || '\\uD83D\\uDCC4';
  }

  function renderResults(results) {
    var container = document.getElementById('search-results');
    if (results.length === 0) {
      var input = document.getElementById('search-input');
      container.innerHTML = input.value
        ? '<div class="search-empty">No files found</div>'
        : '<div class="search-empty">Type to search files...</div>';
      return;
    }
    container.innerHTML = results.slice(0, 50).map(function(r, i) {
      var parts = r.split('/');
      var name = parts.pop();
      var dir = parts.join('/');
      return '<a class="search-result' + (i === activeIndex ? ' active' : '') + '" href="/' + r + '">' +
        '<span class="search-result-icon">' + getIcon(r) + '</span>' +
        '<span>' + name + (dir ? '<span class="search-result-path">' + dir + '</span>' : '') + '</span>' +
        '</a>';
    }).join('');
  }

  document.addEventListener('keydown', function(e) {
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault();
      openSearch();
      return;
    }

    var overlay = document.getElementById('search-overlay');
    if (overlay.style.display === 'none') return;

    if (e.key === 'Escape') {
      closeSearch();
      return;
    }

    var results = document.querySelectorAll('.search-result');

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      activeIndex = Math.min(activeIndex + 1, results.length - 1);
      results.forEach(function(r, i) { r.classList.toggle('active', i === activeIndex); });
      if (results[activeIndex]) results[activeIndex].scrollIntoView({ block: 'nearest' });
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      activeIndex = Math.max(activeIndex - 1, 0);
      results.forEach(function(r, i) { r.classList.toggle('active', i === activeIndex); });
      if (results[activeIndex]) results[activeIndex].scrollIntoView({ block: 'nearest' });
    } else if (e.key === 'Enter' && activeIndex >= 0 && results[activeIndex]) {
      results[activeIndex].click();
    }
  });

  document.addEventListener('input', function(e) {
    if (e.target.id !== 'search-input') return;
    var query = e.target.value.trim();
    activeIndex = -1;

    if (!query || !files) {
      renderResults([]);
      return;
    }

    var scored = files
      .map(function(f) { var m = fuzzyMatch(query, f); return { path: f, match: m.match, score: m.score }; })
      .filter(function(f) { return f.match; })
      .sort(function(a, b) { return a.score - b.score; });

    renderResults(scored.map(function(s) { return s.path; }));
  });
})();
</script>
`;
}

module.exports = { getSearchOverlayHtml };
