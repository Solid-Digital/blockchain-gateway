defmodule TbgNodes.Repo.Migrations.AlterNetworksAddDelete do
  use Ecto.Migration

  def change do
    alter table(:network_interfaces) do
      modify :network_id, :bigint, on_delete: :delete_all
    end
  end
end
