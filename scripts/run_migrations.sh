#!/usr/bin/env bash

echo "Running dba migrations"
migrate --source file://../migrations/dba -database "postgres://kolo:Pass00@localhost:5432/postgres?sslmode=disable" up
echo "Dba migrations finished"

echo "Running lgn migrations"
migrate --source file://../migrations/lgn -database "postgres://lgn:lgn@localhost:5432/lgn_service?sslmode=disable" up
echo "Lgn migrations finished"
