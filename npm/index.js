#!/usr/bin/env node

const os = require('os');
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

try {
  // Execute the native binary with all forwarded arguments
  execFileSync(binaryPath, process.argv.slice(2), { stdio: 'inherit' });
} catch (err) {
  process.exit(err.status || 1);
}
