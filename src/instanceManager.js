const fs = require('fs');
const path = require('path');
const os = require('os');

const CONFIG_DIR = path.join(os.homedir(), '.config', 'lookit');
const INSTANCES_FILE = path.join(CONFIG_DIR, 'instances.json');

// Ensure config directory exists
function ensureConfigDir() {
  if (!fs.existsSync(CONFIG_DIR)) {
    fs.mkdirSync(CONFIG_DIR, { recursive: true });
  }
}

// Read instances file
function readInstances() {
  ensureConfigDir();
  if (!fs.existsSync(INSTANCES_FILE)) {
    return {};
  }
  try {
    return JSON.parse(fs.readFileSync(INSTANCES_FILE, 'utf8'));
  } catch (err) {
    console.warn('⚠️  Could not read instances file, starting fresh');
    return {};
  }
}

// Write instances file
function writeInstances(instances) {
  ensureConfigDir();
  fs.writeFileSync(INSTANCES_FILE, JSON.stringify(instances, null, 2));
}

// Register this instance
function registerInstance(port, directory, protocol = 'http') {
  const instances = readInstances();
  instances[port] = {
    pid: process.pid,
    dir: directory,
    port: port,
    protocol: protocol,
    started: new Date().toISOString()
  };
  writeInstances(instances);
}

// Unregister this instance
function unregisterInstance(port) {
  const instances = readInstances();
  delete instances[port];
  writeInstances(instances);
}

// Clean up stale instances (PIDs that no longer exist)
function cleanStaleInstances() {
  const instances = readInstances();
  let cleaned = false;

  for (const [port, instance] of Object.entries(instances)) {
    try {
      // Check if process exists (throws if not)
      process.kill(instance.pid, 0);
    } catch (err) {
      // Process doesn't exist, remove it
      delete instances[port];
      cleaned = true;
    }
  }

  if (cleaned) {
    writeInstances(instances);
  }
}

// Get all running instances
function listInstances() {
  cleanStaleInstances();
  return readInstances();
}

module.exports = {
  registerInstance,
  unregisterInstance,
  cleanStaleInstances,
  listInstances
};
