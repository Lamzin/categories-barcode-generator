#!/bin/sh
docker build -t barcode-generator . && \
docker run -it --rm -v "$(pwd)"/out:/out barcode-generator
