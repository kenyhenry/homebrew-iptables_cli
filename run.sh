#!/bin/sh

go build -o iptables-cli ./src/main.go && sudo ./iptables-cli
