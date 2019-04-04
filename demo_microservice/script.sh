
#!/bin/bash
export GOOS=linux
export CGO_ENABLED=0

cd books;go get;go build -o books;echo built `pwd`;cd ..

cd customers;go get;go build -o customers;echo built `pwd`;cd ..

export GOOS=darwin

docker rm demomicroservice_customers_1 -f
docker rm demomicroservice_books_1 -f
docker image rm v/customers
docker image rm v/books

docker-compose up -d
