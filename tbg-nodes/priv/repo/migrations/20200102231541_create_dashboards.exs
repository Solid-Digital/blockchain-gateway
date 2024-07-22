defmodule TbgNodes.Repo.Migrations.CreateDashboards do
  use Ecto.Migration

  def change do
    create table(:dashboards) do
      add :lastUpdate, :string

      timestamps()
    end
  end
end
