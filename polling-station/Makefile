
image := dhrp/blockchain-polling-station:$(shell ksuid)
adapter_url := http://ec2-54-175-224-230.compute-1.amazonaws.com/adapter
auth := 248dd60d-63d8-400e-abbc-f75906a909dd

image:
	docker build -t ${image} .

push:
	docker push $(image)

remove:
	docker rm -f vote-app

run:
	docker run -p 3000:3000 -d --name vote-app $(image) 

restart: remove run

test/tostation:
	curl -d @payload_tostation.json localhost:3000/vote -H 'Content-Type: application/json' -H 'Accept: */*'

test/adapter:
	curl -d @payload_toadapter.json ${adapter_url} -H 'Authorization: Bearer ${auth}' -H 'Content-Type: application/json' -H 'Accept: */*'

deploy: 
	sed -E 's;image.*;image: ${image};' deployment.yaml > deployment-latest.yaml
	kubectl apply -f deployment-latest.yaml

all: image push deploy 
	kubectl get pods