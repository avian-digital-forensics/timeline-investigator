#!/usr/bin/env bash

oto -template api.go.plush \
	-out ../../pkg/api/api.gen.go \
	-pkg api \
	../api.def.go
gofmt -w ../../pkg/api/api.gen.go ../../pkg/api/api.gen.go
echo "generated api.gen.go"

oto -template client.go.plush \
	-out ../../tests/client/client.gen.go \
	-pkg client \
	../api.def.go
gofmt -w ../../tests/client/client.gen.go ../../tests/client/client.gen.go
echo "generated client.gen.go"

oto -template readme.md.plush \
	-out ../../pkg/services/readme.gen.md \
	../api.def.go
echo "generated readme.gen.md"