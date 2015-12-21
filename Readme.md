
[![GoDoc](https://godoc.org/github.com/apex/log?status.svg)](https://godoc.org/github.com/apex/log)[![Build Status](https://semaphoreci.com/api/v1/projects/d8a8b1c0-45b0-4b89-b066-99d788d0b94c/642077/badge.svg)](https://semaphoreci.com/tj/log)

# log

Package log implements a simple structured logging API designed with few assumptions.

## About

This package is designed for centralized logging solutions such as Kinesis which require encoding and decoding before fanning-out to handlers. The API is very similar to Logrus, however does not make the same formatting assumptions which make it difficult to marshal/unmarshal an entry over the wire.

You may use this package just like Logrus, with inline handlers, however it's recommended that a centralized solution is used. This allows you to filter, add, or remove logging service providers or "sinks" without re-configuring and re-deploying dozens of applications. This is especially important when using AWS Lambda which encourages many small programs.

# License

MIT