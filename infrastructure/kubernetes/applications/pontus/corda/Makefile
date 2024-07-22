

run:
	docker run -ti -v `pwd`/node.conf:/etc/corda/node.conf corda/corda-zulu-5.0-snapshot:latest

corda-config:
	docker run -ti \
	-e MY_PUBLIC_ADDRESS="peer.corda.bc.unchain.io" \
	-e ONE_TIME_DOWNLOAD_KEY="9dd64fd6-90a6-4e8e-ba84-6b54c0593bc0" \
	-e LOCALITY="Amsterdam" -e COUNTRY="NL" \
	-v `pwd`/newconfig2:/etc/corda \
	-v `pwd`/newcerts2:/opt/corda/certificates \
	corda/corda-zulu-4.0 config-generator --testnet

	# sudo ONE_TIME_DOWNLOAD_KEY=9dd64fd6-90a6-4e8e-ba84-6b54c0593bc0 bash -c "$(curl -L https://testnet.corda.network/api/user/node/install.sh)"

run4:
	docker run -ti \
		--memory=2048m \
		--cpus=2 \
		-v `pwd`/newconfig:/etc/corda \
		-v `pwd`/newcerts:/opt/corda/certificates \
		-v `pwd`/persistence:/opt/corda/persistence \
		-v `pwd`/logs:/opt/corda/logs \
		-p 10200:10200 \
		-p 10201:10201 \
		corda/corda-zulu-4.0

		# -v /home/user/corda/samples/bank-of-corda-demo/build/nodes/BankOfCorda/cordapps:/opt/corda/cordapps \

deployment:
	kubectl apply -f deployment.yaml