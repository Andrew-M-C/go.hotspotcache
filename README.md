# go.hotspotcache

[![Build Status](https://travis-ci.org/Andrew-M-C/go.hotspotcache.svg?branch=master)](https://travis-ci.org/Andrew-M-C/go.hotspotcache)
[![GoDoc](https://godoc.org/github.com/Andrew-M-C/go.hotspotcache?status.svg)](https://godoc.org/github.com/Andrew-M-C/go.hotspotcache)
[![Coverage Status](https://coveralls.io/repos/github/Andrew-M-C/go.hotspotcache/badge.svg?branch=master)](https://coveralls.io/github/Andrew-M-C/go.hotspotcache?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/Andrew-M-C/go.hotspotcache)](https://goreportcard.com/report/github.com/Andrew-M-C/go.hotspotcache)

Package hotspotcache stores values like sync.Map does.

Meanwhile, it handles a aging queue: every access to a existing object in the cache, it will be lifted up to he top of the queue and marked as the newest element. When the size in cache exceeds preset limit size, the bottom elements of the queue would be removed.
