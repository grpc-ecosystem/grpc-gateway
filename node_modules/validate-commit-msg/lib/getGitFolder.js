'use strict';

var fs = require('fs');
var path = require('path');
var findParentDir = require('find-parent-dir');

module.exports = function getGitFolder() {
  var dir = findParentDir.sync(process.cwd(), '.git');
  if (!dir) throw new Error('Cannot find .git folder');

  var gitDir = path.join(dir, '.git');
  var stats = fs.lstatSync(gitDir);

  if (!stats.isDirectory()) {
    // Expect following format
    // git: pathToGit
    var pathToGit = fs
      .readFileSync(gitDir, 'utf-8')
      .split(':')[1]
      .trim();
    gitDir = path.join(dir, pathToGit);

    if (!fs.existsSync(gitDir)) {
      throw new Error('Cannot find file ' + pathToGit);
    }
  }

  return gitDir;
};
