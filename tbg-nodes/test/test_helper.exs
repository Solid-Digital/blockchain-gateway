ExUnit.start()
Ecto.Adapters.SQL.Sandbox.mode(TbgNodes.Repo, :manual)

defmodule TbgNodes.TestHelpers do
  import Ecto.Repo
  import Phoenix.LiveViewTest

  import ExUnit.Assertions

  @valid_network %{
    name: "test-network-name",
    consensus: "clique"
  }

  @valid_node %{
    name: "i am a node",
    node_type: "validator"
  }

  def create_user(%{}) do
    {:ok, user: user_fixture()}
  end

  def count_users(%{}) do
    {:ok, user_count: TbgNodes.Users.list_users() |> length()}
  end

  def create_admin(%{}) do
    {:ok, user: admin_fixture()}
  end

  def auth_user(%{conn: conn, user: user}) do
    {:ok, conn: Pow.Plug.assign_current_user(conn, user, otp_app: :tbg_nodes)}
  end

  def create_github_user(%{}) do
    user = user_fixture()
    {:ok, user_identity} = github_user_identity_fixture(user)
    {:ok, user: %TbgNodes.Users.User{user | user_identities: [user_identity]}}
  end

  def create_public_ethereum_network_with_interfaces(%{user: user}) do
    {:ok, public_network: public_ethereum_network_with_interfaces_fixture(user)}
  end

  def create_permissioned_ethereum_network(%{user: user}) do
    {:ok, permissioned_network: permissioned_ethereum_network_fixture(user)}
  end

  def create_permissioned_ethereum_network_with_interfaces(%{user: user}) do
    {:ok, permissioned_network: permissioned_network_with_interfaces_fixture(user)}
  end

  def with_network_monitor(%{}) do
    {:ok, _} = Singleton.start_child(TbgNodes.NetworkMonitor, [1], {:network_monitor, 1})

    ExUnit.Callbacks.on_exit(fn ->
      :ok = Singleton.stop_child(TbgNodes.NetworkMonitor, [1])
    end)
  end

  def create_permissioned_node(%{permissioned_network: permissioned_network}) do
    {:ok, permissioned_node: besu_node_fixture(permissioned_network)}
  end

  def assert_until(fun, opts \\ []) do
    assert TbgNodes.LiveTestHelpers.check_until(fun, opts)
  end

  def user_fixture(attrs \\ %{}) do
    TbgNodes.LiveTestHelpers.user_fixture(attrs)
  end

  def admin_fixture(attrs \\ %{}) do
    {:ok, user} =
      attrs
      |> Enum.into(%{
        email: "user#{System.unique_integer([:positive])}@test.com",
        password: "supersecret_password"
      })
      |> Pow.Operations.create(otp_app: :tbg_nodes)

    {:ok, user} =
      user.email
      |> TbgNodes.Users.set_admin_role(TbgNodes.Repo)

    user
  end

  def github_user_identity_fixture(user \\ %{}) do
    user_identity_params =
      Enum.into(%{}, %{
        provider: "github",
        uid: "random-uid",
        user_id: user.id
      })

    user_identity =
      PowAssent.Ecto.UserIdentities.Schema.changeset(
        %TbgNodes.Users.UserIdentity{},
        user_identity_params,
        %{}
      )

    {:ok, _user_identity} = TbgNodes.Repo.insert(user_identity)
  end

  def public_ethereum_network_with_interfaces_fixture(%TbgNodes.Users.User{} = user, attrs \\ %{}) do
    attrs =
      Enum.into(attrs, %{
        name: "name",
        network_configuration: "mainnet",
        deployment_type: "shared"
      })

    {:ok, public_network} =
      TbgNodes.PublicEthereumNetworks.create_network_with_interfaces(user.id, attrs)

    public_network
  end

  def permissioned_ethereum_network_user_input(user_input) do
    TbgNodes.LiveTestHelpers.permissioned_ethereum_network_user_input(user_input)
  end

  def permissioned_ethereum_network_fixture(%TbgNodes.Users.User{} = user) do
    {:ok, network} =
      TbgNodes.PermissionedEthereumNetworks.create_network_db(TbgNodes.Repo, user, @valid_network)

    network
  end

  def network_external_interface_fixture(
        %TbgNodes.PermissionedEthereumNetworks.Network{} = network,
        type \\ "http"
      ) do
    scheme =
      case type do
        "http" -> "https://"
        "websocket" -> "wss://"
      end

    {:ok, interface = %TbgNodes.PermissionedEthereumNetworks.ExternalInterface{}} =
      %TbgNodes.PermissionedEthereumNetworks.ExternalInterface{}
      |> TbgNodes.PermissionedEthereumNetworks.ExternalInterface.changeset(network, %{
        url: scheme <> "example.com",
        protocol: type,
        target: %{node_type: "normal", network_uuid: network.uuid}
      })
      |> TbgNodes.Repo.insert()

    {:ok, interface}
  end

  def permissioned_network_with_interfaces_fixture(%TbgNodes.Users.User{} = user) do
    {:ok, network} =
      TbgNodes.PermissionedEthereumNetworks.create_network_db(TbgNodes.Repo, user, @valid_network)

    {:ok, interface} = network_external_interface_fixture(network, "http")
    {:ok} = TbgNodes.PermissionedEthereumNetworks.add_basicauth_creds(TbgNodes.Repo, interface)
    {:ok, interface} = network_external_interface_fixture(network, "websocket")
    {:ok} = TbgNodes.PermissionedEthereumNetworks.add_basicauth_creds(TbgNodes.Repo, interface)
    TbgNodes.PermissionedEthereumNetworks.get_network_for_user_by_uuid!(user.id, network.uuid)
  end

  def besu_node_fixture(%TbgNodes.PermissionedEthereumNetworks.Network{} = network) do
    {:ok, node} =
      TbgNodes.PermissionedEthereumNetworks.create_node_db(TbgNodes.Repo, network, @valid_node)

    node
  end

  def create_k8s_conn(_context) do
    {:ok, k8s_conn} = K8s.Conn.lookup(:default)
    {:ok, k8s_conn: k8s_conn}
  end

  def create_test_ns(%{k8s_conn: k8s_conn}) do
    ns_name = "test-ns-#{:rand.uniform(1000)}"

    ns_k8s = get_ns_k8s(ns_name)

    create_ns = K8s.Client.create(ns_k8s)
    {:ok, _} = K8s.Client.run(create_ns, k8s_conn)

    # K8s.Client.Runner.Wait block until the eval functions evals to true.
    is_active_op = K8s.Client.get(ns_k8s)
    eval = fn val -> val == %{"phase" => "Active"} end
    opts = [find: ["status"], eval: eval, timeout: 5]
    {:ok, _} = K8s.Client.Runner.Wait.run(is_active_op, k8s_conn, opts)

    {:ok, test_ns: ns_name}
  end

  def delete_ns_on_exit(%{test_ns: test_ns, k8s_conn: k8s_conn}) do
    if System.get_env("CLEANUP", "true") == "true" do
      ExUnit.Callbacks.on_exit(fn -> empty_test_namespace(test_ns, k8s_conn) end)
    else
      :ok
    end
  end

  def delete_network_on_exit(network, get_infra_api) do
    if System.get_env("CLEANUP", "true") == "true" do
      ExUnit.Callbacks.on_exit(fn ->
        try do
          {:ok} = get_infra_api.().delete_network(network)
        rescue
          err -> err
        end
      end)
    else
      :ok
    end
  end

  def empty_test_namespace(ns_name, k8s_conn) do
    ns_k8s = get_ns_k8s(ns_name)
    delete_test_ns = K8s.Client.delete(ns_k8s)
    {:ok, _} = K8s.Client.run(delete_test_ns, k8s_conn)
  end

  defp get_ns_k8s(name) do
    _ns_k8s = %{
      "apiVersion" => "v1",
      "kind" => "Namespace",
      "metadata" => %{
        "name" => name
      }
    }
  end

  def create_genesis_config(number_of_validators \\ 3) do
    addresses =
      1..number_of_validators
      |> Enum.map(fn _ ->
        {:ok, private_key} = TbgNodes.ETH.generate_private_key(:hex)
        private_key
      end)
      |> Enum.map(fn private_key ->
        {:ok, address} = TbgNodes.ETH.get_address(:hex, private_key)
        address
      end)

    TbgNodes.ETH.genesis(addresses)
  end

  def put_live_session(conn, key, value) do
    conn |> put_connect_info(%{session: %{key => value}})
  end
end
