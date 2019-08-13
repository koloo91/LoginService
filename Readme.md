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

## Environment
| Key  | Default  |
|---|---|
| JWT_KEY  |  s3cr3t |
|  DBA_USER |  kolo |
|  DBA_PASSWORD |  Pass00 |
| DB_USER  | lgn  |
|  DB_PASSWORD | lgn  |
|  DB_HOST | localhost  |
| DB_NAME  |  lgn_service |
| PORT  | 8080  |
