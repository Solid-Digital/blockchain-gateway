defmodule TbgNodes.PublicEthereumNetworks.NetworkCreateForm do
  @moduledoc false

  use Ecto.Schema
  import Ecto.Changeset

  embedded_schema do
    field :name, :string
    field :network_configuration, :string
    field :deployment_type, :string
    field :archive_data, :boolean
  end

  @doc false
  def changeset(network, attrs) do
    network
    |> cast(attrs, [:name, :network_configuration, :deployment_type, :archive_data])
    |> validate_required([:name, :network_configuration, :deployment_type])
    |> validate_inclusion(:network_configuration, [
      "mainnet",
      "ropsten",
      "mainnet_archive",
      "ropsten_archive"
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
