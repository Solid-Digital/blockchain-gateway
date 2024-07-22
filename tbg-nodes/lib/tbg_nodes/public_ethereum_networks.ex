defmodule TbgNodes.PublicEthereumNetworks do
  @moduledoc """
  The Networks context.
  """

  import Ecto.Query, warn: false
  alias TbgNodes.Repo

  alias TbgNodes.NetworkMonitor.Status
  alias TbgNodes.PublicEthereumNetworks
  alias TbgNodes.Utils.NameGenerator

  @spec get_status_for_user(String.t(), String.t()) :: Status.t()
  def get_status_for_user(network_uuid, user_id) do
    network =
      get_network_with_interfaces_for_user_by_uuid!(
        user_id,
        network_uuid
      )

    TbgNodes.NetworkMonitor.get_cached_status(network)
  end

  def get_network_url_config do
    Application.get_env(:tbg_nodes, TbgNodes.Networks)[:network_url_templates]
  end

  @doc """
  Returns the list of networks.

  ## Examples

      iex> list_networks()
      [%Network{}, ...]

  """
  def list_networks do
    Repo.all(PublicEthereumNetworks.Network)
    |> Repo.preload(:user)
  end

  # String comes from trusted source (internal config)
  # sobelow_skip ["RCE.EEx"]
  @spec get_network_url(map(), String.t(), String.t(), String.t(), String.t()) :: any
  def get_network_url(config, tag, protocol, network_uuid, network_name) do
    EEx.eval_string(config[tag][protocol],
      tag: tag,
      network_name: network_name,
      network_uuid: network_uuid
    )
  end

  # String comes from trusted source (internal config)
  # sobelow_skip ["RCE.EEx"]
  @spec get_liveness_url(%TbgNodes.PublicEthereumNetworks.Network{}) :: String.t()
  def get_liveness_url(network) do
    res =
      EEx.eval_string(
        PublicEthereumNetworks.get_network_url_config()[network.network_configuration][
          "liveness"
        ],
        tag: network.network_configuration,
        network_uuid: network.uuid,
        network_name: network.network_configuration
      )

    if is_binary(res) do
      res
    else
      ""
    end
  end

  # String comes from trusted source (internal config)
  # sobelow_skip ["RCE.EEx"]
  @spec get_readiness_url(%TbgNodes.PublicEthereumNetworks.Network{}) :: String.t()
  def get_readiness_url(network) do
    res =
      EEx.eval_string(
        PublicEthereumNetworks.get_network_url_config()[network.network_configuration][
          "readiness"
        ],
        tag: network.network_configuration,
        network_uuid: network.uuid,
        network_name: network.network_configuration
      )

    if is_binary(res) do
      res
    else
      ""
    end
  end

  def find_verified_network_interface(network_uuid, protocol, username, password) do
    ni =
      Repo.one(
        from c in PublicEthereumNetworks.NetworkExternalInterface,
          where: c.network_uuid == ^network_uuid and c.protocol == ^protocol
      )
      |> Repo.preload(:basicauth_creds)
      |> Repo.preload(:network)

    verified =
      Enum.find(
        ni.basicauth_creds,
        fn c -> c.username == username and c.password == password end
      )

    {ni, verified}
  end

  def list_networks_with_interfaces do
    Repo.all(PublicEthereumNetworks.Network)
    |> Repo.preload(:network_external_interfaces)
  end

  def list_networks_for_user(user) do
    q = from n in PublicEthereumNetworks.Network, where: n.user_id == ^user.id
    Repo.all(q)
  end

  def list_networks_with_interfaces_for_user(user_id) do
    user = TbgNodes.Users.get_user_by_id(user_id)
    q = from n in PublicEthereumNetworks.Network, where: n.user_id == ^user.id
    Repo.all(q) |> Repo.preload(:network_external_interfaces)
  end

  def get_network_by_uuid!(uuid), do: Repo.get_by!(PublicEthereumNetworks.Network, uuid: uuid)

  def get_network_for_user_by_uuid!(user_id, uuid) do
    Repo.get_by!(PublicEthereumNetworks.Network, uuid: uuid, user_id: user_id)
  end

  def get_network_with_interfaces_for_user_by_uuid!(user_id, uuid) do
    Repo.get_by!(PublicEthereumNetworks.Network, uuid: uuid, user_id: user_id)
    |> Repo.preload(network_external_interfaces: :basicauth_creds)
  end

  @doc """
  Creates a network.

  ## Examples

      iex> create_network(%{field: value})
      {:ok, %Network{}}

      iex> create_network(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_network(attrs \\ %{}) do
    %PublicEthereumNetworks.Network{}
    |> PublicEthereumNetworks.Network.changeset(attrs)
    |> Repo.insert(returning: [:uuid])
  end

  def create_network_with_interfaces(user_id, attrs) do
    user = TbgNodes.Users.get_user_by_id(user_id)

    network_changeset =
      %PublicEthereumNetworks.Network{}
      |> PublicEthereumNetworks.Network.changeset(attrs)
      |> Ecto.Changeset.put_assoc(:user, user)

    cond do
      length(network_changeset.errors) > 0 ->
        {:error, network_changeset}

      network_changeset.valid? ->
        Repo.transaction(fn ->
          {:ok, network} = Repo.insert(network_changeset, returning: [:uuid])

          {:ok, _} =
            create_network_external_interface(
              network,
              %{protocol: "http"}
            )

          {:ok, _} =
            create_network_external_interface(
              network,
              %{protocol: "websocket"}
            )

          # Load newly added interfaces as return value
          Repo.preload(network, :network_external_interfaces)
        end)

      # no changeset errors but not valid, return error
      true ->
        {:error, network_changeset}
    end
  end

  defp create_network_external_interface(network, attrs) do
    nei_changeset =
      Ecto.build_assoc(
        network,
        :network_external_interfaces,
        attrs
      )

    {:ok, nei} = Repo.insert(nei_changeset)

    length = 20

    name = NameGenerator.generate_name(12)

    bac_changeset =
      Ecto.build_assoc(
        nei,
        :basicauth_creds,
        %{
          username: name,
          password:
            :crypto.strong_rand_bytes(length)
            |> Base.encode32()
            |> binary_part(0, length)
            |> String.downcase()
        }
      )

    {:ok, _bac} = Repo.insert(bac_changeset)
  end

  @doc """
  Updates a network.

  ## Examples

      iex> update_network(network, %{field: new_value})
      {:ok, %Network{}}

      iex> update_network(network, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_network(%PublicEthereumNetworks.Network{} = network, attrs) do
    network
    |> PublicEthereumNetworks.Network.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Deletes a Network.

  ## Examples

      iex> delete_network(network)
      {:ok, %Network{}}

      iex> delete_network(network)
      {:error, %Ecto.Changeset{}}

  """
  @spec delete_network(Ecto.UUID.t(), atom()) :: {:ok} | {:ok, any} | {:error, Ecto.Changeset.t()}
  def delete_network(uuid, user_id) do
    network = Repo.get_by!(PublicEthereumNetworks.Network, uuid: uuid, user_id: user_id)

    Repo.delete(network)
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking network changes.

  ## Examples

      iex> change_network(network)
      %Ecto.Changeset{source: %Network{}}

  """
  def change_network(%PublicEthereumNetworks.Network{} = network) do
    PublicEthereumNetworks.Network.changeset(network, %{})
  end
end
