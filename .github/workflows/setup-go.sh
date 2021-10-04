#!/bin/bash
set -exo pipefail

git clone https://github.com/golang/go.git /opt/go
cd /opt/go

set +e
git checkout go1.18
set -e

cd src
./make.bash
ln -vsf /opt/go/bin/* /usr/local/bin
