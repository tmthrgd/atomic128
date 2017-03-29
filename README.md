# atomic128

[![GoDoc](https://godoc.org/github.com/tmthrgd/atomic128?status.svg)](https://godoc.org/github.com/tmthrgd/atomic128)
[![Build Status](https://travis-ci.org/tmthrgd/atomic128.svg?branch=master)](https://travis-ci.org/tmthrgd/atomic128)
[![Go Report Card](https://goreportcard.com/badge/github.com/tmthrgd/atomic128)](https://goreportcard.com/report/github.com/tmthrgd/atomic128)

128-bit atomic operations using [CMPXCHG16B](http://www.felixcloutier.com/x86/CMPXCHG8B:CMPXCHG16B.html)
for Golang. **Don't use this. It is not feature complete, nor is it safe. It is strictly an experiment.**
