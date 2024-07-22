start:
	docker-compose up -d

stop:
	docker-compose rm --stop --force

deploy-whoami:
	KUBECONFIG=kubeconfig.yaml kubectl apply -f ./templates/k3s/manifests/whoami.yaml
