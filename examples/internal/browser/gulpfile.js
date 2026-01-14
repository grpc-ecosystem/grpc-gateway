"use strict";

const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

// Module paths for the example builds
const SERVER_MODULE = 'github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-grpc-server';
const GATEWAY_MODULE = 'github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-gateway-server';

// Validate module path to ensure it's safe
function validateModulePath(modulePath) {
  // Module path should only contain alphanumeric, hyphens, underscores, dots, and slashes
  if (!/^[a-zA-Z0-9\-_.\/]+$/.test(modulePath)) {
    throw new Error(`Invalid module path: ${modulePath}`);
  }
}

async function buildServer() {
  console.log('Building example server...');
  validateModulePath(SERVER_MODULE);
  try {
    const { stdout, stderr } = await execPromise(
      `go build -o bin/example-server ${SERVER_MODULE}`
    );
    // Only log if there's actual output (Go builds are usually silent on success)
    if (stdout) console.log(stdout);
    if (stderr) console.error(stderr);
    console.log('Server built successfully');
  } catch (error) {
    console.error('Failed to build server:', error.message);
    throw error;
  }
}

async function buildGateway() {
  console.log('Building gateway...');
  validateModulePath(GATEWAY_MODULE);
  try {
    const { stdout, stderr } = await execPromise(
      `go build -o bin/example-gw ${GATEWAY_MODULE}`
    );
    // Only log if there's actual output (Go builds are usually silent on success)
    if (stdout) console.log(stdout);
    if (stderr) console.error(stderr);
    console.log('Gateway built successfully');
  } catch (error) {
    console.error('Failed to build gateway:', error.message);
    throw error;
  }
}

async function build() {
  await buildServer();
  await buildGateway();
  console.log('All builds completed successfully');
}

exports.build = build;
exports.default = build;
