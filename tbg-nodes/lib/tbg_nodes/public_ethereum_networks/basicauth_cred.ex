defmodule TbgNodes.PublicEthereumNetworks.BasicauthCred do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset
  alias TbgNodes.PublicEthereumNetworks

  schema "public_ethereum_networks_basicauth_creds" do
    field :username, :string
    field :password, :string

    belongs_to :network_external_interface,
               PublicEthereumNetworks.NetworkExternalInterface,
               foreign_key: :public_ethereum_networks_network_external_interface_id

    timestamps()
  end

  @doc false
  def changeset(public_network_basicauth_creds, attrs) do
    public_network_basicauth_creds
    |> cast(attrs, [:username, :password])
    |> validate_required([:username, :password])
  end
end
