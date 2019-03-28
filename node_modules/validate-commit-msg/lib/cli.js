#!/usr/bin/env node

/**
 * Git COMMIT-MSG hook for validating commit message
 * See https://docs.google.com/document/d/1rk04jEuGfk9kYzfqCuOlPTSJw3hEDZJTBN5E5f1SALo/edit
 *
 * This CLI supports 3 usage ways:
 * 1. Default usage is not passing any argument. It will automatically read from COMMIT_EDITMSG file.
 * 2. Passing a file name argument from git directory. For instance GIT GUI stores commit msg @GITGUI_EDITMSG file.
 * 3. Passing commit message as argument. Useful for testing quickly a commit message from CLI.
 *
 * Installation:
 * >> cd <angular-repo>
 * >> ln -s ../../validate-commit-msg.js .git/hooks/commit-msg
 */

'use strict';

var fs = require('fs');

var getGitFolder = require('./getGitFolder');
var validateMessage = require('../index');

// hacky start if not run by mocha :-D
// istanbul ignore next
if (process.argv.join('').indexOf('mocha') === -1) {
  var bufferToString = function (buffer) {
    var hasToString = buffer && typeof buffer.toString === 'function';
    return hasToString && buffer.toString();
  };

  var getFileContent = function (filePath) {
    try {
      var buffer = fs.readFileSync(filePath);
      return bufferToString(buffer);
    } catch (err) {
      // Ignore these error types because is most likely it is validating
      // a commit from a text instead of a file
      if(err && err.code !== 'ENOENT' && err.code !== 'ENAMETOOLONG') {
        throw err;
      }
    }
  };

  var getCommit = function() {
    var file;
    var fileContent;
    var gitDirectory;

    var commitMsgFileOrText = process.argv[2];
    var commitErrorLogPath = process.argv[3];

    var commit = {
      // if it is running from git directory or for a file from there
      // these info might change ahead
      message: commitMsgFileOrText,
      errorLog: commitErrorLogPath || null,
      file: null,
    };

    // On running the validation over a text instead of git files such as COMMIT_EDITMSG and GITGUI_EDITMSG
    // is possible to be doing that the from anywhere. Therefore the git directory might not be available.
    try {
      gitDirectory = getGitFolder();

      // Try to load commit from a path passed as argument
      if (commitMsgFileOrText) {
        file = gitDirectory + '/' + commitMsgFileOrText;
        fileContent = getFileContent(file);
      }

      // If no file or message is available then try to load it from the default commit file
      if (!fileContent && !commitMsgFileOrText) {
        file = gitDirectory + '/COMMIT_EDITMSG';
        fileContent = getFileContent(file);
      }

      // Could resolve the content from a file
      if (fileContent) {
        commit.file = file;
        commit.message = fileContent;
      }

      // Default error log path
      if (!commit.errorLog) {
        commit.errorLog = gitDirectory + '/logs/incorrect-commit-msgs';
      }
    } catch (err) {}

    return commit;
  };

  var validate = function (commit) {
    if (!validateMessage(commit.message, commit.file) && commit.errorLog) {
      fs.appendFileSync(commit.errorLog, commit.message + '\n');
      process.exit(1);
    } else {
      process.exit(0);
    }
  };
  validate(getCommit());
}
