#!/bin/bash

VERSION=$(cat VERSION)
docker build -t "daffaputranarendra/sysy:$VERSION" .