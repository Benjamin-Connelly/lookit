// Modern CSS styles inspired by GitHub/Vercel/Linear

const baseStyles = `
  /* CSS Custom Properties */
  :root {
    /* Dark Theme Colors */
    --bg-primary: #0d1117;
    --bg-secondary: #161b22;
    --bg-tertiary: #21262d;
    --bg-hover: #30363d;

    --border-primary: #30363d;
    --border-secondary: #21262d;

    --text-primary: #e6edf3;
    --text-secondary: #8b949e;
    --text-tertiary: #6e7681;

    --accent-blue: #58a6ff;
    --accent-blue-hover: #79b8ff;
    --accent-green: #3fb950;
    --accent-purple: #a371f7;
    --accent-orange: #ff8c00;

    /* Spacing Scale */
    --space-1: 0.25rem;   /* 4px */
    --space-2: 0.5rem;    /* 8px */
    --space-3: 0.75rem;   /* 12px */
    --space-4: 1rem;      /* 16px */
    --space-5: 1.25rem;   /* 20px */
    --space-6: 1.5rem;    /* 24px */
    --space-7: 2rem;      /* 32px */
    --space-8: 2.5rem;    /* 40px */
    --space-9: 3rem;      /* 48px */
    --space-10: 4rem;     /* 64px */
    --space-11: 5rem;     /* 80px */
    --space-12: 6rem;     /* 96px */

    /* Typography */
    --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Noto Sans', Helvetica, Arial, sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji';
    --font-mono: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, 'Liberation Mono', monospace;

    /* Font Sizes */
    --text-xs: 0.75rem;    /* 12px */
    --text-sm: 0.875rem;   /* 14px */
    --text-base: 1rem;     /* 16px */
    --text-lg: 1.125rem;   /* 18px */
    --text-xl: 1.25rem;    /* 20px */
    --text-2xl: 1.5rem;    /* 24px */
    --text-3xl: 1.875rem;  /* 30px */

    /* Shadows */
    --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.3);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.4), 0 2px 4px -1px rgba(0, 0, 0, 0.3);
    --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.5), 0 4px 6px -2px rgba(0, 0, 0, 0.4);
    --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.6), 0 10px 10px -5px rgba(0, 0, 0, 0.4);

    /* Border Radius */
    --radius-sm: 0.25rem;  /* 4px */
    --radius-md: 0.375rem; /* 6px */
    --radius-lg: 0.5rem;   /* 8px */
    --radius-xl: 0.75rem;  /* 12px */

    /* Transitions */
    --transition-fast: 150ms cubic-bezier(0.4, 0, 0.2, 1);
    --transition-base: 200ms cubic-bezier(0.4, 0, 0.2, 1);
    --transition-slow: 300ms cubic-bezier(0.4, 0, 0.2, 1);
  }

  /* Base Reset */
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  /* Body Styles */
  body {
    font-family: var(--font-sans);
    background-color: var(--bg-primary);
    color: var(--text-primary);
    line-height: 1.6;
    font-size: var(--text-base);
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    padding: var(--space-6);
  }

  /* Container */
  .container {
    max-width: 1200px;
    margin: 0 auto;
    background-color: var(--bg-secondary);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-lg);
    border: 1px solid var(--border-primary);
    overflow: hidden;
  }

  /* Header */
  header {
    padding: var(--space-6);
    border-bottom: 1px solid var(--border-primary);
    background-color: var(--bg-tertiary);
  }

  h1 {
    font-size: var(--text-3xl);
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--space-2);
    letter-spacing: -0.025em;
  }

  /* Breadcrumb Navigation */
  .breadcrumb {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: var(--space-2);
    padding: var(--space-4);
    background-color: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-primary);
    font-size: var(--text-sm);
  }

  .breadcrumb a {
    color: var(--accent-blue);
    text-decoration: none;
    transition: color var(--transition-fast);
    display: inline-flex;
    align-items: center;
    gap: var(--space-1);
  }

  .breadcrumb a:hover {
    color: var(--accent-blue-hover);
    text-decoration: underline;
  }

  .breadcrumb .separator {
    color: var(--text-tertiary);
    user-select: none;
  }

  .breadcrumb .current {
    color: var(--text-secondary);
    font-weight: 500;
  }

  /* Content Area */
  .content {
    padding: var(--space-6);
  }

  /* Directory Listing */
  .directory-list {
    list-style: none;
    display: grid;
    gap: var(--space-3);
  }

  .directory-item {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-4);
    background-color: var(--bg-tertiary);
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-md);
    transition: all var(--transition-base);
  }

  .directory-item:hover {
    background-color: var(--bg-hover);
    border-color: var(--border-primary);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }

  .directory-item a {
    color: var(--accent-blue);
    text-decoration: none;
    font-weight: 500;
    flex: 1;
    transition: color var(--transition-fast);
  }

  .directory-item a:hover {
    color: var(--accent-blue-hover);
  }

  /* File Icon */
  .file-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    flex-shrink: 0;
  }

  .file-icon.folder {
    color: var(--accent-blue);
  }

  .file-icon.file {
    color: var(--text-secondary);
  }

  /* Code Container */
  .code-container {
    background-color: var(--bg-primary);
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-md);
    overflow: hidden;
    margin: var(--space-4) 0;
  }

  .code-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--space-3) var(--space-4);
    background-color: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-primary);
  }

  .code-filename {
    font-family: var(--font-mono);
    font-size: var(--text-sm);
    color: var(--text-primary);
    font-weight: 500;
  }

  .language-badge {
    display: inline-flex;
    align-items: center;
    padding: var(--space-1) var(--space-3);
    background-color: var(--bg-hover);
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-sm);
    font-size: var(--text-xs);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-secondary);
  }

  .language-badge.javascript,
  .language-badge.js {
    background-color: rgba(247, 223, 30, 0.1);
    color: #f7df1e;
    border-color: rgba(247, 223, 30, 0.2);
  }

  .language-badge.python,
  .language-badge.py {
    background-color: rgba(53, 114, 165, 0.1);
    color: #3572a5;
    border-color: rgba(53, 114, 165, 0.2);
  }

  .language-badge.typescript,
  .language-badge.ts {
    background-color: rgba(43, 116, 137, 0.1);
    color: #2b7489;
    border-color: rgba(43, 116, 137, 0.2);
  }

  .language-badge.html {
    background-color: rgba(227, 76, 38, 0.1);
    color: #e34c26;
    border-color: rgba(227, 76, 38, 0.2);
  }

  .language-badge.css {
    background-color: rgba(86, 61, 124, 0.1);
    color: #563d7c;
    border-color: rgba(86, 61, 124, 0.2);
  }

  .language-badge.json {
    background-color: rgba(41, 128, 185, 0.1);
    color: #2980b9;
    border-color: rgba(41, 128, 185, 0.2);
  }

  .language-badge.markdown,
  .language-badge.md {
    background-color: rgba(163, 113, 247, 0.1);
    color: var(--accent-purple);
    border-color: rgba(163, 113, 247, 0.2);
  }

  pre {
    margin: 0;
    padding: var(--space-4);
    overflow-x: auto;
    font-family: var(--font-mono);
    font-size: var(--text-sm);
    line-height: 1.5;
    background-color: var(--bg-primary);
  }

  code {
    font-family: var(--font-mono);
    font-size: var(--text-sm);
    color: var(--text-primary);
  }

  /* Inline Code */
  p code,
  li code {
    background-color: var(--bg-tertiary);
    padding: var(--space-1) var(--space-2);
    border-radius: var(--radius-sm);
    border: 1px solid var(--border-primary);
    font-size: 0.9em;
  }

  /* Buttons */
  .btn {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-4);
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 500;
    line-height: 1.5;
    text-decoration: none;
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-md);
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
    cursor: pointer;
    transition: all var(--transition-base);
    white-space: nowrap;
  }

  .btn:hover {
    background-color: var(--bg-hover);
    border-color: var(--text-tertiary);
    transform: translateY(-1px);
    box-shadow: var(--shadow-sm);
  }

  .btn:active {
    transform: translateY(0);
  }

  .btn-primary {
    background-color: var(--accent-blue);
    color: #ffffff;
    border-color: var(--accent-blue);
  }

  .btn-primary:hover {
    background-color: var(--accent-blue-hover);
    border-color: var(--accent-blue-hover);
    color: #ffffff;
  }

  .btn-success {
    background-color: var(--accent-green);
    color: #ffffff;
    border-color: var(--accent-green);
  }

  .btn-success:hover {
    background-color: #4aca65;
    border-color: #4aca65;
    color: #ffffff;
  }

  /* Links */
  a {
    color: var(--accent-blue);
    text-decoration: none;
    transition: color var(--transition-fast);
  }

  a:hover {
    color: var(--accent-blue-hover);
  }

  /* Tables */
  table {
    width: 100%;
    border-collapse: collapse;
    margin: var(--space-4) 0;
  }

  th {
    text-align: left;
    padding: var(--space-3) var(--space-4);
    background-color: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-primary);
    font-weight: 600;
    color: var(--text-primary);
    font-size: var(--text-sm);
  }

  td {
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--border-secondary);
    color: var(--text-secondary);
  }

  tr:hover {
    background-color: var(--bg-tertiary);
  }

  /* Loading State */
  .loading {
    display: inline-block;
    width: 16px;
    height: 16px;
    border: 2px solid var(--border-primary);
    border-top-color: var(--accent-blue);
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Error Message */
  .error {
    padding: var(--space-4);
    background-color: rgba(248, 81, 73, 0.1);
    border: 1px solid rgba(248, 81, 73, 0.3);
    border-radius: var(--radius-md);
    color: #f85149;
    margin: var(--space-4) 0;
  }

  /* Success Message */
  .success {
    padding: var(--space-4);
    background-color: rgba(63, 185, 80, 0.1);
    border: 1px solid rgba(63, 185, 80, 0.3);
    border-radius: var(--radius-md);
    color: var(--accent-green);
    margin: var(--space-4) 0;
  }

  /* Info Message */
  .info {
    padding: var(--space-4);
    background-color: rgba(88, 166, 255, 0.1);
    border: 1px solid rgba(88, 166, 255, 0.3);
    border-radius: var(--radius-md);
    color: var(--accent-blue);
    margin: var(--space-4) 0;
  }

  /* Scrollbar Styles */
  ::-webkit-scrollbar {
    width: 12px;
    height: 12px;
  }

  ::-webkit-scrollbar-track {
    background-color: var(--bg-secondary);
  }

  ::-webkit-scrollbar-thumb {
    background-color: var(--bg-hover);
    border-radius: var(--radius-md);
    border: 2px solid var(--bg-secondary);
  }

  ::-webkit-scrollbar-thumb:hover {
    background-color: var(--text-tertiary);
  }

  /* Selection */
  ::selection {
    background-color: var(--selection-bg, rgba(88, 166, 255, 0.3));
    color: var(--text-primary);
  }

  /* Light Theme - System Detection */
  @media (prefers-color-scheme: light) {
    :root:not([data-theme="dark"]) {
      --bg-primary: #ffffff;
      --bg-secondary: #f6f8fa;
      --bg-tertiary: #f0f3f6;
      --bg-hover: #e8ebef;
      --border-primary: #d0d7de;
      --border-secondary: #e8ebef;
      --text-primary: #1f2328;
      --text-secondary: #656d76;
      --text-tertiary: #8b949e;
      --accent-blue: #0969da;
      --accent-blue-hover: #0550ae;
      --accent-green: #1a7f37;
      --accent-purple: #8250df;
      --accent-orange: #bc4c00;
      --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.06);
      --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.08), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
      --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.08);
      --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.08);
      --selection-bg: rgba(9, 105, 218, 0.2);
    }
  }

  /* Light Theme - Explicit */
  [data-theme="light"] {
    --bg-primary: #ffffff;
    --bg-secondary: #f6f8fa;
    --bg-tertiary: #f0f3f6;
    --bg-hover: #e8ebef;
    --border-primary: #d0d7de;
    --border-secondary: #e8ebef;
    --text-primary: #1f2328;
    --text-secondary: #656d76;
    --text-tertiary: #8b949e;
    --accent-blue: #0969da;
    --accent-blue-hover: #0550ae;
    --accent-green: #1a7f37;
    --accent-purple: #8250df;
    --accent-orange: #bc4c00;
    --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.06);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.08), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.08);
    --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.08);
    --selection-bg: rgba(9, 105, 218, 0.2);
  }

  /* Dark Theme - Explicit (overrides prefers-color-scheme) */
  [data-theme="dark"] {
    --bg-primary: #0d1117;
    --bg-secondary: #161b22;
    --bg-tertiary: #21262d;
    --bg-hover: #30363d;
    --border-primary: #30363d;
    --border-secondary: #21262d;
    --text-primary: #e6edf3;
    --text-secondary: #8b949e;
    --text-tertiary: #6e7681;
    --accent-blue: #58a6ff;
    --accent-blue-hover: #79b8ff;
    --accent-green: #3fb950;
    --accent-purple: #a371f7;
    --accent-orange: #ff8c00;
    --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.3);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.4), 0 2px 4px -1px rgba(0, 0, 0, 0.3);
    --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.5), 0 4px 6px -2px rgba(0, 0, 0, 0.4);
    --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.6), 0 10px 10px -5px rgba(0, 0, 0, 0.4);
    --selection-bg: rgba(88, 166, 255, 0.3);
  }

  /* Responsive Design */
  @media (max-width: 768px) {
    body {
      padding: var(--space-3);
    }

    .container {
      border-radius: var(--radius-md);
    }

    header {
      padding: var(--space-4);
    }

    h1 {
      font-size: var(--text-2xl);
    }

    .content {
      padding: var(--space-4);
    }

    .breadcrumb {
      padding: var(--space-3);
      font-size: var(--text-xs);
    }

    .directory-item {
      padding: var(--space-3);
    }

    pre {
      padding: var(--space-3);
      font-size: var(--text-xs);
    }

    .code-header {
      padding: var(--space-2) var(--space-3);
    }
  }

  @media (max-width: 480px) {
    h1 {
      font-size: var(--text-xl);
    }

    .breadcrumb {
      flex-direction: column;
      align-items: flex-start;
    }

    .breadcrumb .separator {
      display: none;
    }
  }

  /* Print Styles */
  @media print {
    body {
      background-color: white;
      color: black;
    }

    .container {
      box-shadow: none;
      border: none;
    }

    .btn {
      display: none;
    }

    a {
      color: #0066cc;
      text-decoration: underline;
    }
  }

  /* Git Status Badges */
  .git-badge {
    display: inline-flex;
    align-items: center;
    padding: 0.125rem 0.5rem;
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    font-weight: 700;
    border-radius: var(--radius-sm);
    margin-right: var(--space-2);
    letter-spacing: 0.05em;
  }

  .git-badge.git-orange {
    background-color: rgba(255, 140, 0, 0.15);
    color: #ff8c00;
    border: 1px solid rgba(255, 140, 0, 0.3);
  }

  .git-badge.git-green {
    background-color: rgba(63, 185, 80, 0.15);
    color: #3fb950;
    border: 1px solid rgba(63, 185, 80, 0.3);
  }

  .git-badge.git-purple {
    background-color: rgba(163, 113, 247, 0.15);
    color: #a371f7;
    border: 1px solid rgba(163, 113, 247, 0.3);
  }

  .git-badge.git-red {
    background-color: rgba(248, 81, 73, 0.15);
    color: #f85149;
    border: 1px solid rgba(248, 81, 73, 0.3);
  }

  .git-badge.git-blue {
    background-color: rgba(88, 166, 255, 0.15);
    color: #58a6ff;
    border: 1px solid rgba(88, 166, 255, 0.3);
  }

  /* Git Branch Badge */
  .git-branch {
    display: inline-flex;
    align-items: center;
    gap: var(--space-1);
    padding: var(--space-1) var(--space-3);
    margin-left: var(--space-4);
    background-color: var(--bg-hover);
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-xl);
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    font-weight: 600;
    color: var(--accent-blue);
  }

  /* Repository Statistics Panel */
  .repo-stats {
    margin: var(--space-6) var(--space-6) 0;
    padding: var(--space-5);
    background-color: var(--bg-tertiary);
    border: 1px solid var(--border-primary);
    border-radius: var(--radius-lg);
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: var(--space-5);
  }

  .repo-stat {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .repo-stat-value {
    font-size: var(--text-2xl);
    font-weight: 700;
    color: var(--text-primary);
    font-family: var(--font-mono);
  }

  .repo-stat-label {
    font-size: var(--text-xs);
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    font-weight: 600;
  }

  /* Responsive Git Styles */
  @media (max-width: 768px) {
    .repo-stats {
      margin: var(--space-4) var(--space-4) 0;
      padding: var(--space-4);
      grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
      gap: var(--space-4);
    }

    .repo-stat-value {
      font-size: var(--text-xl);
    }

    .git-branch {
      margin-left: var(--space-2);
      padding: var(--space-1) var(--space-2);
    }
  }

  @media (max-width: 480px) {
    .repo-stats {
      grid-template-columns: repeat(2, 1fr);
    }

    .git-badge {
      font-size: 0.65rem;
      padding: 0.1rem 0.4rem;
      margin-right: var(--space-1);
    }
  }
`;

module.exports = { baseStyles };
