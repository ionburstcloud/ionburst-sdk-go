#!/bin/bash

printf "[+] - upgrade dependencies\n"
go get -u -v
printf "[+] - update mod file\n"
go mod tidy
