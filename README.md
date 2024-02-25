# Snippetbox
Creating and sharing snippets of text with Golang.

This is an educational project developed following the book "Let's Go" by Alex Edwards.

## Folder structure
```
.
├── cmd         // application-specific code
│   └── web     // web app
├── internal    // non-application-specific
└── ui          // user-interface assets
    ├── html    // html templates
    └── static  // static files
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