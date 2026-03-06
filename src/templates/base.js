// Base HTML template
const { baseStyles } = require('../styles.js');
const { getSearchOverlayHtml } = require('./search.js');

/**
 * Generate breadcrumb navigation HTML
 * @param {string} urlPath - The current URL path
 * @param {Function} escapeHtml - Function to escape HTML special characters
 * @returns {string} Breadcrumb HTML
 */
function generateBreadcrumb(urlPath, escapeHtml) {
  // Split path into parts and filter out empty strings
  const parts = urlPath.split('/').filter(Boolean);

  // Start with home link
  let breadcrumbHtml = '<a href="/">🏠</a>';

  // Build path incrementally
  let currentPath = '';

  for (let i = 0; i < parts.length; i++) {
    const part = parts[i];
    currentPath += '/' + part;

    // Add separator
    breadcrumbHtml += ' <span class="separator">/</span> ';

    // If this is the last part, show it as current (non-clickable)
    if (i === parts.length - 1) {
      breadcrumbHtml += `<span class="current">${escapeHtml(decodeURIComponent(part))}</span>`;
    } else {
      // Otherwise, make it a clickable link
      breadcrumbHtml += `<a href="${escapeHtml(currentPath)}">${escapeHtml(decodeURIComponent(part))}</a>`;
    }
  }

  return breadcrumbHtml;
}

/**
 * Create base HTML template
 * @param {Object} options - Template options
 * @param {string} options.title - Page title
 * @param {string} options.breadcrumb - Breadcrumb HTML
 * @param {string} options.content - Main content HTML
 * @param {string} [options.extraStyles=''] - Additional CSS styles
 * @param {string} [options.extraHead=''] - Additional head content
 * @returns {string} Complete HTML document
 */
function getHljsLinks(theme) {
  const base = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles';
  if (theme === 'dark') {
    return `<link rel="stylesheet" href="${base}/github-dark.min.css">`;
  }
  if (theme === 'light') {
    return `<link rel="stylesheet" href="${base}/github.min.css">`;
  }
  // auto: both with media queries
  return `<link rel="stylesheet" href="${base}/github-dark.min.css" media="(prefers-color-scheme: dark)">
  <link rel="stylesheet" href="${base}/github.min.css" media="(prefers-color-scheme: light)">`;
}

function createBaseTemplate({ title, breadcrumb, content, extraStyles = '', extraHead = '', theme = 'auto' }) {
  const themeAttr = (theme === 'light' || theme === 'dark') ? ` data-theme="${theme}"` : '';
  const hljsLinks = getHljsLinks(theme);
  return `<!DOCTYPE html>
<html lang="en"${themeAttr}>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${title} - lookit</title>
  ${hljsLinks}
  <style>${baseStyles}
  .theme-toggle {
    margin-left: auto;
    background: none;
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-md);
    padding: var(--space-1) var(--space-2);
    cursor: pointer;
    color: var(--text-secondary);
    font-size: var(--text-sm);
    transition: all var(--transition-fast);
    display: inline-flex;
    align-items: center;
    gap: var(--space-1);
  }
  .theme-toggle:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }
  ${extraStyles}</style>
  ${extraHead}
</head>
<body>
  <div class="container">
    <div class="header">
      <div class="breadcrumb">${breadcrumb}<button class="theme-toggle" onclick="cycleTheme()" aria-label="Toggle theme" title="Toggle theme"><span class="theme-icon"></span></button></div>
    </div>
    <div class="content">
      ${content}
    </div>
  </div>
  ${getSearchOverlayHtml()}
  <script>
  (function() {
    var stored = localStorage.getItem('lookit-theme');
    var html = document.documentElement;

    function getSystemTheme() {
      return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }

    function applyTheme(mode) {
      if (mode === 'auto' || !mode) {
        html.removeAttribute('data-theme');
      } else {
        html.setAttribute('data-theme', mode);
      }
      updateHljsLinks(mode || 'auto');
      updateIcon(mode || 'auto');
    }

    function updateIcon(mode) {
      var icon = document.querySelector('.theme-icon');
      if (!icon) return;
      var icons = { light: '\\u2600\\uFE0F', dark: '\\uD83C\\uDF19', auto: '\\uD83D\\uDCBB' };
      icon.textContent = icons[mode] || icons.auto;
    }

    function updateHljsLinks(mode) {
      var base = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles';
      var head = document.head;
      var old = head.querySelectorAll('link[href*="highlight.js"][rel="stylesheet"]');
      old.forEach(function(l) { l.remove(); });

      if (mode === 'dark') {
        addLink(base + '/github-dark.min.css');
      } else if (mode === 'light') {
        addLink(base + '/github.min.css');
      } else {
        addLink(base + '/github-dark.min.css', '(prefers-color-scheme: dark)');
        addLink(base + '/github.min.css', '(prefers-color-scheme: light)');
      }
    }

    function addLink(href, media) {
      var link = document.createElement('link');
      link.rel = 'stylesheet';
      link.href = href;
      if (media) link.media = media;
      document.head.appendChild(link);
    }

    window.cycleTheme = function() {
      var current = localStorage.getItem('lookit-theme') || 'auto';
      var next = current === 'auto' ? 'light' : current === 'light' ? 'dark' : 'auto';
      if (next === 'auto') {
        localStorage.removeItem('lookit-theme');
      } else {
        localStorage.setItem('lookit-theme', next);
      }
      applyTheme(next);
    };

    window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', function() {
      if (!localStorage.getItem('lookit-theme')) {
        applyTheme('auto');
      }
    });

    applyTheme(stored || 'auto');
  })();
  </script>
</body>
</html>`;
}

module.exports = { createBaseTemplate, generateBreadcrumb };
