defmodule Blyver.Repo.Migrations.CreateUsers do
  use Ecto.Migration

  def change do
    create table(:users) do
      add :first_name, :string
      add :last_name, :string
      add :email, :string, null: false
      add :email_confirmed, :boolean
      add :account_status, :string
      add :password_hash, :string
      add :street_address, :string
      add :city, :string
      add :phonenumber, :string

      timestamps()
    end

    create unique_index(:users, [:email])
  end
end
