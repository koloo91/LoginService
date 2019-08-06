#!/usr/bin/env bash

echo "Running dba migrations"
migrate --source file://../migrations/dba -database "postgres://kolo:Pass00@localhost:5433/postgres?sslmode=disable" up
echo "Dba migrations finished"

echo "Running lgn migrations"
migrate --source file://../migrations/lgn -database "postgres://lgn:lgn@localhost:5433/lgn_service?sslmode=disable" up
echo "Lgn migrations finished"
