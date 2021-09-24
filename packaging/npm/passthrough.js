#!/usr/bin/env node

var path = require('path');

// os and arch restrictions are handled by the package.json
var os = process.platform;
var arch = 'amd64';

// Select the right binary for this platform, then exec it with the original
// arguments. This is a true exec(3), which will take over the pid, env, and
// file descriptors.
var iacCustomRulesPath = path.join(__dirname, './snyk-iac-rules-' + os + '-' + arch);
if (os === 'win32') {
  iacCustomRulesPath = path.join(__dirname, './snyk-iac-rules.exe');
}

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
