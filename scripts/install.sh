#!/bin/bash

# Author Lei Wang, yiak.wy@gmail.com
# Date: 2017/12/12
# Update: 2018/9/19
#	  2018/12/26, add Ubuntu support

set -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd  )"
OS=
COLOR_OFF='\033[0m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'

function info() {
	(echo >&2 -e "${GREEN} [INFO] ${COLOR_OFF} $*")
}

function err() {
	(echo >&2 -e "${RED} [ERROR] ${COLOR_OFF} $*" )
}

function warn() {
	(echo >&2 -e "${YELLOW} [WARNING] ${COLOR_OFF} $*" )
}

# install npm, yarn using macOS brew
function OSX() {
	if test ! $(which brew); then
		warn "Homebrew not installed."
		info "Installing Homebrew ..."
		/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
	fi
	info "update brew ..."
	# brew update
	if brew ls --versions node > /dev/null; then
		info "node installed."
	else
		warn "node not installed. Installing ..."
		brew install node
		
	fi

	init_node
	(install_node_modules_on_MacOSX) || {
		err "Failing to install dependancies on OSX"
		mv /usr/local/bin/libtool.bak /usr/local/bin/libtool
	}
	# insall_pcl

  # install gitbook plugins

  brew install graphviz
  # instal PlantUML

  # mkdir for uml
  mkdir -p ${ROOT}/pedia/assets/images/uml
}

function Ubuntu() {
if dpkg-query -l nodejs > /dev/null; then
  info "node installed"
else
  warn "node not installed. installing ..."
  # https://github.com/nodesource/distributions#debinstall
  curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
  sudo apt-get install -y nodejs
fi

sudo aptitude install -y \
	libgflags-dev \
	libgoogle-glog-dev

init_node
install_node_modules
}

Normal_opt=--save
Dev_opt=--save-dev 
Global_opt=-g

INSTALLER=npm # or yarn

node_modules=(
	"cross-env@5.1.4"
	# commit utils
	"commitizen"
	"conventional-changelog-cli"
	# update json
	"json"
  "gitbook-cli"
  # gitbook plugins
  "gitbook-pdf"
)

node_dev_modules=(
	@commitlint/{config-conventional,cli}
	"husky"
	"validate-commit-msg"
	"cz-conventional-changelog"

)

brew_global_libraries=(
	""
)

function init_node() {
	if [ ! -f 'package.json' ]; then
		$INSTALLER init --yes
	fi
}

function configGitCommitter() {
	info "init committing message format."
	commitizen init --force cz-conventional-changelog --save --save-exact
	info "init (force) changelog."
	if [ ${OS} = "OSX" ]; then
	~/.npm-packages/bin/conventional-changelog -p angular -i CHANGELOG.md -s
	fi
		
	if [ ${OS} = "Ubuntu" ]; then
	conventional-changelog -p angular -i CHANGELOG.md -s
	fi
	echo "
module.exports = {
	extends: ['@commitlint/config-conventional']
}
" > commitlint.config.js
	info "editing package.json"
	edit_json 'this.husky={"hooks": {"commit-msg": "commitlint -E HUSKY_GIT_PARAMS"}}'
	edit_json 'this.config["validate-commit-msg"]={
	"types": [
    "feat",
    "fix",
    "docs",
    "style",
    "refactor",
    "perf",
    "test",
    "chore",
    "revert",
    "ci",
    "build"
  ],
  "warnOnFail": false,
  "maxSubjectLength": 100,
  "subjectPattern": "^[A-Za-z]+-[0-9]+((,| - ).*)?",
  "subjectPatternErrorMsg": "Wrong commit format! Subject does not match subject pattern!",
  "helpMessage": ""}'
}

function edit_json() {
	# in-place editing, see https://github.com/trentm/json#readme
	json -I -f package.json -e "$1"
}

function npm_package_is_installed() {
  # set to 1 initially
  local ret=1
  # set to 0 if not found
  local pkg=$(echo $1 | cut -d/ -f 1)

  ls node_modules | grep "$pkg" >/dev/null 2>&1||
  npm list -g --depth=0 | grep "$pkg" >/dev/null 2>&1 || 
  { local ret=0; }
  # return value
  echo "$ret"
}

function install_node_modules() {
	npm config set registry="https://registry.npmjs.org"
	if [ ! -d "${ROOT}/node_modules" ]; then
		mkdir "${ROOT}/node_modules"
	fi
	info "Installing global node dependancies..." 
	for plugin in ${node_modules[@]}; do 
		if [ $(npm_package_is_installed $plugin) == 1 ]; then
		echo "$plugin installed, skipping ..."
		else
		  if [ ${OS} = "Ubuntu" ]; then
		  sudo $INSTALLER install $Global_opt $plugin
		  else
	          $INSTALLER install $Global_opt $plugin
		  fi	
		fi 
	done 

	info "Installing local node dev dependancies"
	for plugin in ${node_dev_modules[@]}; do 
		if [ $(npm_package_is_installed $plugin) == 1 ]; then
		echo "$plugin installed, skipping ..."
		else
		$INSTALLER install $Dev_opt $plugin 
		fi 
	done 
	info "done."
}

function install_node_modules_on_MacOSX() {
	set -e
	# dealing with mac libtool conflictions
	# see discussion https://github.com/barrysteyn/node-scrypt/issues/113
	mv /usr/local/bin/libtool /usr/local/bin/libtool.bak
	$INSTALLER rebuild node-sass
	install_node_modules
	mv /usr/local/bin/libtool.bak /usr/local/bin/libtool

}

function install_pcl() {
# http://www.pointclouds.org/documentation/tutorials/installing_homebrew.php
	brew install pcl
# https://github.com/totakke/homebrew-openni, homebrew/science has been discarded.
	brew tap brewsci/science
	brew tap toakke/openni
	brew install openni
	brew install sensor
	brew install sensor-kinect
	brew install nite

	brew tap toakke/openni2
	brew install openni2
	brew install openni2-freenectdriver

}

main() {
	if [ $(uname -s) == "Darwin" ]; then
		OS="OSX"
		info "<$(uname -s)> detected. checking Max OSX versions: Major.Minor ..."
		info "installing ..."
		OSX

		configGitCommitter

	else
		if [ -f /etc/lsb-release ]; then
		OS="Ubuntu"
		Ubuntu

		configGitCommitter
	
		else
		warn "platform <$(uname -s)> Not supported yet. Please install the package dependencies manually. Pull requests are welcome!"
		fi
	fi 

}

main
