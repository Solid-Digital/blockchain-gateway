defmodule TbgNodes.Repo.Migrations.AlterPermissionedEthereumNetworksBesuNodesAddAddress do
  use Ecto.Migration

  def change do
    alter table(:permissioned_ethereum_networks_besu_nodes) do
      add :address, :string
    end
  end
end
