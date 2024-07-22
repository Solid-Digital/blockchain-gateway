defmodule TbgNodes.PublicEthereumNetworks.Network do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset
  alias TbgNodes.PublicEthereumNetworks
  @timestamps_opts [type: :utc_datetime]

  schema "public_ethereum_networks" do
    field :name, :string
    field :network_configuration, :string
    field :deployment_type, :string
    field :uuid, Ecto.UUID

    belongs_to :user, TbgNodes.Users.User

    has_many :network_external_interfaces,
             PublicEthereumNetworks.NetworkExternalInterface,
             on_delete: :delete_all,
             foreign_key: :network_id

    timestamps()
  end

  @doc false
  def changeset(network, attrs) do
    network
    |> cast(attrs, [:name, :network_configuration, :deployment_type])
    |> assoc_constraint(:user)
    |> validate_required([:name, :network_configuration, :deployment_type])
    |> validate_inclusion(:network_configuration, [
      "mainnet",
      "ropsten",
      "mainnet-archive",
      "ropsten-archive"
    ])
    |> validate_inclusion(:deployment_type, ["shared"])
  end

  def validate_fields(changeset, fields) do
    if Enum.all?(fields, &present?(changeset, &1)) do
      changeset
    else
      # Add the error to the first field only since Ecto requires a field name for each error.
      add_error(
        changeset,
        hd(fields),
        "All of these fields must be present: Network Name, Protocol, Network Configuration, and Deployment type"
      )
    end
  end

  def present?(changeset, field) do
    value = get_field(changeset, field)
    value && value != ""
  end

  def changeset_new(network, attrs) do
    network
    |> cast(attrs, [:name, :network_configuration, :deployment_type])
    |> validate_fields([:name, :network_configuration, :deployment_type])
  end
end
