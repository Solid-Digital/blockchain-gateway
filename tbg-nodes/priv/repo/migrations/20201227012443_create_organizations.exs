defmodule TbgNodes.Repo.Migrations.CreateOrganizations do
  use Ecto.Migration

  def change do
    create table(:organizations) do
      add :name, :string

      timestamps()
    end

    create unique_index(:organizations, [:name])

    create table(:organization_members, primary_key: false) do
      add :organization_id, references(:organizations, on_delete: :delete_all), primary_key: true
      add :user_id, references(:users), primary_key: true
      add :role, :string

      timestamps()
    end

    create unique_index(:organization_members, [:organization_id, :user_id])
  end
end
