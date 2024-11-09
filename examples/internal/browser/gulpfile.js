"use strict";

var gulp = require('gulp');

var path = require('path');

var bower = require('gulp-bower');
var exit = require('gulp-exit');
var shell = require('gulp-shell');
var jasmineBrowser = require('gulp-jasmine-browser');
var webpack = require('webpack-stream');
const child = require('child_process');

gulp.task('bower', function () {
  return bower();
});

gulp.task('server', shell.task([
  'go build -o bin/example-server github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-grpc-server',
]));

gulp.task('gateway', shell.task([
  'go build -o bin/example-gw github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-gateway-server',
]));

gulp.task('serve-server', ['server'], function () {
  let server = child.spawn('bin/example-server', [], { stdio: 'inherit' });
  process.on('exit', function () {
    server.kill();
  });
});

gulp.task('serve-gateway', ['gateway', 'serve-server'], function () {
  let gw = child.spawn('bin/example-gw', [
    '--openapi_dir', path.join(__dirname, "../proto/examplepb"),
  ], { stdio: 'inherit' });
  process.on('exit', function () {
    gw.kill();
  });
});

gulp.task('backends', ['serve-gateway', 'serve-server']);

var specFiles = ['*.spec.js'];
gulp.task('test', ['backends'], function (done) {
  let s = gulp.src(specFiles)
  console.log(s);
  return s
    .pipe(webpack({ output: { filename: 'spec.js' } }))
    .pipe(jasmineBrowser.specRunner({
      console: true,
      sourceMappedStacktrace: true,
    }))
    .pipe(jasmineBrowser.headless({
      driver: 'phantomjs',
      findOpenPort: true,
      catch: true,
      throwFailures: true,
    }))
    .on('error', function (err) {
      done(err);
      process.exit(1);
    })
    .pipe(exit());
});

gulp.task('serve', ['backends'], function (done) {
  var JasminePlugin = require('gulp-jasmine-browser/webpack/jasmine-plugin');
  var plugin = new JasminePlugin();

  return gulp.src(specFiles)
    .pipe(webpack({
      output: { filename: 'spec.js' },
      watch: true,
      plugins: [plugin],
    }))
    .pipe(jasmineBrowser.specRunner({
      sourceMappedStacktrace: true,
    }))
    .pipe(jasmineBrowser.server({
      port: 8000,
      whenReady: plugin.whenReady,
    }));
});

gulp.task('default', ['test']);
