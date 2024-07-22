defmodule TbgNodes.PermissionedEthereumNetworks.BasicauthCred do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset
  alias TbgNodes.PermissionedEthereumNetworks
  alias TbgNodes.Utils.NameGenerator

  schema "permissioned_ethereum_networks_basicauth_creds" do
    field :username, :string
    field :password, :string

    belongs_to :external_interface, PermissionedEthereumNetworks.ExternalInterface,
      foreign_key: :external_interface_uuid,
      type: Ecto.UUID,
      references: :uuid

    timestamps()
  end

  @spec changeset(atom(), %__MODULE__{}, map()) :: Ecto.Changeset.t()
  def changeset(
        :generate,
        %__MODULE__{} = basicauth_cred,
        %PermissionedEthereumNetworks.ExternalInterface{} = external_interface
      ) do
    name = NameGenerator.generate_name(12)

    length = 20

    password =
      :crypto.strong_rand_bytes(length)
      |> Base.encode32()
      |> binary_part(0, length)
      |> String.downcase()

    basicauth_cred
    |> cast(%{}, [])
    |> put_change(:username, name)
    |> put_change(:password, password)
    |> put_assoc(:external_interface, external_interface)
    |> validate_required([:username, :password])
  end
end
