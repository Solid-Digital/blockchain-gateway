defmodule TbgNodes.Repo.Migrations.AlterNetworksAddUuid do
  use Ecto.Migration

  def change do
    execute "create extension if not exists \"uuid-ossp\"",
            "drop extension if exists \"uuid-ossp\""

    alter table(:networks) do
      add :uuid, :uuid, null: false, default: fragment("uuid_generate_v4()")
    end

    create unique_index(:networks, [:uuid])
  end
end
