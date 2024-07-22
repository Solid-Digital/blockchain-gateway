# TbgNodes

## Description

This is an Elixir Phoenix web application through which users can interact with unchain's Blockchain-as-a-Service.

## Setting up

When you first cloned this project:

  * Install dependencies with `mix deps.get`
  * Create and migrate your database with `make db-migrate`
  * Install Node.js dependencies with `cd assets && npm install`
  
To make a user admin:

   * Run `mix run -e "TbgNodes.Users.set_admin_role(~s(user@example.com))"`

## Running in development

When you want to start the project for development:
	* Run `make start`
	
Now you can visit [`localhost:4000`](http://localhost:4000) from your browser to develop the app with live-reload.
PgAdmin runs on [`localhost:49276`](http://localhost:49276) and is accessible with the credentials found in `docker-compose.yml`. 

## Deploying to staging or production

The github actions are setup to automatically deploy. When merging into develop, the new version of the application will be 
deployed on the [staging environment](https://tbg-nodes.staging.dgo.unchain.io/networks); likewise, merging to master will 
deploy to the [production environment](https://tbg.unchain.io).

### PRs and merging into master

Feature branches are rebased and merged into develop with a squash merge. Before offering your PR make sure you manually rebase it on top of develop. This way the new feature is added on top of develop as one commit.

In order to keep a tidy linear commit history on master we have to perform some manual steps because Github does not support fast-forward only merges.

Take the following steps: 
1. Create a PR from develop to master
2. Merge it through the github UI with the merge commit strategy (make sure not to delete develop in the process)
3. Immediately after merging, pull the latest develop and latest master on your machine. For both branches run `git checkout BRANCH-NAME`, pull the latest changes with `git pull --tags`.
4. Checkout develop and fast forward merge master back into it with `git merge --ff-only master`
5. Add a new version tag to the latest commit with `git tag vX.X.X`. Depending on the changes introduced update the patch, minor or major version of the latest tag that you can list with `git tag`.
6. Push your locally updated develop to github with `git push --tags`
    
After accepting the PR with the changes to master the github actions will automatically deploy to the production environment and run the migrations. Make sure to checkout [the application](tbg.unchain.io) to check that it runs OK. 

## Running phoenix application in k3s

* Make sure docker is logged into `registry.unchain.io` (ie: run `docker login registry.unchain.io`)
* Make sure all docker-compose services are running (ie: run `docker-compose up`)
* To deploy the phoenix application to k3s, run: `make quick-full-deploy-local`
* Once this is done, navigate to `http://tbg-nodes.localhost`
* If you make changes to the code and want to see them take effect on `http://tbg-nodes.localhost`,
you will need to re-run `make quick-full-deploy-local`
* If you get a server error page when first going to `http://tbg-nodes.localhost`, don't forget
you're dev database has to be migrated, this can be done with `MIX_ENV=dev mix ecto.migrate`

FYI: The phoenix application running in k3s is connected to the redis and postgres running
via docker-compose.

## Network certificates
  Kubed is used to copy the TLS certificates from letsencrypt to each tbg-nodes network's namespace. In order to enable this on a new deployment target for tbg-nodes, first add the appropriate yaml under `deployments_kustomize/<deployment-target>` (check the existing deployment targets for examples) and after deploying it, run:

  ```
    kubectl --namespace tbg-nodes-<deployment-target> annotate secret <network-host> \
      kubed.appscode.com/sync="tbg.unchain.io/network-type=permissioned,tbg.unchain.io/deployment-target=<deployment-target>"
  ```

  Where `<network-host>` is for example `permissioned.<deployment-target>.nodes.unchain.io`

## Learn more

  * Official website: http://www.phoenixframework.org/
  * Guides: https://hexdocs.pm/phoenix/overview.html
  * Docs: https://hexdocs.pm/phoenix
  * Mailing list: http://groups.google.com/group/phoenix-talk
  * Source: https://github.com/phoenixframework/phoenix
