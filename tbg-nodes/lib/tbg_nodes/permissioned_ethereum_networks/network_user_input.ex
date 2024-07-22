defmodule TbgNodes.PermissionedEthereumNetworks.NetworkUserInput do
  @moduledoc false
  use Ecto.Schema
  import Ecto.Changeset

  @max_besu_validators 5
  @min_besu_validators 1
  @max_besu_normal_nodes 5
  @min_besu_normal_nodes 1
  @min_besu_boot_nodes 0

  embedded_schema do
    field :network_name, :string
    field :number_besu_validators, :integer
    field :number_besu_normal_nodes, :integer
    field :number_besu_boot_nodes, :integer
    field :deployment_option, :string
    field :consensus, :string
    timestamps()
  end

  @spec changeset(map(), map()) :: Ecto.Changeset.t()
  def changeset(struct, params) do
    struct
    |> cast(params, [
      :network_name,
      :number_besu_validators,
      :number_besu_normal_nodes,
      :number_besu_boot_nodes,
      :deployment_option,
      :consensus
    ])
    |> validate_required(
      [
        :number_besu_validators,
        :number_besu_normal_nodes,
        :number_besu_boot_nodes,
        :deployment_option,
        :consensus
      ],
      message: "This field can't be blank."
    )
    |> validate_required(:network_name, message: "Please give your network a name.")
    |> validate_number(:number_besu_validators,
      less_than_or_equal_to: @max_besu_validators,
      message: "You can not have more than #{@max_besu_validators} validators."
    )
    |> validate_number(:number_besu_validators,
      greater_than_or_equal_to: @min_besu_validators,
      message: "Negatives are not possible."
    )
    |> validate_number(:number_besu_normal_nodes,
      greater_than_or_equal_to: @min_besu_normal_nodes,
      message: "You need at least #{@min_besu_normal_nodes} node."
    )
    |> validate_number(:number_besu_normal_nodes,
      less_than_or_equal_to: @max_besu_normal_nodes,
      message: "You can not have more than #{@max_besu_normal_nodes} nodes"
    )
    |> validate_number(:number_besu_boot_nodes, greater_than: @min_besu_boot_nodes)
    |> validate_inclusion(:deployment_option, ["cloud"])
    |> validate_inclusion(:consensus, ["IBFT"])
  end
end
