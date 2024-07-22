defmodule TbgNodes.Organizations.OrganizationMember do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset

  @primary_key false
  schema "organization_members" do
    field :role, :string

    belongs_to :user, TbgNodes.Users.User
    belongs_to :organization, TbgNodes.Organizations.Organization

    timestamps()
  end

  @doc false
  def changeset(organization_member, attrs) do
    organization_member
    |> cast(attrs, [:role, :user_id, :organization_id])
    |> validate_required([:role, :user_id, :organization_id])
    |> validate_inclusion(:role, ["member", "admin"])
  end
end
