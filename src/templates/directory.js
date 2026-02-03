// Directory listing template
const { createBaseTemplate, generateBreadcrumb } = require('./base.js');

/**
 * Get file icon emoji based on entry type
 * @param {Object} entry - Directory entry object
 * @param {string} fileType - File type string
 * @returns {string} Icon emoji
 */
function getFileIcon(entry, fileType) {
  // Directory icon
  if (entry.isDirectory) {
    return '📁';
  }

  // File type icons
  const iconMap = {
    'markdown': '📝',
    'code': '💻',
    'image': '🖼️',
    'video': '🎬',
    'audio': '🎵',
    'pdf': '📄',
    'binary': '📦',
    'text': '📄'
  };

  return iconMap[fileType] || '📄';
}

/**
 * Format file size in human-readable format
 * @param {number} bytes - File size in bytes
 * @returns {string} Formatted size string
 */
function formatFileSize(bytes) {
  if (bytes === 0) return '0 B';

  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const k = 1024;
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + units[i];
}

/**
 * Format date relative to now
 * @param {Date} date - Date to format
 * @returns {string} Formatted date string
 */
function formatDate(date) {
  const now = new Date();
  const diffMs = now - date;
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) {
    return 'today';
  } else if (diffDays === 1) {
    return 'yesterday';
  } else if (diffDays < 7) {
    return `${diffDays} days ago`;
  } else {
    // Format as "Jan 30, 2026"
    const options = { year: 'numeric', month: 'short', day: 'numeric' };
    return date.toLocaleDateString('en-US', options);
  }
}

/**
 * Render git status badge
 * @param {string} status - Git status code
 * @returns {string} HTML for git badge
 */
function renderGitBadge(status) {
  const badges = {
    'M': '<span class="git-badge git-orange">[M]</span>',
    'A': '<span class="git-badge git-green">[A]</span>',
    '??': '<span class="git-badge git-purple">[??]</span>',
    'D': '<span class="git-badge git-red">[D]</span>',
    'R': '<span class="git-badge git-blue">[R]</span>'
  };
  return badges[status] || '';
}

/**
 * Render git status legend
 * @returns {string} HTML for git status legend
 */
function renderGitLegend() {
  return `
    <div class="git-legend">
      <div class="git-legend-title">Git Status:</div>
      <div class="git-legend-items">
        <span class="git-legend-item"><span class="git-badge git-green">[A]</span> Added</span>
        <span class="git-legend-item"><span class="git-badge git-orange">[M]</span> Modified</span>
        <span class="git-legend-item"><span class="git-badge git-red">[D]</span> Deleted</span>
        <span class="git-legend-item"><span class="git-badge git-blue">[R]</span> Renamed</span>
        <span class="git-legend-item"><span class="git-badge git-purple">[??]</span> Untracked</span>
      </div>
    </div>
  `;
}

/**
 * Render repository statistics panel
 * @param {Object} stats - Repository statistics
 * @param {Function} escapeHtml - HTML escape function
 * @returns {string} HTML for repo stats panel
 */
function renderRepoStats(stats, escapeHtml) {
  if (!stats) return '';

  return `
    <div class="repo-stats">
      <div class="repo-stat">
        <div class="repo-stat-value">${stats.trackedFiles}</div>
        <div class="repo-stat-label">Tracked Files</div>
      </div>
      <div class="repo-stat">
        <div class="repo-stat-value">${stats.modified}</div>
        <div class="repo-stat-label">Modified</div>
      </div>
      <div class="repo-stat">
        <div class="repo-stat-value">${stats.staged}</div>
        <div class="repo-stat-label">Staged</div>
      </div>
      <div class="repo-stat">
        <div class="repo-stat-value">${stats.untracked}</div>
        <div class="repo-stat-label">Untracked</div>
      </div>
      <div class="repo-stat">
        <div class="repo-stat-value">${stats.totalCommits}</div>
        <div class="repo-stat-label">Total Commits</div>
      </div>
      <div class="repo-stat">
        <div class="repo-stat-value">${escapeHtml(stats.lastCommit)}</div>
        <div class="repo-stat-label">Last Commit</div>
      </div>
    </div>
  `;
}

/**
 * Create directory listing template
 * @param {Object} options - Template options
 * @param {string} options.dirName - Name of the directory
 * @param {Array} options.entries - Array of directory entries
 * @param {string} options.urlPath - URL path for breadcrumb
 * @param {boolean} options.showAll - Whether to show ignored files
 * @param {Function} options.escapeHtml - Function to escape HTML
 * @param {string} options.currentBranch - Current git branch name
 * @param {Object} options.repoStats - Repository statistics
 * @returns {string} Complete HTML document
 */
function createDirectoryTemplate({ dirName, entries, urlPath, showAll, escapeHtml, currentBranch, repoStats }) {
  const breadcrumb = generateBreadcrumb(urlPath, escapeHtml);

  // Sort entries: directories first, then alphabetically
  const sortedEntries = [...entries].sort((a, b) => {
    // Directories before files
    if (a.isDirectory && !b.isDirectory) return -1;
    if (!a.isDirectory && b.isDirectory) return 1;

    // Then alphabetically by name (case-insensitive)
    return a.name.toLowerCase().localeCompare(b.name.toLowerCase());
  });

  // Check if all entries have the same commit
  const commits = sortedEntries.map(e => e.lastCommit).filter(Boolean);
  const uniqueCommits = [...new Set(commits)];
  const hasCommonCommit = uniqueCommits.length === 1 && commits.length === sortedEntries.length;
  const commonCommit = hasCommonCommit ? uniqueCommits[0] : null;

  // Check if any entries have git status
  const hasGitStatus = sortedEntries.some(e => e.gitStatus);

  // Add parent directory link if not root
  const isRoot = urlPath === '/' || urlPath === '';
  let fileListHtml = '';

  if (!isRoot) {
    fileListHtml += `
      <a href=".." class="file-entry parent-dir">
        <div class="file-icon">📁</div>
        <div class="file-info">
          <div class="file-name">..</div>
        </div>
      </a>
    `;
  }

  // Generate file list HTML
  for (const entry of sortedEntries) {
    const icon = getFileIcon(entry, entry.fileType);
    const sizeStr = entry.isDirectory ? '' : formatFileSize(entry.size);
    const dateStr = formatDate(entry.mtime);

    // Gray out ignored files if they're shown
    const ignoredClass = entry.ignored ? ' ignored' : '';

    // Show commit info unless there's a common commit in header
    const showCommit = !commonCommit && entry.lastCommit;

    fileListHtml += `
      <a href="${escapeHtml(entry.url)}" class="file-entry${ignoredClass}">
        <div class="file-icon">${icon}</div>
        <div class="file-info">
          <div class="file-name">
            ${entry.gitStatus ? renderGitBadge(entry.gitStatus) : ''}${escapeHtml(entry.name)}
          </div>
          <div class="file-meta">
            <span class="file-size">${sizeStr}</span>
            ${sizeStr ? '<span class="separator">•</span>' : ''}
            <span class="file-date">${dateStr}</span>
          </div>
          ${showCommit ? `<div class="file-commit">📝 ${escapeHtml(entry.lastCommit)}</div>` : ''}
        </div>
      </a>
    `;
  }

  const itemCount = entries.length;
  const itemLabel = itemCount === 1 ? 'item' : 'items';

  const content = `
    <div class="directory-header">
      <div class="directory-icon">📁</div>
      <div class="directory-info">
        <div class="directory-name">${escapeHtml(dirName)}</div>
        <div class="directory-meta">
          ${itemCount} ${itemLabel}
          ${currentBranch ? `<span class="git-branch">🌿 ${escapeHtml(currentBranch)}</span>` : ''}
          ${commonCommit ? `<span class="separator">•</span><span class="last-commit">📝 ${escapeHtml(commonCommit)}</span>` : ''}
        </div>
      </div>
    </div>
    ${hasGitStatus ? renderGitLegend() : ''}
    ${renderRepoStats(repoStats, escapeHtml)}
    <div class="file-list">
      ${fileListHtml}
    </div>
  `;

  const extraStyles = `
    .directory-header {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1.5rem;
      background: #f8f9fa;
      border-radius: 8px 8px 0 0;
      border-bottom: 2px solid #e9ecef;
    }

    .directory-icon {
      font-size: 2rem;
      line-height: 1;
    }

    .directory-info {
      flex: 1;
    }

    .directory-name {
      font-size: 1.25rem;
      font-weight: 600;
      color: #212529;
      margin-bottom: 0.25rem;
    }

    .directory-meta {
      font-size: 0.875rem;
      color: #6c757d;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      flex-wrap: wrap;
    }

    .last-commit {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
      font-size: 0.8125rem;
      color: #495057;
    }

    .git-legend {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 0.75rem 1.5rem;
      background: #fff9e6;
      border-bottom: 1px solid #e9ecef;
      font-size: 0.8125rem;
    }

    .git-legend-title {
      font-weight: 600;
      color: #495057;
    }

    .git-legend-items {
      display: flex;
      gap: 1rem;
      flex-wrap: wrap;
    }

    .git-legend-item {
      display: inline-flex;
      align-items: center;
      gap: 0.25rem;
      color: #6c757d;
    }

    .file-list {
      background: white;
      border-radius: 0 0 8px 8px;
      overflow: hidden;
    }

    .file-entry {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem 1.5rem;
      text-decoration: none;
      color: inherit;
      border-bottom: 1px solid #e9ecef;
      transition: background-color 0.15s ease;
      position: relative;
    }

    .file-entry:last-child {
      border-bottom: none;
    }

    .file-entry:hover {
      background-color: #f8f9fa;
    }

    .file-entry.parent-dir {
      background-color: #f8f9fa;
      font-weight: 500;
    }

    .file-entry.parent-dir:hover {
      background-color: #e9ecef;
    }

    .file-entry.ignored {
      opacity: 0.5;
    }

    .file-icon {
      font-size: 1.5rem;
      line-height: 1;
      flex-shrink: 0;
    }

    .file-info {
      flex: 1;
      min-width: 0;
    }

    .file-name {
      font-size: 1rem;
      font-weight: 500;
      color: #212529;
      margin-bottom: 0.25rem;
      word-wrap: break-word;
    }

    .file-meta {
      display: flex;
      gap: 0.5rem;
      align-items: center;
      font-size: 0.875rem;
      color: #6c757d;
    }

    .file-size {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
      font-size: 0.8125rem;
    }

    .separator {
      color: #dee2e6;
    }

    .file-date {
      color: #6c757d;
    }

    .file-commit {
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
      font-size: 0.8125rem;
      color: #6c757d;
      margin-top: 0.25rem;
    }

    /* Responsive adjustments */
    @media (max-width: 768px) {
      .file-entry {
        padding: 0.875rem 1rem;
      }

      .file-icon {
        font-size: 1.25rem;
      }

      .file-name {
        font-size: 0.9375rem;
      }

      .file-meta {
        font-size: 0.8125rem;
      }

      .git-legend {
        padding: 0.5rem 1rem;
        font-size: 0.75rem;
      }

      .git-legend-items {
        gap: 0.5rem;
      }

      .last-commit {
        font-size: 0.75rem;
      }

      .file-commit {
        font-size: 0.75rem;
      }
    }
  `;

  return createBaseTemplate({
    title: dirName,
    breadcrumb,
    content,
    extraStyles
  });
}

module.exports = { createDirectoryTemplate, getFileIcon, formatFileSize, formatDate };
