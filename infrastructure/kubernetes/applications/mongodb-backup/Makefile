build:
	docker build -t registry.unchain.io/unchain/mongodb-backup .

push:
	docker push registry.unchain.io/unchain/mongodb-backup

job:
	kubectl create -f s3_backup_job.yaml