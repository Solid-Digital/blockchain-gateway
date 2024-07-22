defmodule TbgNodes.PermissionedEthereumNetworks.BesuNode do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset

  alias TbgNodes.PermissionedEthereumNetworks

  schema "permissioned_ethereum_networks_besu_nodes" do
    field :name, :string
    field :public_key, :string
    field :private_key, :string
    field :address, :string
    belongs_to :network, PermissionedEthereumNetworks.Network, foreign_key: :network_id
    has_one :network_config, through: [:network, :config]
    field :node_type, :string
    field :uuid, Ecto.UUID, autogenerate: true
    field :managed_by, :string

    timestamps()
  end

  @spec changeset(%__MODULE__{}, %PermissionedEthereumNetworks.Network{}, %{}) ::
          Ecto.Changeset.t()
  def changeset(besu_node, network, attrs \\ %{})

  def changeset(
        %__MODULE__{} = besu_node,
        %PermissionedEthereumNetworks.Network{} = network,
        attrs
      ) do
    besu_node
    |> cast(attrs, [:name, :node_type, :managed_by])
    |> validate_required([:name, :node_type])
    |> put_assoc(:network, network)
    |> validate_inclusion(:node_type, ["validator", "normal", "boot"])
    |> validate_inclusion(:managed_by, ["unchain", "external"])
    |> put_keys()
  end

  def node_attrs(node_number, :normal_node) do
    %{name: "normal-node-#{node_number}", node_type: "normal"}
  end

  def node_attrs(node_number, :boot_node) do
    %{name: "boot-node-#{node_number}", node_type: "boot"}
  end

  def node_attrs(node_number, :validator_node) do
    %{name: "validator-node-#{node_number}", node_type: "validator"}
  end

  @spec put_keys(Ecto.Changeset.t()) :: Ecto.Changeset.t()
  defp put_keys(changeset) do
    case changeset do
      %Ecto.Changeset{
        valid?: true,
        data: %{private_key: nil, public_key: nil, address: nil}
      } ->
        {:ok, private_key} = TbgNodes.ETH.generate_private_key(:hex)
        {:ok, public_key} = TbgNodes.ETH.get_public_key(:hex, private_key)
        {:ok, address} = TbgNodes.ETH.get_address(:hex, private_key)

        changeset
        |> put_change(:private_key, private_key)
        |> put_change(:public_key, public_key)
        |> put_change(:address, address |> ETH.encode16())

      _ ->
        changeset
    end
  end
end
