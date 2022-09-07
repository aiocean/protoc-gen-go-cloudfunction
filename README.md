```shell
protoc  --go_out=. --go_opt=paths=source_relative \
        --go-cloudfunction_out=. --go-cloudfunction_opt=paths=source_relative \
    example/example.proto
```
