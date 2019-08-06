#!/usr/bin/env bash
./increment_version.sh -m $(cat version) > version
docker build -t koloooo/lgn:$(cat version) -t koloooo/lgn:latest ../.
