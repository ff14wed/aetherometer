# Contributing to Aetherometer

## Getting Started

First, make sure you have all the dependencies for building the project. Docs for building can be located [here](docs/building.md).

If you're looking for ways to help (much appreciated), the [issue
tracker](https://github.com/ff14wed/aetherometer/issues) will have some
things to potentially work on.

## Plugin Development

Go to the documentation [here](docs/plugin_work.md).

## Reporting Bugs/Issues

To keep track of things that need fixing, any issues should be reported
through the [issue tracker](https://github.com/ff14wed/aetherometer/issues)
on GitHub (if communicated through Discord, issues could easily get lost).

## Making Changes

If making changes to `core`, you should also write tests to
cover this change. Currently, we use [Ginkgo](https://github.com/onsi/ginkgo)
which is a neat BDD framework for testing Golang.

To run all tests in `core`, cd into the `core` directory and run `ginkgo -r
-p -race` (flags are recursive, in parallel, turn on race detector).

If making changes to `ui`, there currently aren't any tests right now, but
feel free to contribute tests there too.

When you're ready to contribute the changes, please submit a pull request
and the change will be reviewed.