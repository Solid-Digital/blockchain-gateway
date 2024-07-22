defmodule TbgNodes.Users.NpsFeedback do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset
  @timestamps_opts [type: :utc_datetime]

  schema "nps_feedbacks" do
    field :score, :integer

    belongs_to :user, TbgNodes.Users.User
    timestamps()
  end

  @doc false
  def changeset(nps_feedback, attrs) do
    nps_feedback
    |> cast(attrs, [:user_id, :score])
    |> validate_inclusion(:score, 0..10)
    |> validate_required([:user_id, :score])
  end
end
