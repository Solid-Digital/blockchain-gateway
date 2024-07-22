defmodule TbgNodes.Repo.Migrations.DeleteDashboard do
  use Ecto.Migration

  def change do
    drop table(:dashboards)
  end
end
