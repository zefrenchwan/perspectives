#!/bin/sh
echo -n code lines counter '   '
cat commons/*.go | wc -l 
echo -n test lines counter '   '
cat commons_test/*.go | wc -l 
echo
echo launch tests
go clean -testcache 
go test ./commons_test/
echo