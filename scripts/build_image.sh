#!/usr/bin/env bash
./increment_version.sh -m "$(cat version)" > version
#docker build -t koloooo/lgn:"$(cat version)" -t koloooo/lgn:latest -f docker/x86/Dockerfile ../.
docker build -t koloooo/lgn_arm:"$(cat version)" -t koloooo/lgn_arm:latest -f docker/arm/Dockerfile ../.

#docker push koloooo/lgn:"$(cat version)"
#docker push koloooo/lgn:latest
docker push koloooo/lgn_arm:"$(cat version)"
docker push koloooo/lgn_arm:latest
