defmodule TbgNodes.Repo.Migrations.UpdateAddUserFkToNetwork do
  use Ecto.Migration

  def change do
    alter table(:networks) do
      add :user_id, references("users", on_delete: :nothing)
    end
  end
end
