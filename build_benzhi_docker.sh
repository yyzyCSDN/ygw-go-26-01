#!/bin/sh
set -eu
docker build -f benzhi.Dockerfile -t tracelink:latest .
docker run --rm tracelink:latest
