defmodule TbgNodes.Organizations.Organization do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset

  schema "organizations" do
    field :name, :string

    many_to_many :members, TbgNodes.Users.User,
      join_through: TbgNodes.Organizations.OrganizationMember

    timestamps()
  end

  @doc false
  def changeset(organization, attrs) do
    organization
    |> cast(attrs, [:name])
    |> validate_required([:name])
  end
end
