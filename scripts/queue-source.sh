#!/bin/bash

endpoint_url="http://localhost:4566"
region="ap-southeast-1"

queues=("bank-transfer-request" "record-transaction" "bank-transfer-request_dlq" "record-transaction_dlq")

for queue in "${queues[@]}"; do
    aws --endpoint-url="${endpoint_url}" sqs create-queue --queue-name "${queue}"
    if [ $? -ne 0 ]; then
        echo "Failed to create queue: ${queue}"
        exit 1
    fi
done

url="http://localhost:4566/000000000000"

aws --endpoint-url="${endpoint_url}" sqs set-queue-attributes \
 --queue-url "${url}/bank-transfer-request" \
 --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"arn:aws:sqs:'${region}':000000000000:bank-transfer-request_dlq\",\"maxReceiveCount\":\"5\"}"}'
if [ $? -ne 0 ]; then
    echo "Failed to set RedrivePolicy for queue: bank-transfer-request"
    exit 1
fi

aws --endpoint-url="${endpoint_url}" sqs set-queue-attributes \
 --queue-url "${url}/record-transaction" \
 --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"arn:aws:sqs:'${region}':000000000000:record-transaction_dlq\",\"maxReceiveCount\":\"5\"}"}'
if [ $? -ne 0 ]; then
echo "Failed to set RedrivePolicy for queue: record-transaction"
    exit 1
fi