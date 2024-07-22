defmodule TbgNodes.Repo.Migrations.AddPermissionedBesuContext do
  use Ecto.Migration

  def change do
    create table(:permissioned_ethereum_networks) do
      add :name, :string
      add :config, :map
      add :uuid, :uuid, null: false
      add :consensus, :string
      add :user_id, references("users", on_delete: :nothing)
      timestamps()
    end

    create unique_index(:permissioned_ethereum_networks, [:uuid])

    create table(:permissioned_ethereum_networks_besu_nodes) do
      add :name, :string
      add :public_key, :string
      add :private_key, :string
      add :node_type, :string
      add :network_id, references("permissioned_ethereum_networks", on_delete: :nothing)
      add :uuid, :uuid, null: false
      add :managed_by, :string
      timestamps()
    end

    create unique_index(:permissioned_ethereum_networks_besu_nodes, [:uuid])

    create table(:permissioned_ethereum_networks_external_interfaces) do
      add :network_id, references("permissioned_ethereum_networks", on_delete: :nothing)
      add :protocol, :string
      add :url, :string
      add :uuid, :uuid, null: false
      add :target, :map
      timestamps()
    end

    create unique_index(:permissioned_ethereum_networks_external_interfaces, :uuid)

    create index(:permissioned_ethereum_networks_external_interfaces, [
             :uuid,
             :protocol,
             :network_id
           ])

    create table(:permissioned_ethereum_networks_basicauth_creds) do
      add :username, :string
      add :password, :string

      add :external_interface_uuid,
          references("permissioned_ethereum_networks_external_interfaces",
            type: :uuid,
            column: :uuid,
            on_delete: :nothing
          )

      timestamps()
    end

    create unique_index(:permissioned_ethereum_networks_basicauth_creds, [:username, :password])

    create unique_index(:permissioned_ethereum_networks_basicauth_creds, [
             :external_interface_uuid
           ])
  end
end
