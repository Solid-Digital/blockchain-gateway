defmodule TbgNodes.Repo.Migrations.CreateNetworks do
  use Ecto.Migration

  def change do
    create table(:networks) do
      add :name, :string
      add :protocol, :string

      timestamps()
    end
  end
end
