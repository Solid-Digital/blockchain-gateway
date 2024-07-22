#!/bin/bash

adapter_url='https://unchain-staging.eu2.adapter.unchain.io/dan--develop'

function transaction() {
    echo $1
    curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x5E5Af7670253D57A1fB1CbFFF19f47318E827CA6","pending"],"id":67}' https://admin:aEwvI2fYmq7z5f6e@dan.private.nodes.unchain.io/v0/node1

# domain="unchain23.org3"
domain="unchain${RANDOM}.org3"

echo $1
echo $domain

echo $1
curl --max-time 900 --location --request POST $adapter_url \
--header 'Content-Type: application/json' \
--data-raw '{
	"function": "createDomain",
	"params": {"_name":"'$domain'", "_owners": [["unchain", "100"], ["dan", "100"]], "_isPrivate": false}
}'

curl --max-time 900 --location --request POST $adapter_url \
--header 'Content-Type: application/json' \
--data-raw '{
  "function": "createAndSignPurchaseAgreement",
  "params": {
    "_name": "'$domain'",
    "_agreementID": "118743",
    "_totalAmount": 153000,
    "_vatAmount": 0,
    "_buyer": "unchain",
    "_seller": "unchain",
    "_domainShares": 100,
    "_paymentScheme": [
      [1,12800,0,1584533737,1584533800],
      [2,12800,0,1587160800,0],
      [3,12800,0,1589752800,0],
      [4,12800,0,1592431200,0],
      [5,12800,0,1595023200,0],
      [6,12800,0,1597701600,0],
      [7,12800,0,1600380000,0],
      [8,12800,0,1602972000,0],
      [9,12800,0,1605654000,0],
      [10,12800,0,1608246000,0],
      [11,12800,0,1610924400,0],
      [12,12200,0,1613602800,0]
    ]
  },
  "username": "Undeveloped"
}'

curl --max-time 900 --location --request POST $adapter_url \
--header 'Content-Type: application/json' \
--data-raw '{
	"function": "getDomain",
	"params": {"_name":"'$domain'"}
}'



    echo $1

    curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x5E5Af7670253D57A1fB1CbFFF19f47318E827CA6","pending"],"id":67}' https://admin:aEwvI2fYmq7z5f6e@dan.private.nodes.unchain.io/v0/node1

    echo '--------------'
}

for i in {1..40}
do
    transaction $i
done


# synchronousMode doesn't help with nonces because it still allows processing multiple requests at once
# but it does help with making sure the transactions are committed before responding
