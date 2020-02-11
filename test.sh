#!/usr/bin/env bash

go test -coverprofile cp.out ./...
go tool cover -html=cp.out