# Contributing

[gls](https://github.com/ozansz/gls) is a [MIT](LICENSE) licensed minimal file manager project written in Go, and accepts contributions via GitHub pull requests. This document outlines some of the conventions on development workflow, pull request message formatting, and other resources to make it easier to get your contribution accepted.

## Table of Contents
* [Getting Started](#getting-started)
* [Contribution Flow](#contribution-flow)
    + [Opening a pull request](#opening-a-pull-request)
* [Code of Conduct](#code-of-conduct)

## Getting Started

1. Fork the project [github.com/ozansz/gls](https://github.com/ozansz/gls) to your own account
2. Clone the repository on your development machine

```bash
git clone git@github.com:<your-username>/gls.git
```

3. Make sure the CLI builds flawlessly

```bash
cd gls
go build cmd/gls.go
```

## Contribution Flow

This is a outline of what a contributor's workflow looks like:

1. Select an issue to work on from [github.com/ozansz/gls/issues](https://github.com/ozansz/gls/issues), **OR** open an issue on the same page if you found another issue to fix, or have an idea to improve the project

2. Create a separate branch from `master` branch to base your work.

```bash
git checkout -b this-is-a-super-cool-feature
```

3. Work on your fix/implementation. Also please do proper commenting on your code where it may be hard for people to understand at first sight

4. Update the documentation (currently only the [README.md](README.md)) if needed

4. Push your changes to the branch you have created

5. Submit a pull request to the original repository. Please see [the section below](#opening-a-pull-request) section before opening your PR

### Opening a pull request

Please follow the below format while creating your pull request:

* **Title**: Make sure that your PR's title summarizes your contribution in a short simple sentence. Ex: "update installation section of documentation"

Note: Your PR title doesn't have to be the same as the title of the issue you're intended to fix, if you're fixing a specific issue in the PR.

* **Content**: We have a PR template which includes some related information about your implementation and fixes which make it easier for us to review your contribution. Please make sure to follow the template.

## Code of Conduct

Please check the [code of conduct](CODE_OF_CONDUCT.md) before contribution.
