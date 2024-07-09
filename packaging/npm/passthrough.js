#!/usr/bin/env node

var path = require('path');

// os and arch restrictions are handled by the package.json
var os = '';
switch (process.platform) {
  case 'darwin':
    os = 'darwin';
    break;
  case 'win32':
    os = 'windows';
    break;
  default:
    os = 'linux';
};

var arch = ''
switch (process.arch) {
  case 'arm' :
  case 'arm64':
    arch = 'arm64';
    break;
  case 'x64':
    arch = 'amd64_v1';
    break;
  default:
    throw new Error(`Architecture not supported: ${process.arch}`)
}

// Select the right binary for this platform, then exec it with the original
// arguments. This is a true exec(3), which will take over the pid, env, and
// file descriptors.
var iacCustomRulesPath = path.join(__dirname, './snyk-iac-rules' + '_' + os + '_' + arch, 'snyk-iac-rules');

try {
  var spawn = require('child_process').spawn;
  var proc = spawn(iacCustomRulesPath, process.argv.slice(2), { stdio: 'inherit' });
  proc.on('exit', function (code, signal) {
    process.on('exit', function () {
      if (signal) {
        process.kill(process.pid, signal);
      } else {
        process.exit(code);
      }
    });
  });
} catch (err) {
  console.error(err);
  process.exit(code);
}
