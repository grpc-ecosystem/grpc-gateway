const { chromium } = require('playwright');
const path = require('path');
const fs = require('fs');

async function runTests() {
  console.log('Launching browser...');
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext();
  const page = await context.newPage();

  console.log('Setting page timeout...');
  page.setDefaultTimeout(120000);  // Set default timeout to 120 seconds

  // Capture console output (set before loading page)
  page.on('console', msg => {
    const text = msg.text();
    console.log('[Browser]:', text);
  });

  // Capture page errors
  page.on('pageerror', error => {
    console.error('[Page Error]:', error.message);
  });

  console.log('Loading test page...');

  // Load Jasmine and the bundled specs
  const jasmineCorePath = require.resolve('jasmine-core/lib/jasmine-core/jasmine.js');
  const jasmineCssPath = require.resolve('jasmine-core/lib/jasmine-core/jasmine.css');
  const jasmineHtmlPath = require.resolve('jasmine-core/lib/jasmine-core/jasmine-html.js');
  const jasmineBootPath = require.resolve('jasmine-core/lib/jasmine-core/boot1.js');
  
  console.log('Reading jasmine files...');
  const jasmineCore = fs.readFileSync(jasmineCorePath, 'utf8');
  const jasmineHtml = fs.readFileSync(jasmineHtmlPath, 'utf8');
  const jasmineBoot = fs.readFileSync(jasmineBootPath, 'utf8');
  const specBundle = fs.readFileSync(path.join(__dirname, 'bin', 'spec.js'), 'utf8');

  console.log('Creating test page...');
  // Create a test page
  await page.setContent(`
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
            window.jasmineResults = result;
          }
        };
        env.addReporter(consoleReporter);
        
        // Make jasmine interface available globally
        var jasmineInterface = jasmineRequire.interface(jasmine, env);
        Object.assign(window, jasmineInterface);
      </script>
      <script>${specBundle}</script>
      <script>
        env.execute();
      </script>
    </body>
    </html>
  `);

  console.log('Waiting for tests to complete...');
  // Wait for tests to complete (increased timeout for potentially slow tests)
  await page.waitForFunction(() => window.jasmineResults !== undefined, { timeout: 120000 });

  // Get results
  const results = await page.evaluate(() => window.jasmineResults);
  
  await browser.close();

  console.log(`\nTest Results: ${results.overallStatus}`);
  console.log(`Total specs: ${results.totalCount}`);
  console.log(`Failed specs: ${results.failedCount || results.overallStatus === 'failed' ? 'some' : 0}`);

  if (results.overallStatus !== 'passed') {
    process.exit(1);
  }
}

runTests().catch(error => {
  console.error('Test runner failed:', error);
  process.exit(1);
});
