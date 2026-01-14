"use strict";

const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

async function buildServer() {
  console.log('Building example server...');
  const { stdout, stderr } = await execPromise(
    'go build -o bin/example-server github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-grpc-server'
  );
  if (stdout) console.log(stdout);
  if (stderr) console.error(stderr);
  console.log('Server built successfully');
}

async function buildGateway() {
  console.log('Building gateway...');
  const { stdout, stderr } = await execPromise(
    'go build -o bin/example-gw github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-gateway-server'
  );
  if (stdout) console.log(stdout);
  if (stderr) console.error(stderr);
  console.log('Gateway built successfully');
}

async function build() {
  await buildServer();
  await buildGateway();
  console.log('All builds completed successfully');
}

exports.build = build;
exports.default = build;
