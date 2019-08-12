#!/usr/bin/env bash
./increment_version.sh -m "$(cat version)" > version
docker build -t koloooo/lgn:"$(cat version)" -t koloooo/lgn:latest ../.
docker build -f ../arm.Dockerfile -t koloooo/lgn_arm:"$(cat version)" -t koloooo/lgn_arm:latest ../.

docker push koloooo/lgn:"$(cat version)"
docker push koloooo/lgn:latest koloooo/lgn_arm:"$(cat version)" koloooo/lgn_arm:latest
docker push koloooo/lgn_arm:"$(cat version)"
docker push koloooo/lgn_arm:latest
