"use strict";

const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

// Module paths for the example builds
const SERVER_MODULE = 'github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-grpc-server';
const GATEWAY_MODULE = 'github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-gateway-server';

async function buildServer() {
  console.log('Building example server...');
  try {
    const { stdout, stderr } = await execPromise(
      `go build -o bin/example-server ${SERVER_MODULE}`
    );
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
  try {
    const { stdout, stderr } = await execPromise(
      `go build -o bin/example-gw ${GATEWAY_MODULE}`
    );
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
