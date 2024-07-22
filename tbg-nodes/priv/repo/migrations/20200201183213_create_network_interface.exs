defmodule TbgNodes.Repo.Migrations.CreateNetworkInterface do
  use Ecto.Migration

  def change do
    create table(:network_interfaces) do
      add :protocol, :string
      add :url, :string
      add :configuration_link, :string
      add :doc_link, :string

      add :network_id, references(:networks)

      timestamps()
    end
  end
end
