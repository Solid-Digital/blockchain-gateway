defmodule TbgNodes.Users.User do
  @moduledoc false

  use Ecto.Schema
  @timestamps_opts [type: :utc_datetime]

  use Pow.Ecto.Schema,
    user_id_field: :email,
    password_min_length: 12,
    password_max_length: 4096

  use Pow.Extension.Ecto.Schema,
    extensions: [PowResetPassword]

  use PowAssent.Ecto.Schema

  import Pow.Ecto.Schema.Changeset, only: [new_password_changeset: 3]

  schema "users" do
    has_many :user_identities, TbgNodes.Users.UserIdentity,
      on_delete: :delete_all,
      foreign_key: :user_id

    many_to_many :organization_memberships, TbgNodes.Organizations.Organization,
      join_through: TbgNodes.Organizations.OrganizationMember

    field :role, :string, default: "user"
    field :username, :string

    has_many :permissioned_ethereum_networks, TbgNodes.PermissionedEthereumNetworks.Network,
      foreign_key: :user_id

    has_many :public_ethereum_networks, TbgNodes.PublicEthereumNetworks.Network,
      foreign_key: :user_id

    pow_user_fields()

    timestamps()
  end

  def changeset(user_or_changeset, attrs) do
    user_or_changeset
    |> pow_user_id_field_changeset(attrs)
    |> pow_current_password_changeset(attrs)
    |> new_password_changeset(attrs, @pow_config)
    |> pow_extension_changeset(attrs)
    |> Ecto.Changeset.cast(attrs, [:role, :email, :password_hash])
    |> Ecto.Changeset.unique_constraint(:email)
    |> Ecto.Changeset.validate_format(
      :email,
      ~r/^[A-Za-z0-9\._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}$/
    )
    |> Ecto.Changeset.validate_inclusion(:role, ~w(user admin))
  end

  # changeset for 3rd party auth users like github
  def user_identity_changeset(user_or_changeset, user_identity, attrs, user_id_attrs) do
    user_or_changeset
    |> Ecto.Changeset.cast(attrs, [:username])
    |> pow_assent_user_identity_changeset(user_identity, attrs, user_id_attrs)
  end

  def changeset_edit(user) do
    Ecto.Changeset.change(user)
  end

  def changeset_update_email(user, %{"user" => attrs}) do
    if attrs["email"] == user.email do
      user
      |> Ecto.Changeset.cast(attrs, [:email, :password])
      |> Ecto.Changeset.add_error(:email, "Email address not changed")
    else
      Pow.Ecto.Schema.Password.pbkdf2_verify(attrs["password"], user.password_hash, [])
      |> case do
        true ->
          user
          |> Ecto.Changeset.cast(attrs, [:email, :password])
          |> Ecto.Changeset.validate_required([:email, :password])
          |> Ecto.Changeset.unique_constraint(:email)
          |> Ecto.Changeset.validate_format(
            :email,
            ~r/^[A-Za-z0-9\._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}$/
          )

        false ->
          user
          |> Ecto.Changeset.cast(attrs, [:email, :password])
          |> Ecto.Changeset.add_error(:password, "Incorrect password")
      end
    end
  end

  def changeset_change_password(user, %{"user" => attrs}) do
    user
    |> Ecto.Changeset.cast(attrs, [:password, :current_password])
    |> Ecto.Changeset.validate_required([:password, :current_password])
    |> new_password_changeset(attrs, @pow_config)
  end

  def changeset_change_password_without_current_password(user, %{"user" => attrs}) do
    user
    |> Ecto.Changeset.cast(attrs, [:password])
    |> Ecto.Changeset.validate_required([:password])
    |> new_password_changeset(attrs, @pow_config)
  end

  def changeset_new_login(user_or_changeset, attrs) do
    user_or_changeset
    |> Ecto.Changeset.cast(attrs, [])
  end

  def changeset_email_validation(user_or_changeset, attrs) do
    user_or_changeset
    |> Ecto.Changeset.cast(attrs, [:email])
    |> Ecto.Changeset.unique_constraint(:email)
    |> Ecto.Changeset.validate_required([:email])
    |> Ecto.Changeset.validate_format(:email, ~r/@/)
  end

  def changeset_login_attempt(user_or_changeset, attrs) do
    user_or_changeset
    |> Ecto.Changeset.cast(attrs, [:email, :password])
    |> Ecto.Changeset.add_error(:password, "Invalid credentials")
  end

  @spec changeset_role(Ecto.Schema.t() | Ecto.Changeset.t(), map()) :: Ecto.Changeset.t()
  def changeset_role(user_or_changeset, attrs) do
    user_or_changeset
    |> Ecto.Changeset.cast(attrs, [:role])
    |> Ecto.Changeset.validate_inclusion(:role, ~w(user admin))
  end
end
