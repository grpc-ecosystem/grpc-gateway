'use strict';

var findup = require('findup');
var fs = require('fs');
var resolve = require('path').resolve;

function getConfigObject(filename) {
  try {
    var rcFile = findup.sync(process.cwd(), filename);
    return JSON.parse(fs.readFileSync(resolve(rcFile, filename)));
  } catch (e) {
    return null;
  }
}

function getRcConfig() {
  return getConfigObject('.vcmrc');
}

function getPackageConfig() {
  var configObject = getConfigObject('package.json');
  return configObject && configObject.config && configObject.config['validate-commit-msg'];
}

function getConfig() {
  return getRcConfig() || getPackageConfig() || {};
}

module.exports = {
  getConfig: getConfig,
  getRcConfig: getRcConfig,
  getPackageConfig: getPackageConfig,
  getConfigObject: getConfigObject
};
