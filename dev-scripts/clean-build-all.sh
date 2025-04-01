#!/bin/bash
rm -rf ./bin
rm -rf ./go-typer
go build -o ./bin/gotyper
./bin/gotyper start
