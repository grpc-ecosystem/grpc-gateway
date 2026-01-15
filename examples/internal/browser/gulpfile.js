const { exec, spawn } = require('child_process');
const { parallel, series } = require('gulp');
const { join } = require('path');
const { stat, readFile } = require('fs/promises');
const { H3, handleCors, serve, serveStatic } = require('h3');

const buildGS = () => {
  const cmd = 'go build -o bin/example-server github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-grpc-server';
  return exec(cmd, (error, stdout, stderr) => {
    if (error) {
      console.error(`Clean failed: ${error}`);
    } else {
      console.log('build grpc server success');
    }
  });
}

const buildGW = () => {
  const cmd = 'go build -o bin/example-gw github.com/grpc-ecosystem/grpc-gateway/v2/examples/internal/cmd/example-gateway-server';
  return exec(cmd, (error, stdout, stderr) => {
    if (error) {
      console.error(`Clean failed: ${error}`);
    } else {
      console.log('build gw server success');
    }
  });
}

const runGS = (done, cb) => {
  const gs = spawn('bin/example-server', [], { stdio: 'inherit' }, cb);
  process.on('exit', () => {
    gs.kill();
  });
  done();
}

const runGW = (done, cb) => {
  const gw = spawn('bin/example-gw', ['--openapi_dir', join(__dirname, "../proto/examplepb"),], { stdio: 'inherit' }, cb);
  process.on('exit', () => {
    gw.kill();
  });
  done();
}

const server = (done) => {
  const app = new H3({ debug: true });
  return serve(
    app.all('/**', (event) => {
      if (handleCors(event, { origin: "*" })) {
        return;
      }
      return serveStatic(event, {
        indexNames: ['/index.html'],
        getContents: (id) => readFile(join('public', id)),
        getMeta: async (id) => {
          const stats = await stat(join('public', id)).catch(() => { });
          if (stats?.isFile()) {
            return {
              size: stats.size,
              mtime: stats.mtimeMs,
            };
          }
        },
      });
    }
    )
  );
}

exports.default = series(parallel(buildGS, buildGW), runGS, runGW, server);
