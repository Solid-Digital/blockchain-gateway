First database migrations need to be run for blockscout. 

Go into the blockscout repo and set env variable to:

`export MIX_ENV=prod`

`export DATABASE_URL=postgresql://blockscout:AVNS_-BnKY9FnFXsNiJ420lA@block-explorer-test-db-do-user-2339835-0.b.db.ondigitalocean.com:25061/blockscout_pool?sslmode=require`

Then run `mix ecto.migrate`

If deployment needs to be removed and restarted, the db will be out of sync and needs to be deleted. Delete `blockscout` database from database pool in digital ocean, recreate it and rerun migrations

All other information is in the manifests files