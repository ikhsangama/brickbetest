.PHONY: runscript
init: # recreate resource including container, db schema, queue, and data
	rm -rf ./tmp
	go build -o main ./app
	docker-compose up -d --build
	chmod +x ./scripts/queue-source.sh
	./scripts/queue-source.sh

build:
	go build -o main ./app

run:
	./main start-rest-server

worker-1:
	./main start-transferrequest-consumer

worker-2:
	./main start-recordtxn-consumer

worker-3:
	./main start-transferstatuschecker-cron

check-queue-list:
	aws --profile localstack --endpoint-url=http://localhost:4566 sqs list-queues

check-record-transaction-q:
	aws --profile localstack --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/record-transaction
	# aws --endpoint-url=http://localhost:4566 sqs get-queue-attributes --queue-url http://localhost:4566/000000000000/record-transaction --attribute-names All

check-record-transaction-dlq:
	aws --profile localstack --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/record-transaction

check-transfer-request-q:
	aws --profile localstack --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/bank-transfer-request_dlq

check-transfer-request-dlq:
	aws --profile localstack --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/bank-transfer-request_dlq

stop:
	docker-compose down --volumes
	rm -rf ./tmp