# lgn

## Update
To update current version run:
``` bash
./increment_version.sh -m $(cat version) > version
```

## Database migrations
``` bash
migrate --source file://migrations/lgn -database "postgres://lgn:lgn@localhost:5433/lgn_service?sslmode=disable" up
```
