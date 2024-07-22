defmodule TbgNodes.PermissionedEthereumNetworks.ExternalInterface do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset

  alias TbgNodes.PermissionedEthereumNetworks

  @valid_protocols ["http", "websocket"]

  schema "permissioned_ethereum_networks_external_interfaces" do
    belongs_to :network, PermissionedEthereumNetworks.Network, foreign_key: :network_id
    field :protocol, :string
    field :url, :string
    field :uuid, Ecto.UUID, autogenerate: true
    field :target, :map

    has_many :basicauth_creds, PermissionedEthereumNetworks.BasicauthCred,
      foreign_key: :external_interface_uuid,
      on_delete: :delete_all,
      references: :uuid

    timestamps()
  end

  def changeset(external_interface, network, attrs \\ %{})

  def changeset(
        %__MODULE__{} = external_interface,
        %PermissionedEthereumNetworks.Network{} = network,
        attrs
      ) do
    external_interface
    |> cast(attrs, [:protocol, :url, :target])
    |> put_assoc(:network, network)
    |> validate_required([:protocol, :target])
    |> validate_inclusion(:protocol, @valid_protocols)
    |> validate_change(:target, &validate/2)
  end

  defp validate(:target, target) do
    case target do
      %{network_uuid: _, node_uuid: _} ->
        []

      %{network_uuid: _, node_type: _} ->
        []

      _ ->
        [
          target:
            "target missing keys: :network_uuid and one of :node_uuid or :node_type. Got #{
              Map.keys(target) |> Enum.join(",")
            }."
        ]
    end
  end
end
