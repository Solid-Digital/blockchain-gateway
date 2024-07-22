
## Flow of the process:
- Create private/public keys for the validators & update the secrets/validator-keys-secret.yaml with the validator private keys
- Update the configmap/configmap.yml with the public keys & genesis file
- Update the number of nodes you would like in deployments/node-deployment.yaml
- Run kubectl

## Overview of Setup
![Image ibft](../../images/ibft.png)

## NOTE:
1. validators1 and 2 serve as bootnodes as well. Adjust according to your needs
2. If you add more validators in past the initial setup, they need to be voted in to be validators i.e they will serve as normal nodes and not validators until they've been voted in.

#### 1. Boot nodes private keys
Create private/public keys for the validators using nodejs script. The private keys are put into secrets and the public keys go into a configmap to get the bootnode enode address easily
Repeat this process for as many validators as you would like to provision i.e keys and replicate the deployment & service

```bash
npm install besu-clique-config
besu-clique-config --config-file=cliqueSetup/cliqueConfigFile.json --to=cliqueSetup/networkFiles --private-key-file-name=key
sudo chown -R $USER:$USER ./cliqueSetup
mv ./cliqueSetup/networkFiles/genesis.json ./cliqueSetup/
```

Update the secrets/validator-key-secret.yaml with the private keys. The private keys are put into secrets and the public keys go into a configmap that other nodes use to create the enode address
Update the configmap/configmap.yaml with the public keys
**Note:** Please remove the '0x' prefix of the public keys

#### 2. Genesis.json
Copy the genesis.json file and copy its contents into the configmap/configmap as shown

#### 3. Update any more config if required
eg: To alter the number of nodes on the network, alter the `replicas: 2` in the deployments/node-deployments.yaml to suit

#### 4. Deploy:
```bash

./deploy.sh

```

## Add new validators:

1. Generate public and private key with above steps
2. add the private key in vault.tools.unchain.io
3. `kubectl apply` the secret with the private key (as shown in `deploy.sh`)
4. Create new stateful set copying the existing validator manifests.
5. Create new service for stateful set
6. apply the stateful set, copy the address of the validator (will show up in the logs next to `DefaultP2PNetwork | Node address: ADDRESS`)
7. call the `clique_propose` method with the address in the params with the RPC endpoint of more than 50% of the validators.
8. confirm with the RPC method `clique_getSigners` and the address of the new validator should show up


#### 5. In the dashboard, you will see each bootnode deployment & service, nodes & a node service, miner if enabled, secrets(opaque) and a configmap

If using minikube
```bash
minikube dashboard &
```

#### 6. Verify that the nodes are communicating:
```bash
minikube ssh

# once in the terminal
curl -X POST --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' <besu_NODE_SERVICE_HOST>:8545

# which should return:
The result confirms that the node running the JSON-RPC service has two peers:
{
  "jsonrpc" : "2.0",
  "id" : 1,
  "result" : "0x5"
}

```


#### 8. Delete
```
./remove.sh
```
