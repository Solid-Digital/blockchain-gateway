defmodule TbgNodes.Repo.Migrations.CreateNpsFeedback do
  use Ecto.Migration

  def change do
    create table(:nps_feedbacks) do
      add :score, :integer
      add :user_id, references("users", on_delete: :nothing)

      timestamps()
    end
  end
end
