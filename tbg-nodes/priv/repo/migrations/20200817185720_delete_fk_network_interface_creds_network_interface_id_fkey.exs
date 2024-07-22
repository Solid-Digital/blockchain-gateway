defmodule TbgNodes.Repo.Migrations.DeleteFkNetworkInterfaceCredsNetworkInterfaceIdFkey do
  use Ecto.Migration

  def change do
    execute "ALTER TABLE public_ethereum_networks_basicauth_creds DROP CONSTRAINT public_ethereum_networks_basicauth_creds_public_ethereum_networ"
  end
end
