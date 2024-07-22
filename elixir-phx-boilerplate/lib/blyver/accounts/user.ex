defmodule Blyver.Accounts.User do
  @moduledoc false
  use Ecto.Schema

  use Pow.Ecto.Schema,
    password_min_length: 8

  use Pow.Extension.Ecto.Schema,
    extensions: [PowEmailConfirmation, PowResetPassword]

  import Ecto.Changeset
  import Pow.Ecto.Schema.Changeset, only: [confirm_password_changeset: 3]

  @timestamps_opts [type: :utc_datetime]

  schema "users" do
    field :first_name, :string
    field :last_name, :string
    field :email_confirmed, :boolean
    field :account_status, :string
    field :password_hash, :string
    field :street_address, :string
    field :city, :string
    field :phonenumber, :string

    pow_user_fields()
    timestamps()
  end

  @doc false
  def changeset(user, attrs) do
    user
    |> pow_changeset(attrs)
    |> pow_extension_changeset(attrs)
    |> validate_password()
  end

  def login_changeset(user, attrs) do
    user
    |> cast(attrs, [:email, :password])
    |> validate_required([:email, :password])
  end

  def new_reset_password_changeset(user, attrs) do
    user
    |> cast(attrs, [:email])
    |> validate_required([:email])
  end

  def submit_reset_password_changeset(user_or_changeset, attrs) do
    user_or_changeset
    |> confirm_password_changeset(attrs, @pow_config)
    |> validate_password()
  end

  defp validate_password(changeset) do
    changeset
    |> validate_required(:password)
    |> validate_format(:password, ~r/[A-Z]+/,
      message: "Password should contain at least one upper case letter"
    )
    |> validate_format(:password, ~r/[a-z]+/,
      message: "Password should contain at least one lower case letter"
    )
    |> validate_format(:password, ~r/[#\!\?&@\$%^&*\(\)]+/,
      message: "Password should contain at least one special character"
    )
  end
end
