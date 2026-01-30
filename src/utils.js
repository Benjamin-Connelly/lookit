// General utilities
const { spawn } = require('child_process');

function parseArgs(argv) {
  const args = {};

  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];

    if (arg === '--port' && argv[i + 1]) {
      args.port = parseInt(argv[++i], 10);
    } else if (arg === '--host' && argv[i + 1]) {
      args.host = argv[++i];
    } else if (arg === '--cert' && argv[i + 1]) {
      args.cert = argv[++i];
    } else if (arg === '--key' && argv[i + 1]) {
      args.key = argv[++i];
    } else if (arg === '--open') {
      args.open = true;
    } else if (arg === '--https-only') {
      args.httpsOnly = true;
    } else if (arg === '--no-https') {
      args.noHttps = true;
    } else if (arg === '--no-dirlist') {
      args.noDirlist = true;
    } else if (arg === '--all') {
      args.showAll = true;
    } else if (arg === '--quiet' || arg === '-q') {
      args.quiet = true;
    } else if (arg === '--help' || arg === '-h') {
      printHelp();
      process.exit(0);
    }
  }

  return args;
}

function printHelp() {
  console.log(`
lookit - Browse and view code files with syntax highlighting

USAGE:
  lookit [OPTIONS]

OPTIONS:
  --port <number>      Port to listen on (default: 7777) - Lucky number! 🍀
  --host <address>     Host to bind to (default: 127.0.0.1)
  --open               Open browser automatically
  --all                Show hidden files (starting with .)
  --https-only         Fail if TLS certificates are not found
  --no-https           Use HTTP only, skip HTTPS even if certificates exist
  --no-dirlist         Disable directory listings
  --cert <path>        Path to TLS certificate (default: ~/.config/lookit/localhost.pem)
  --key <path>         Path to TLS private key (default: ~/.config/lookit/localhost-key.pem)
  -q, --quiet          Suppress TLS certificate warnings
  -h, --help           Show this help message

FILE SUPPORT:
  Markdown (.md)       Rendered as HTML with syntax-highlighted code blocks
  Code files           Displayed with syntax highlighting (YAML, JSON, JS, Python, etc.)
  Images/PDFs          Displayed natively in browser
  Other files          Downloaded

EXAMPLES:
  lookit                           # Serve current directory on https://localhost:7777
  lookit --port 8080               # Use port 8080
  lookit --no-https                # Use HTTP only
  lookit --open                    # Auto-open browser
  lookit --https-only              # Require HTTPS, fail if no certs
  lookit --all                     # Show hidden files

TLS SETUP (Ubuntu):
  sudo apt install -y mkcert libnss3-tools
  mkcert -install
  mkdir -p ~/.config/lookit
  mkcert -cert-file ~/.config/lookit/localhost.pem \\
         -key-file ~/.config/lookit/localhost-key.pem \\
         localhost 127.0.0.1 ::1
`);
}

function printCertInstructions() {
  console.log(`📝 To enable HTTPS, install mkcert and generate certificates:\n`);
  console.log(`   1. Install mkcert:`);
  console.log(`      sudo apt install -y mkcert libnss3-tools\n`);
  console.log(`   2. Install the local CA:`);
  console.log(`      mkcert -install\n`);
  console.log(`   3. Generate certificates:`);
  console.log(`      mkdir -p ~/.config/lookit`);
  console.log(`      mkcert -cert-file ~/.config/lookit/localhost.pem \\`);
  console.log(`             -key-file ~/.config/lookit/localhost-key.pem \\`);
  console.log(`             localhost 127.0.0.1 ::1`);
}

function openBrowser(url) {
  const commands = {
    linux: 'xdg-open',
    darwin: 'open',
    win32: 'start'
  };

  const command = commands[process.platform];
  if (command) {
    spawn(command, [url], { detached: true, stdio: 'ignore' }).unref();
  }
}

function escapeHtml(text) {
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return text.replace(/[&<>"']/g, m => map[m]);
}

module.exports = {
  parseArgs,
  printHelp,
  printCertInstructions,
  openBrowser,
  escapeHtml
};