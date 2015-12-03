# gocompat [![Build Status](https://travis-ci.org/s2gatev/gocompat.svg?branch=master)](https://travis-ci.org/s2gatev/gocompat) [![Coverage Status](https://coveralls.io/repos/s2gatev/gocompat/badge.svg?branch=master&service=github)](https://coveralls.io/github/s2gatev/gocompat?branch=master)

Backwards compatibility checker for Go APIs.

## Introduction

**gocompat** allows you to verify backwards compatibility of you project API by continuously storing an index of all exported symbols. This index is stored in a `.gocompat` file that is used to later compare against the current API state.

## Installation

Execute `go get -u "github.com/s2gatev/gocompat"`.

## Usage

Execute `gocompat` inside your project directory.

## License

gocompat is licensed under the [MIT License](LICENSE).