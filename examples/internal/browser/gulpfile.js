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
  return new Promise((resolve, reject) => {
    const serverPath = path.join(__dirname, 'bin', 'example-server');
    console.log(`Starting server: ${serverPath}`);
    serverProcess = spawn(serverPath, [], { 
      stdio: ['ignore', 'pipe', 'pipe'],
      cwd: __dirname
    });
    
    serverProcess.stdout.on('data', (data) => {
      console.log(`[Server]: ${data.toString().trim()}`);
    });
    
    serverProcess.stderr.on('data', (data) => {
      console.error(`[Server Error]: ${data.toString().trim()}`);
    });
    
    serverProcess.on('error', (error) => {
      console.error(`Failed to start server: ${error.message}`);
      reject(error);
    });
    
    serverProcess.on('exit', (code, signal) => {
      console.log(`Server exited with code ${code} and signal ${signal}`);
    });
    
    // Give server time to start
    setTimeout(resolve, 2000);
  });
}

async function startGateway() {
  await buildGateway();
  return new Promise((resolve, reject) => {
    const gatewayPath = path.join(__dirname, 'bin', 'example-gw');
    const openApiDir = path.join(__dirname, '..', 'proto', 'examplepb');
    console.log(`Starting gateway: ${gatewayPath} with openapi_dir: ${openApiDir}`);
    gatewayProcess = spawn(gatewayPath, [
      '--openapi_dir', openApiDir
    ], { 
      stdio: ['ignore', 'pipe', 'pipe'],
      cwd: __dirname
    });
    
    gatewayProcess.stdout.on('data', (data) => {
      console.log(`[Gateway]: ${data.toString().trim()}`);
    });
    
    gatewayProcess.stderr.on('data', (data) => {
      console.error(`[Gateway Error]: ${data.toString().trim()}`);
    });
    
    gatewayProcess.on('error', (error) => {
      console.error(`Failed to start gateway: ${error.message}`);
      reject(error);
    });
    
    gatewayProcess.on('exit', (code, signal) => {
      console.log(`Gateway exited with code ${code} and signal ${signal}`);
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

// Generate HTML content for Jasmine test runner
function generateJasmineHTML(jasmineCore, jasmineHtml, includeResultsCapture = false) {
  const resultsScript = includeResultsCapture 
    ? 'window.jasmineResults = result;' 
    : '';
  
  return `
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>Jasmine Spec Runner</title>
      <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .jasmine_html-reporter { margin: 0; padding: 0; }
      </style>
    </head>
    <body>
      <script>${jasmineCore}</script>
      <script>${jasmineHtml}</script>
      <script>
        window.jasmine = jasmineRequire.core(jasmineRequire);
        var env = jasmine.getEnv();
        
        // Console reporter for logging
        var consoleReporter = {
          jasmineStarted: function() { console.log('Jasmine started'); },
          suiteStarted: function(result) { console.log('Suite: ' + result.description); },
          specStarted: function(result) { console.log('  Spec: ' + result.description); },
          specDone: function(result) {
            console.log('  Spec done: ' + result.description + ' - ' + result.status);
            if (result.status === 'failed') {
              result.failedExpectations.forEach(function(expectation) {
                console.log('    FAILED: ' + expectation.message);
              });
            }
          },
          suiteDone: function(result) {},
          jasmineDone: function(result) {
            console.log('Jasmine done: ' + result.overallStatus);
            ${resultsScript}
          }
        };
        env.addReporter(consoleReporter);
        
        // Make jasmine interface available globally
        var jasmineInterface = jasmineRequire.interface(jasmine, env);
        Object.assign(window, jasmineInterface);
      </script>
      <script src="/spec.js"></script>
      <script>
        env.execute();
      </script>
    </body>
    </html>
  `;
}

// Load Jasmine library files
function loadJasmineFiles() {
  const fs = require('fs');
  const jasmineCorePath = require.resolve('jasmine-core/lib/jasmine-core/jasmine.js');
  const jasmineHtmlPath = require.resolve('jasmine-core/lib/jasmine-core/jasmine-html.js');
  
  return {
    jasmineCore: fs.readFileSync(jasmineCorePath, 'utf8'),
    jasmineHtml: fs.readFileSync(jasmineHtmlPath, 'utf8')
  };
}

// Create HTTP server that serves test page and spec bundle
function createTestServer(htmlContent, specBundle) {
  const http = require('http');
  
  return http.createServer((req, res) => {
    if (req.url === '/spec.js') {
      res.writeHead(200, { 'Content-Type': 'application/javascript' });
      res.end(specBundle);
    } else {
      res.writeHead(200, { 'Content-Type': 'text/html' });
      res.end(htmlContent);
    }
  });
}

async function runTests() {
  const fs = require('fs');
  const { chromium } = require('playwright');
  
  let server = null;
  let browser = null;
  
  try {
    // Start the backends
    console.log('Starting server and gateway...');
    await startServer();
    await startGateway();
    
    // Bundle the specs
    await bundleSpecs();
    
    // Set up test server
    console.log('Setting up test server...');
    
    // Load Jasmine and the bundled specs
    const { jasmineCore, jasmineHtml } = loadJasmineFiles();
    const specBundle = fs.readFileSync(path.join(__dirname, 'bin', 'spec.js'), 'utf8');

    // Create HTML page content with results capture
    const htmlContent = generateJasmineHTML(jasmineCore, jasmineHtml, true);

    // Create HTTP server
    server = createTestServer(htmlContent, specBundle);

    // Start server on port 8000
    await new Promise((resolve) => {
      server.listen(8000, () => {
        console.log('Test server listening on http://localhost:8000');
        resolve();
      });
    });

    // Run tests in Playwright
    console.log('Launching browser...');
    browser = await chromium.launch({ headless: true });
    const context = await browser.newContext();
    const page = await context.newPage();

    console.log('Setting page timeout...');
    page.setDefaultTimeout(120000);  // Set default timeout to 120 seconds

    // Capture console output
    page.on('console', msg => {
      const text = msg.text();
      console.log('[Browser]:', text);
    });

    // Capture page errors
    page.on('pageerror', error => {
      console.error('[Page Error]:', error.message);
    });

    console.log('Navigating to test page...');
    // Navigate to the test server
    await page.goto('http://localhost:8000');

    console.log('Waiting for tests to complete...');
    // Wait for tests to complete
    await page.waitForFunction(() => window.jasmineResults !== undefined, { timeout: 120000 });

    // Get results
    const results = await page.evaluate(() => window.jasmineResults);
    
    await browser.close();
    browser = null;

    console.log(`\nTest Results: ${results.overallStatus}`);
    console.log(`Total specs: ${results.totalCount || 'N/A'}`);
    console.log(`Failed specs: ${results.failedExpectations?.length || 0}`);

    if (results.overallStatus !== 'passed') {
      throw new Error('Tests failed');
    }
    
    console.log('Tests completed successfully');
  } catch (error) {
    console.error('Tests failed:', error.message);
    throw error;
  } finally {
    // Cleanup
    if (browser) {
      await browser.close();
    }
    if (server) {
      server.close();
    }
    cleanupProcesses();
  }
}

async function serve() {
  const fs = require('fs');
  
  try {
    // Start the backends
    console.log('Starting server and gateway...');
    await startServer();
    await startGateway();
    
    // Bundle the specs
    await bundleSpecs();
    
    // Start a development server
    console.log('Starting development server on http://localhost:8000...');
    
    // Load Jasmine files
    const { jasmineCore, jasmineHtml } = loadJasmineFiles();
    
    // Create server that dynamically loads spec bundle on each request (for development)
    const http = require('http');
    const server = http.createServer((req, res) => {
      if (req.url === '/spec.js') {
        const specBundle = fs.readFileSync(path.join(__dirname, 'bin', 'spec.js'), 'utf8');
        res.writeHead(200, { 'Content-Type': 'application/javascript' });
        res.end(specBundle);
        return;
      }
      
      // Generate HTML without results capture for development
      const htmlContent = generateJasmineHTML(jasmineCore, jasmineHtml, false);
      res.writeHead(200, { 'Content-Type': 'text/html' });
      res.end(htmlContent);
    });
    
    server.listen(8000, () => {
      console.log('Development server running at http://localhost:8000');
      console.log('Backend server running at http://localhost:9090');
      console.log('Gateway running at http://localhost:8080');
      console.log('Press Ctrl+C to stop');
    });
    
    // Keep the process running
    await new Promise(() => {});
  } catch (error) {
    console.error('Serve task failed:', error.message);
    throw error;
  } finally {
    cleanupProcesses();
  }
}

exports.build = build;
exports.buildServer = buildServer;
exports.buildGateway = buildGateway;
exports.test = runTests;
exports.serve = serve;
exports.default = runTests;
