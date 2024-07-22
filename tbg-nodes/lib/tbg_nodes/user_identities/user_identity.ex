defmodule TbgNodes.Users.UserIdentity do
  @moduledoc false

  use Ecto.Schema
  use PowAssent.Ecto.UserIdentities.Schema, user: TbgNodes.Users.User
  @timestamps_opts [type: :utc_datetime]

  schema "user_identities" do
    pow_assent_user_identity_fields()

    timestamps()
  end
end
