# Snippetbox
Creating and sharing snippets of text with Golang.

This is an educational project developed following the book "Let's Go" by Alex Edwards.

## Folder structure
```
.
├── cmd                 // application-specific code
│   ├── db              // SQL files
│   └── web             // web app
├── internal            // non-application-specific
│   ├── models          // DB models
│   ├── utils           // misc utils
│   └── validator       // validation tools
├── tls                 // TLS certificate files
└── ui                  // user-interface assets
    ├── html            // html templates
    └── static          // static files
```

## Requirements
### Environment Variables
Create and fill in the environment variables file
```shell
cp .env.dev ./.env
```

### TLS certificate
Store the TLS certificate files in `./tls`

To generate self-signed certificates
```shell
cd tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## Running the project
Start the DB container
```shell
docker-compose up -d
```

Run the web server
```shell
go run ./cmd/web
```

To load dummy data in the DB
```shell
cat ./cmd/db/load_dummy_data.sql | docker exec -i <container_name> mysql -usnippetbox -p<pwd> snippetbox
```