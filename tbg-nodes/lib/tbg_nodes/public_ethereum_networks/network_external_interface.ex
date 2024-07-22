defmodule TbgNodes.PublicEthereumNetworks.NetworkExternalInterface do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset
  alias TbgNodes.PublicEthereumNetworks
  @timestamps_opts [type: :utc_datetime]

  schema "public_ethereum_networks_network_external_interfaces" do
    belongs_to :network, PublicEthereumNetworks.Network, foreign_key: :network_id
    field :configuration_link, :string
    field :doc_link, :string
    field :protocol, :string
    field :url, :string

    has_many :basicauth_creds, PublicEthereumNetworks.BasicauthCred,
      foreign_key: :public_ethereum_networks_network_external_interface_id

    timestamps()
  end

  def changeset(
        %__MODULE__{} = external_interface,
        %PublicEthereumNetworks.Network{} = network,
        attrs
      ) do
    external_interface
    |> cast(attrs, [:protocol, :url, :configuration_link, :doc_link])
    |> put_assoc(:network, network)
    |> validate_inclusion(:protocol, ["http", "websocket"])
  end
end
