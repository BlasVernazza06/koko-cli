#!/usr/bin/env node

const os = require('os');
const fs = require('fs');
const { execFileSync } = require('child_process');
const path = require('path');

const platform = os.platform();
const arch = os.arch();

let binaryName = '';

if (platform === 'win32') {
  binaryName = 'koko-windows-amd64.exe';
} else if (platform === 'darwin') {
  if (arch === 'arm64') {
    binaryName = 'koko-darwin-arm64';
  } else {
    binaryName = 'koko-darwin-amd64';
  }
} else if (platform === 'linux') {
  binaryName = 'koko-linux-amd64';
} else {
  console.error(`Unsupported platform/architecture: ${platform}-${arch}`);
  process.exit(1);
}

const binaryPath = path.join(__dirname, 'bin', binaryName);

// Ensure binary exists
if (!fs.existsSync(binaryPath)) {
  console.error(`Koko binary not found for platform ${platform}-${arch} at: ${binaryPath}`);
  process.exit(1);
}

// Ensure execution permissions on Unix systems (macOS / Linux)
if (platform !== 'win32') {
  try {
    fs.chmodSync(binaryPath, 0o755);
  } catch (_) {
    // Ignore if permission change is not permitted
  }
}

try {
  // Execute the native binary with all forwarded arguments
  execFileSync(binaryPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (err) {
  process.exit(err.status || 1);
}
