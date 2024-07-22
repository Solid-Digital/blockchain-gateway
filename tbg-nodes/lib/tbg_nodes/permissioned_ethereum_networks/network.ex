defmodule TbgNodes.PermissionedEthereumNetworks.Network do
  @moduledoc false

  alias TbgNodes.PermissionedEthereumNetworks

  alias TbgNodes.Users.User

  use Ecto.Schema
  import Ecto.Changeset

  schema "permissioned_ethereum_networks" do
    field :name, :string
    field :config, :map
    field :consensus, :string
    field :uuid, Ecto.UUID, autogenerate: true

    has_many :besu_nodes, PermissionedEthereumNetworks.BesuNode,
      foreign_key: :network_id,
      on_delete: :delete_all

    has_many :external_interfaces, PermissionedEthereumNetworks.ExternalInterface,
      foreign_key: :network_id,
      on_delete: :delete_all

    belongs_to :user, User, foreign_key: :user_id
    timestamps()
  end

  @spec changeset(%__MODULE__{}, %User{}, map()) :: Ecto.Changeset.t()
  def changeset(%__MODULE__{} = network, user, params \\ %{}) do
    network
    |> cast(params, [:name, :consensus, :config])
    |> validate_change(:config, fn :config, config -> validate_config(:config, config) end)
    |> put_assoc(:user, user)
    |> validate_required([:name, :consensus])
  end

  defp validate_config(:config, %{} = config) when config == %{} do
    [config: "cannot be empty map"]
  end

  defp validate_config(:config, _config) do
    []
  end
end
