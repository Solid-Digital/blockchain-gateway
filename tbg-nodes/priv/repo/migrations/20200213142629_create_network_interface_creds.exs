defmodule TbgNodes.Repo.Migrations.AlterNetworkInterfacesAddCreds do
  use Ecto.Migration

  def change do
    create table(:network_interface_creds) do
      add :username, :string
      add :password, :string

      add :network_interface_id, references(:network_interfaces)

      timestamps()
    end

    create unique_index(:network_interface_creds, [:network_interface_id, :username])
  end
end
