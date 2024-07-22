defmodule TbgNodes.Repo.Migrations.AddPermissionedEthereumNetworksCascadingDeletes do
  use Ecto.Migration

  def change do
    alter table(:permissioned_ethereum_networks_basicauth_creds) do
      modify :external_interface_uuid,
             references("permissioned_ethereum_networks_external_interfaces",
               type: :uuid,
               column: :uuid,
               on_delete: :delete_all
             ),
             from:
               references("permissioned_ethereum_networks_external_interfaces",
                 type: :uuid,
                 column: :uuid,
                 on_delete: :nothing
               )
    end
  end
end
