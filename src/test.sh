#!/bin/sh
# For windows users, same script, different extension. 
go clean -testcache 
go test ./models_test/ ./commons_test/