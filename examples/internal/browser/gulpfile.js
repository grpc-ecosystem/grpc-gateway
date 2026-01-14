"use strict";

const { exec, spawn } = require('child_process');
const util = require('util');
const path = require('path');
const execPromise = util.promisify(exec);

// Module paths for the example builds
const SERVER_MODULE = 'github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-grpc-server';
const GATEWAY_MODULE = 'github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-gateway-server';

let serverProcess = null;
let gatewayProcess = null;

// Validate module path to ensure it's safe
function validateModulePath(modulePath) {
  // Module path should only contain alphanumeric, hyphens, underscores, dots, and slashes
  if (!/^[a-zA-Z0-9\-_.\/]+$/.test(modulePath)) {
    throw new Error(`Invalid module path: ${modulePath}`);
  }
}

function cleanupProcesses() {
  if (serverProcess) {
    serverProcess.kill();
    serverProcess = null;
  }
  if (gatewayProcess) {
    gatewayProcess.kill();
    gatewayProcess = null;
  }
}

process.on('exit', cleanupProcesses);
process.on('SIGINT', () => {
  cleanupProcesses();
  process.exit();
});

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

async function startServer() {
  await buildServer();
  return new Promise((resolve) => {
    serverProcess = spawn('bin/example-server', [], { 
      stdio: 'inherit',
      cwd: __dirname
    });
    // Give server time to start
    setTimeout(resolve, 2000);
  });
}

async function startGateway() {
  await buildGateway();
  return new Promise((resolve) => {
    gatewayProcess = spawn('bin/example-gw', [
      '--openapi_dir', path.join(__dirname, '../proto/examplepb')
    ], { 
      stdio: 'inherit',
      cwd: __dirname
    });
    // Give gateway time to start
    setTimeout(resolve, 2000);
  });
}

async function bundleSpecs() {
  console.log('Bundling test specs with webpack...');
  const specFiles = [
    path.join(__dirname, 'a_bit_of_everything_service.spec.js'),
    path.join(__dirname, 'echo_service.spec.js')
  ];
  
  try {
    const { stdout, stderr } = await execPromise(
      `npx webpack --config webpack.config.js --entry ${specFiles.join(' --entry ')} --output-path ${path.join(__dirname, 'bin')} --output-filename spec.js`,
      { cwd: __dirname }
    );
    if (stdout) console.log(stdout);
    if (stderr && !stderr.includes('webpack')) console.error(stderr);
    console.log('Specs bundled successfully');
  } catch (error) {
    console.error('Failed to bundle specs:', error.message);
    throw error;
  }
}

async function runTests() {
  try {
    // Start the backends
    console.log('Starting server and gateway...');
    await startServer();
    await startGateway();
    
    // Bundle the specs
    await bundleSpecs();
    
    // Run jasmine tests in browser using Playwright
    console.log('Running tests with Playwright...');
    const { stdout, stderr } = await execPromise(
      'node test-runner.js',
      { cwd: __dirname }
    );
    if (stdout) console.log(stdout);
    if (stderr) console.error(stderr);
    console.log('Tests completed successfully');
  } catch (error) {
    console.error('Tests failed:', error.message);
    throw error;
  } finally {
    cleanupProcesses();
  }
}

exports.build = build;
exports.buildServer = buildServer;
exports.buildGateway = buildGateway;
exports.test = runTests;
exports.default = runTests;
