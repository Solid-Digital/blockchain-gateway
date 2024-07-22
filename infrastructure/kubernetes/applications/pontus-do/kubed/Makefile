init: helm-add-repo

helm-add-repo:
	helm repo add appscode https://charts.appscode.com/stable/
	helm repo update

deploy:
	helm upgrade \
		--install kubed appscode/kubed \
		--version v0.12.0 \
		--namespace kube-system

check:
	helm upgrade \
		--install kubed appscode/kubed \
		--version v0.12.0 \
		--namespace kube-system \
		--debug \
		--dry-run

