const { listInstances } = require('./instanceManager');

// Handle --list command
function handleListCommand() {
  const instances = listInstances();
  const ports = Object.keys(instances);

  if (ports.length === 0) {
    console.log('📭 No running lookit instances');
    return;
  }

  console.log(`📋 Running lookit instances (${ports.length}):\n`);

  ports.forEach(port => {
    const instance = instances[port];
    const protocol = instance.protocol === 'https' ? '🔒' : '🌐';
    const uptime = getUptime(instance.started);

    console.log(`  ${protocol} Port ${port}`);
    console.log(`     Directory: ${instance.dir}`);
    console.log(`     PID: ${instance.pid}`);
    console.log(`     Uptime: ${uptime}`);
    console.log(`     URL: ${instance.protocol}://localhost:${port}`);
    console.log('');
  });
}

// Handle --stop <port> command
function handleStopCommand(port) {
  const instances = listInstances();
  const instance = instances[port];

  if (!instance) {
    console.error(`❌ No lookit instance running on port ${port}`);
    console.log(`\nRunning instances:`);
    handleListCommand();
    process.exit(1);
  }

  try {
    process.kill(instance.pid, 'SIGTERM');
    console.log(`✅ Stopped lookit on port ${port} (PID ${instance.pid})`);
  } catch (err) {
    if (err.code === 'ESRCH') {
      console.log(`⚠️  Process ${instance.pid} not found (already stopped?)`);
    } else {
      console.error(`❌ Failed to stop process: ${err.message}`);
      process.exit(1);
    }
  }
}

// Handle --stop-all command
function handleStopAllCommand() {
  const instances = listInstances();
  const ports = Object.keys(instances);

  if (ports.length === 0) {
    console.log('📭 No running lookit instances to stop');
    return;
  }

  console.log(`🛑 Stopping ${ports.length} lookit instance(s)...\n`);

  let stopped = 0;
  let failed = 0;

  ports.forEach(port => {
    const instance = instances[port];
    try {
      process.kill(instance.pid, 'SIGTERM');
      console.log(`  ✅ Stopped port ${port} (PID ${instance.pid})`);
      stopped++;
    } catch (err) {
      if (err.code === 'ESRCH') {
        console.log(`  ⚠️  Port ${port} (PID ${instance.pid}) - already stopped`);
        stopped++;
      } else {
        console.log(`  ❌ Port ${port} - failed: ${err.message}`);
        failed++;
      }
    }
  });

  console.log(`\n✅ Stopped ${stopped} instance(s)`);
  if (failed > 0) {
    console.log(`⚠️  Failed to stop ${failed} instance(s)`);
  }
}

// Helper: calculate uptime
function getUptime(startedISO) {
  const started = new Date(startedISO);
  const now = new Date();
  const diffMs = now - started;
  const diffMins = Math.floor(diffMs / 60000);

  if (diffMins < 1) return 'just now';
  if (diffMins < 60) return `${diffMins}m`;

  const diffHours = Math.floor(diffMins / 60);
  const mins = diffMins % 60;
  return `${diffHours}h ${mins}m`;
}

module.exports = {
  handleListCommand,
  handleStopCommand,
  handleStopAllCommand
};
