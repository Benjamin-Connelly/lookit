// Binary file template
const { createBaseTemplate, generateBreadcrumb } = require('./base.js');

/**
 * Format file size in human-readable format
 * @param {number} bytes - File size in bytes
 * @returns {string} Formatted file size (e.g., "1.2 MB", "456 KB")
 */
function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes';

  const units = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const k = 1024;
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  // For bytes, don't show decimal
  if (i === 0) {
    return bytes + ' ' + units[i];
  }

  // For larger units, show 1 decimal place
  return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + units[i];
}

/**
 * Get file icon emoji based on file extension
 * @param {string} ext - File extension (with or without leading dot)
 * @returns {string} Emoji icon for the file type
 */
function getFileIcon(ext) {
  // Normalize extension (remove leading dot, lowercase)
  const normalized = ext.toLowerCase().replace(/^\./, '');

  const iconMap = {
    // Archive files
    'zip': '📦',
    'tar': '📦',
    'gz': '📦',
    'bz2': '📦',
    'xz': '📦',
    '7z': '📦',
    'rar': '📦',

    // Executable files
    'exe': '⚙️',
    'dll': '⚙️',
    'so': '⚙️',
    'dylib': '⚙️',
    'app': '⚙️',
    'bin': '⚙️',

    // Database files
    'db': '🗄️',
    'sqlite': '🗄️',
    'sqlite3': '🗄️',

    // Office documents
    'doc': '📄',
    'docx': '📄',
    'xls': '📊',
    'xlsx': '📊',
    'ppt': '📽️',
    'pptx': '📽️',

    // Other common binary formats
    'iso': '💿',
    'dmg': '💿',
    'pkg': '📦',
    'deb': '📦',
    'rpm': '📦',
    'apk': '📦',
    'jar': '📦'
  };

  return iconMap[normalized] || '🔒';
}

/**
 * Create binary file template with preview card
 * @param {Object} options - Template options
 * @param {string} options.fileName - Name of the file
 * @param {string} options.filePath - Absolute path to the file
 * @param {number} options.fileSize - File size in bytes
 * @param {Date} options.modified - Last modified date
 * @param {string} options.urlPath - URL path for breadcrumb
 * @param {Function} options.escapeHtml - Function to escape HTML
 * @returns {string} Complete HTML document
 */
function createBinaryTemplate({ fileName, filePath, fileSize, modified, urlPath, escapeHtml }) {
  const breadcrumb = generateBreadcrumb(urlPath, escapeHtml);

  // Get file extension and icon
  const ext = fileName.includes('.') ? fileName.split('.').pop() : '';
  const icon = getFileIcon(ext);

  // Format file size
  const formattedSize = formatFileSize(fileSize);

  // Format last modified date
  const formattedDate = modified.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });

  // File type display
  const fileType = ext ? ext.toUpperCase() : 'Unknown';

  const content = `
    <div class="binary-card">
      <div class="binary-icon">${icon}</div>
      <div class="binary-filename">${escapeHtml(fileName)}</div>
      <div class="binary-message">Binary file - cannot preview</div>

      <div class="binary-info">
        <div class="info-row">
          <span class="info-label">File Size:</span>
          <span class="info-value">${escapeHtml(formattedSize)}</span>
        </div>
        <div class="info-row">
          <span class="info-label">Last Modified:</span>
          <span class="info-value">${escapeHtml(formattedDate)}</span>
        </div>
        <div class="info-row">
          <span class="info-label">Type:</span>
          <span class="info-value">${escapeHtml(fileType)}</span>
        </div>
      </div>

      <div class="binary-actions">
        <button class="btn btn-secondary" onclick="copyPath()">Copy Path</button>
        <a href="${escapeHtml(urlPath)}?download=true" class="btn btn-primary">Download</a>
      </div>
    </div>

    <script>
      function copyPath() {
        const path = ${JSON.stringify(filePath)};
        navigator.clipboard.writeText(path).then(() => {
          const btn = event.target;
          const originalText = btn.textContent;
          btn.textContent = 'Copied!';
          btn.style.backgroundColor = '#28a745';
          setTimeout(() => {
            btn.textContent = originalText;
            btn.style.backgroundColor = '';
          }, 2000);
        }).catch(err => {
          alert('Failed to copy path: ' + err.message);
        });
      }
    </script>
  `;

  const extraStyles = `
    .binary-card {
      max-width: 600px;
      margin: 3rem auto;
      padding: 3rem 2rem;
      background: white;
      border: 2px solid #e9ecef;
      border-radius: 12px;
      text-align: center;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
    }

    .binary-icon {
      font-size: 5rem;
      line-height: 1;
      margin-bottom: 1.5rem;
    }

    .binary-filename {
      font-size: 1.5rem;
      font-weight: 600;
      color: #212529;
      margin-bottom: 0.5rem;
      word-break: break-all;
    }

    .binary-message {
      color: #6c757d;
      font-size: 1rem;
      margin-bottom: 2rem;
    }

    .binary-info {
      background: #f8f9fa;
      border-radius: 8px;
      padding: 1.5rem;
      margin-bottom: 2rem;
      text-align: left;
    }

    .info-row {
      display: flex;
      justify-content: space-between;
      padding: 0.75rem 0;
      border-bottom: 1px solid #e9ecef;
    }

    .info-row:last-child {
      border-bottom: none;
    }

    .info-label {
      font-weight: 600;
      color: #495057;
    }

    .info-value {
      color: #212529;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
      font-size: 0.9rem;
    }

    .binary-actions {
      display: flex;
      gap: 1rem;
      justify-content: center;
    }

    .btn {
      padding: 0.75rem 1.5rem;
      border: none;
      border-radius: 6px;
      font-size: 1rem;
      font-weight: 600;
      cursor: pointer;
      text-decoration: none;
      display: inline-block;
      transition: all 0.2s;
    }

    .btn-primary {
      background: #007bff;
      color: white;
    }

    .btn-primary:hover {
      background: #0056b3;
      transform: translateY(-1px);
      box-shadow: 0 2px 4px rgba(0, 123, 255, 0.3);
    }

    .btn-secondary {
      background: #6c757d;
      color: white;
    }

    .btn-secondary:hover {
      background: #5a6268;
      transform: translateY(-1px);
      box-shadow: 0 2px 4px rgba(108, 117, 125, 0.3);
    }
  `;

  return createBaseTemplate({
    title: fileName,
    breadcrumb,
    content,
    extraStyles
  });
}

module.exports = {
  formatFileSize,
  getFileIcon,
  createBinaryTemplate
};
