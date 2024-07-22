defmodule TbgNodes.Repo.Migrations.AlterNetworksAddDeploymentType do
  use Ecto.Migration

  def change do
    alter table(:networks) do
      add :deployment_type, :string
    end
  end
end
