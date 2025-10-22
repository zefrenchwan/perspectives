#./bin/sh
go clean -testcache
for path in `find . -type d -name '*_test'`; do
    go test $path
done