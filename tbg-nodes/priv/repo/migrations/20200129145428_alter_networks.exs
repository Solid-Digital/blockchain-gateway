defmodule TbgNodes.Repo.Migrations.AlterNetworks do
  use Ecto.Migration

  def change do
    alter table(:networks) do
      add :network_configuration, :string
    end
  end
end
