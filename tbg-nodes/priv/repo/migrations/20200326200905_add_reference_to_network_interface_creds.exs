defmodule TbgNodes.Repo.Migrations.AddReferenceToNetworkInterfaceCreds do
  use Ecto.Migration

  def change do
    drop_if_exists constraint(
                     :network_interface_creds,
                     "network_interface_creds_network_interface_id_fkey"
                   )

    alter table(:network_interface_creds) do
      modify :network_interface_id, references(:network_interfaces, on_delete: :delete_all)
    end
  end
end
