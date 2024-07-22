defmodule TbgNodes.PermissionedEthereumNetworks do
  @moduledoc false
  alias TbgNodes.Repo
  alias TbgNodes.Users.User

  alias TbgNodes.NetworkMonitor.Status
  alias TbgNodes.PermissionedEthereumNetworks
  alias TbgNodes.PermissionedEthereumNetworks.InfraAPI

  import Ecto.Query

  @spec get_liveness_url(
          %TbgNodes.PermissionedEthereumNetworks.Network{}
          | %TbgNodes.PermissionedEthereumNetworks.BesuNode{},
          (() -> InfraAPI)
        ) :: String.t()
  def get_liveness_url(network_or_node, get_infra_api \\ &InfraAPI.get_infra_api!/0)

  def get_liveness_url(%TbgNodes.PermissionedEthereumNetworks.Network{} = network, get_infra_api) do
    infra_api = get_infra_api.()

    {:ok, url} = infra_api.get_liveness_url(network)
    url
  end

  def get_liveness_url(%TbgNodes.PermissionedEthereumNetworks.BesuNode{} = node, get_infra_api) do
    infra_api = get_infra_api.()

    {:ok, url} = infra_api.get_liveness_url(node)
    url
  end

  @spec get_readiness_url(
          %TbgNodes.PermissionedEthereumNetworks.Network{}
          | %TbgNodes.PermissionedEthereumNetworks.BesuNode{},
          (() -> InfraAPI)
        ) :: String.t()
  def get_readiness_url(network_or_node, get_infra_api \\ &InfraAPI.get_infra_api!/0)

  def get_readiness_url(%TbgNodes.PermissionedEthereumNetworks.Network{} = network, get_infra_api) do
    infra_api = get_infra_api.()

    {:ok, url} = infra_api.get_readiness_url(network)
    url
  end

  def get_readiness_url(%TbgNodes.PermissionedEthereumNetworks.BesuNode{} = node, get_infra_api) do
    infra_api = get_infra_api.()

    {:ok, url} = infra_api.get_readiness_url(node)
    url
  end

  @spec create_network(
          %PermissionedEthereumNetworks.NetworkUserInput{},
          integer(),
          (() -> InfraAPI)
        ) ::
          {:ok, %PermissionedEthereumNetworks.Network{}} | {:error, String.t()}
  def create_network(
        params,
        user_id,
        get_infra_api \\ &InfraAPI.get_infra_api!/0
      ) do
    infra_api = get_infra_api.()
    user = TbgNodes.Users.get_user_by_id(user_id)

    do_create_network(params, user, infra_api)
    |> TbgNodes.Repo.transaction()
    |> case do
      {:ok, %{network_db: network}} ->
        {:ok, network}

      {:error, _op, _res, %{network_db: network}} ->
        case infra_api.delete_network(network) do
          {:ok} -> {:error, "creating network failed"}
        end
    end
  end

  def do_create_network(params, user, infra_api) do
    Ecto.Multi.new()
    |> Ecto.Multi.run(:network_db, fn repo, _ ->
      attrs = %{name: params.network_name, consensus: params.consensus}
      {:ok, _network} = create_network_db(repo, user, attrs)
    end)
    |> Ecto.Multi.run(:besu_normal_nodes_db, fn repo, %{network_db: network} ->
      nodes =
        1..params.number_besu_normal_nodes
        |> Enum.map(&PermissionedEthereumNetworks.BesuNode.node_attrs(&1, :normal_node))
        |> Enum.map(
          &PermissionedEthereumNetworks.BesuNode.changeset(
            %PermissionedEthereumNetworks.BesuNode{},
            network,
            &1
          )
        )
        |> Enum.map(fn cs -> {:ok, _node} = repo.insert(cs) end)
        |> Enum.map(fn {:ok, node} -> node end)

      {:ok, nodes}
    end)
    |> Ecto.Multi.run(:besu_validator_nodes_db, fn repo, %{network_db: network} ->
      nodes =
        1..params.number_besu_validators
        |> Enum.map(&PermissionedEthereumNetworks.BesuNode.node_attrs(&1, :validator_node))
        |> Enum.map(
          &PermissionedEthereumNetworks.BesuNode.changeset(
            %PermissionedEthereumNetworks.BesuNode{},
            network,
            &1
          )
        )
        |> Enum.map(fn cs -> {:ok, _node} = repo.insert(cs) end)
        |> Enum.map(fn {:ok, node} -> node end)

      {:ok, nodes}
    end)
    |> Ecto.Multi.run(:besu_boot_nodes_db, fn repo, %{network_db: network} ->
      nodes =
        1..params.number_besu_boot_nodes
        |> Enum.map(&PermissionedEthereumNetworks.BesuNode.node_attrs(&1, :boot_node))
        |> Enum.map(
          &PermissionedEthereumNetworks.BesuNode.changeset(
            %PermissionedEthereumNetworks.BesuNode{},
            network,
            &1
          )
        )
        |> Enum.map(fn cs -> {:ok, _node} = repo.insert(cs) end)
        |> Enum.map(fn {:ok, node} -> node end)

      {:ok, nodes}
    end)
    |> Ecto.Multi.run(:add_network_config, fn repo, multi_params ->
      %{network_db: network, besu_validator_nodes_db: nodes} = multi_params

      with {:ok, config} <- create_network_config(nodes, network),
           {:ok, _network} <- add_network_config(repo, user, network, config) do
        {:ok, nil}
      end
    end)
    |> Ecto.Multi.run(:deploy_network_to_infra, fn _, %{network_db: network} ->
      case deploy_network(network, infra_api) do
        {:ok} -> {:ok, nil}
        {:error, err} -> {:error, err}
      end
    end)
    |> Ecto.Multi.run(:add_network_external_interfaces_db, fn repo, %{network_db: network} ->
      with {:ok, http_interface} <-
             add_network_external_interface_db(repo, network, "http"),
           {:ok, websocket_interface} <-
             add_network_external_interface_db(repo, network, "websocket") do
        {:ok, [http_interface, websocket_interface]}
      else
        {:error, err} -> {:error, err}
      end
    end)
    |> Ecto.Multi.run(:deploy_network_external_interfaces, fn repo, multi_params ->
      %{network_db: network, add_network_external_interfaces_db: interfaces} = multi_params

      # Deploys interfaces one by one. If deployment fails, returns the the deployment error.

      Enum.reduce_while(interfaces, nil, fn interface, _acc ->
        deploy_result = deploy_network_external_interface(repo, network, interface, infra_api)
        halt_if_interface_deploy_fails(deploy_result)
      end)
      |> case do
        {:ok} -> {:ok, nil}
        {:error, err} -> {:error, err}
      end
    end)
    |> Ecto.Multi.run(:create_basic_auth_creds_network_external_interfaces, fn repo,
                                                                               multi_params ->
      %{add_network_external_interfaces_db: interfaces} = multi_params

      interfaces
      |> Enum.each(fn interface -> {:ok} = add_basicauth_creds(repo, interface) end)

      {:ok, nil}
    end)
  end

  defp halt_if_interface_deploy_fails({:ok}), do: {:cont, {:ok}}
  defp halt_if_interface_deploy_fails({:error, err}), do: {:halt, {:error, err}}

  @spec create_network_db(Ecto.Repo.t(), %User{}, map()) ::
          {:ok, %PermissionedEthereumNetworks.Network{}} | {:error, Ecto.Changeset.t()}
  def create_network_db(repo, %User{} = user, attrs) do
    %PermissionedEthereumNetworks.Network{}
    |> PermissionedEthereumNetworks.Network.changeset(user, attrs)
    |> repo.insert()
    |> case do
      {:ok, %PermissionedEthereumNetworks.Network{} = network} -> {:ok, network}
      {:error, err} -> {:error, err}
    end
  end

  @spec create_node_db(Ecto.Repo.t(), %PermissionedEthereumNetworks.Network{}, map()) ::
          {:ok, %PermissionedEthereumNetworks.BesuNode{}} | {:error, Ecto.Changeset.t()}
  def create_node_db(repo, %PermissionedEthereumNetworks.Network{} = network, %{} = attrs) do
    %PermissionedEthereumNetworks.BesuNode{}
    |> PermissionedEthereumNetworks.BesuNode.changeset(network, attrs)
    |> repo.insert()
    |> case do
      {:ok, %PermissionedEthereumNetworks.BesuNode{} = node} -> {:ok, node}
      {:error, err} -> {:error, err}
    end
  end

  def create_network_config(
        besu_validator_nodes,
        %PermissionedEthereumNetworks.Network{} = _network
      ) do
    addresses =
      besu_validator_nodes
      |> Enum.map(fn node -> node.private_key end)
      |> Enum.map(fn private_key -> TbgNodes.ETH.get_address(:hex, private_key) end)
      |> Enum.map(fn {:ok, address} -> address end)

    genesis = TbgNodes.ETH.genesis(addresses)

    {:ok, genesis}
  end

  @spec add_network_config(Ecto.Repo.t(), %User{}, %PermissionedEthereumNetworks.Network{}, map()) ::
          {:ok, %PermissionedEthereumNetworks.Network{}} | {:error, Ecto.Changeset.t()}
  def add_network_config(
        repo,
        %User{} = user,
        %PermissionedEthereumNetworks.Network{} = network,
        %{} = config
      ) do
    network
    |> PermissionedEthereumNetworks.Network.changeset(user, %{config: config})
    |> repo.update()
    |> case do
      {:ok, updated_network} -> {:ok, updated_network}
      {:error, err} -> {:error, err}
    end
  end

  @spec deploy_network(%PermissionedEthereumNetworks.Network{}, atom()) ::
          {:ok} | {:error, String.t()}
  def deploy_network(
        %PermissionedEthereumNetworks.Network{} = network,
        infra_api,
        timeout \\ 20_000
      ) do
    query =
      from network in PermissionedEthereumNetworks.Network,
        where: network.id == ^network.id,
        select: network,
        preload: [:besu_nodes, :user]

    network_to_deploy =
      query
      |> Repo.one()

    {deploy_pid, deploy_ref} =
      spawn_monitor(fn ->
        {:ok} = infra_api.deploy_network(network_to_deploy)
        {:ok}
      end)

    receive do
      {:DOWN, ^deploy_ref, :process, ^deploy_pid, :normal} ->
        {:ok}

      {:DOWN, ^deploy_ref, :process, ^deploy_pid, _reason} ->
        {:error, "Creating k8s resources failed."}
    after
      timeout ->
        {:error, "Creating k8s resources timed-out."}
    end
  end

  @spec add_network_external_interface_db(
          Ecto.Repo.t(),
          %PermissionedEthereumNetworks.Network{},
          String.t()
        ) ::
          {:ok, %PermissionedEthereumNetworks.ExternalInterface{}} | {:error, Ecto.Changeset.t()}
  def add_network_external_interface_db(
        repo,
        %PermissionedEthereumNetworks.Network{} = network,
        protocol
      ) do
    target = %{
      network_uuid: network.uuid,
      node_type: "normal"
    }

    attrs = %{protocol: protocol, target: target}

    %PermissionedEthereumNetworks.ExternalInterface{}
    |> PermissionedEthereumNetworks.ExternalInterface.changeset(network, attrs)
    |> repo.insert()
    |> case do
      {:ok, interface} -> {:ok, interface}
      {:error, err} -> {:error, err}
    end
  end

  @spec add_basicauth_creds(Ecto.Repo.t(), %PermissionedEthereumNetworks.ExternalInterface{}) ::
          {:ok} | {:error, Ecto.Changeset.t()}
  def add_basicauth_creds(
        repo,
        %PermissionedEthereumNetworks.ExternalInterface{} = external_interface
      ) do
    PermissionedEthereumNetworks.BasicauthCred.changeset(
      :generate,
      %PermissionedEthereumNetworks.BasicauthCred{},
      external_interface
    )
    |> repo.insert()
    |> case do
      {:ok, _basicauth_cred = %PermissionedEthereumNetworks.BasicauthCred{}} -> {:ok}
      {:error, err} -> {:error, err}
    end
  end

  @spec deploy_network_external_interface(
          Ecto.Repo.t(),
          %PermissionedEthereumNetworks.Network{},
          %PermissionedEthereumNetworks.ExternalInterface{},
          atom()
        ) :: {:ok} | {:error, any()}
  def deploy_network_external_interface(
        repo,
        %PermissionedEthereumNetworks.Network{} = network,
        interface,
        infra_api
      ) do
    with {:ok, {:url, url}} <- infra_api.deploy_external_interface(interface),
         {:ok, _} <-
           repo.update(
             PermissionedEthereumNetworks.ExternalInterface.changeset(
               interface,
               network,
               %{url: url}
             )
           ) do
      {:ok}
    else
      {:error, err} -> {:error, err}
    end
  end

  @spec get_network_by_uuid!(Ecto.UUID.t()) :: %PermissionedEthereumNetworks.Network{}
  def get_network_by_uuid!(uuid) do
    query =
      from network in PermissionedEthereumNetworks.Network,
        where: network.uuid == ^uuid,
        select: network,
        preload: [:besu_nodes, external_interfaces: :basicauth_creds]

    network = Repo.one!(query)
    network
  end

  @spec get_network_for_user_by_uuid!(number(), Ecto.UUID.t()) ::
          %PermissionedEthereumNetworks.Network{}
  def get_network_for_user_by_uuid!(user_id, uuid) do
    Repo.get_by!(PermissionedEthereumNetworks.Network, uuid: uuid, user_id: user_id)
    |> Repo.preload([:besu_nodes, external_interfaces: :basicauth_creds])
  end

  @spec list_networks() :: [%PermissionedEthereumNetworks.Network{}]
  def list_networks do
    Repo.all(PermissionedEthereumNetworks.Network)
    |> Repo.preload([:external_interfaces, :user])
  end

  @spec list_networks_for_user(integer()) :: [%PermissionedEthereumNetworks.Network{}]
  def list_networks_for_user(user_id) do
    query =
      from network in PermissionedEthereumNetworks.Network,
        where: network.user_id == ^user_id

    query
    |> Repo.all()
    |> Repo.preload(:external_interfaces)
  end

  def list_nodes do
    Repo.all(PermissionedEthereumNetworks.BesuNode)
    |> Repo.preload(network: :user)
  end

  @spec get_node_status(number(), Ecto.UUID.t()) :: Status.t()
  def get_node_status(user_id, node_uuid) do
    node =
      from(node in TbgNodes.PermissionedEthereumNetworks.BesuNode,
        join: network in TbgNodes.PermissionedEthereumNetworks.Network,
        on: node.network_id == network.id,
        where: node.uuid == ^node_uuid and network.user_id == ^user_id
      )
      |> Repo.one!()
      |> Repo.preload(:network)

    TbgNodes.NetworkMonitor.get_cached_status(node)
  end

  @spec get_network_status(
          number(),
          Ecto.UUID.t()
        ) :: Status.t()
  def get_network_status(
        user_id,
        network_uuid
      ) do
    network =
      Repo.get_by!(PermissionedEthereumNetworks.Network, uuid: network_uuid, user_id: user_id)
      |> Repo.preload(:external_interfaces)

    TbgNodes.NetworkMonitor.get_cached_status(network)
  end

  @spec delete_network_for_user(
          Ecto.UUID.t(),
          number(),
          (() -> InfraAPI)
        ) :: {:ok} | {:ok, any} | {:error, String.t()}

  def delete_network_for_user(uuid, user_id, get_infra_api \\ &InfraAPI.get_infra_api!/0) do
    network = Repo.get_by!(PermissionedEthereumNetworks.Network, uuid: uuid, user_id: user_id)

    do_delete_network(network, get_infra_api)
  end

  @spec delete_network(
          Ecto.UUID.t(),
          (() -> InfraAPI)
        ) :: {:ok} | {:ok, any} | {:error, String.t()}
  def delete_network(uuid, get_infra_api \\ &InfraAPI.get_infra_api!/0) do
    network = Repo.get_by(PermissionedEthereumNetworks.Network, uuid: uuid)

    do_delete_network(network, get_infra_api)
  end

  @spec do_delete_network(
          network :: %TbgNodes.PermissionedEthereumNetworks.Network{},
          (() -> InfraAPI)
        ) :: {:ok} | {:ok, any} | {:error, String.t()}
  defp do_delete_network(network, get_infra_api) do
    infra_api = get_infra_api.()

    case infra_api.delete_network(network) do
      n when n in [{:ok}, {:error, :not_found}] ->
        Repo.delete(network)

      {:error, msg} ->
        {:error, msg}
    end
  end
end
