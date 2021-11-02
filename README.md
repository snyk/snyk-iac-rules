# `snyk-iac-rules` SDK
---

[![CircleCI](https://circleci.com/gh/snyk/snyk-iac-rules/tree/develop.svg?style=svg&circle-token=5597b9f0189554f754f38400cbe9d8f8b334c72a)](https://circleci.com/gh/snyk/snyk-iac-rules/tree/develop) 
[![Shellspec Tests](https://github.com/snyk/snyk-iac-rules/actions/workflows/main.yml/badge.svg)](https://github.com/snyk/snyk-iac-rules/actions/workflows/main.yml)
[![Golang Version](https://img.shields.io/github/go-mod/go-version/snyk/snyk-iac-rules)](https://github.com/snyk/snyk-iac-rules)

`snyk-iac-rules` is a Golang SDK that provides flags for writing, debugging, testing, and bundling custom rules for the [Snyk IaC CLI](https://github.com/snyk/snyk/).

---

# About
The SDK is a tool for writing, debugging, testing, and bundling custom rules for [Snyk Infrastructure as Code](https://snyk.io/product/infrastructure-as-code-security/). See our [Custom Rules documentation](https://docs.snyk.io/products/snyk-infrastructure-as-code/custom-rules) to learn more.

![system overview](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.github.com/snyk/snyk-iac-rules/main/assets/overview-activity-swimlanes.puml)

# Install
The SDK can be installed through multiple channels.

## Install with npm or Yarn

[snyk-iac-rules available as an npm package](https://www.npmjs.com/package/snyk-iac-rules). If you have Node.js installed locally, you can install it by running:

```bash
npm install snyk-iac-rules@latest -g
```

or if you are using Yarn:

```bash
yarn global add snyk-iac-rules
```
## More installation methods

<details>
  <summary>Standalone executables (macOS, Linux, Windows)</summary>

### Standalone executables

Use [GitHub Releases](https://github.com/snyk/snyk-iac-rules/releases) to download a standalone executable of Snyk CLI for your platform.

For example, to download and run the latest SDK on macOS, you could run:

```bash
wget https://github.com/snyk/snyk-iac-rules/releases/download/v0.1.0/snyk-iac-rules_0.1.0_Darwin_x86_64.tar.gz
chmod +x ./snyk-iac-rules
mv ./snyk-iac-rules /usr/local/bin/
```

Drawback of this method is, that you will have to manually keep the SDK up to date.

</details>

<details>
  <summary>Install with Homebrew (macOS, Linux)</summary>

### Homebrew

Install the SDK from [Snyk tap](https://github.com/snyk/homebrew-tap) with [Homebrew](https://brew.sh) by running:

```bash
brew tap snyk/tap
brew install snyk-iac-rules
```

</details>

<details>
  <summary>Scoop (Windows)</summary>

### Scoop

Install the SDK from our [Snyk bucket](https://github.com/snyk/scoop-snyk) with [Scoop](https://scoop.sh) on Windows:

```
scoop bucket add snyk https://github.com/snyk/scoop-snyk
scoop install snyk-iac-rules
```

</details>

---

# Getting started with snyk-iac-rules

Once you installed the `snyk-iac-rules` SDK, you can verify it's working by running

```bash
snyk-iac-rules --help
```

For more help, read the documentation about [Snyk Infrastructure as Code](https://docs.snyk.io/snyk-infrastructure-as-code).

# Getting support

We recommend reaching out via the [support@snyk.io](mailto:support@snyk.io) email whenever you need help with the SDK or Snyk in general.


* See [DEVELOPMENT.md](DEVELOPMENT.md) for how to setup the environment, add a new command, run the code locally, and run the tests.
* See [RELEASE.md](RELEASE.md) for how to release a new version of the SDK.

---

# Contributing

This project is open source but we don't encourage outside contributors.
