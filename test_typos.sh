#!/bin/bash
wget -qO- https://github.com/crate-ci/typos/releases/download/v1.23.1/typos-v1.23.1-x86_64-unknown-linux-musl.tar.gz | tar xz typos
./typos
