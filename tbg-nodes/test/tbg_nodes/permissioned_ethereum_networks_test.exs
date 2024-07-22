defmodule TbgNodes.PermissionedEthereumNetworksTest do
  use TbgNodes.DataCase
  alias TbgNodes.PermissionedEthereumNetworks

  @moduletag :PermissionedEthereumNetworks

  @valid_user_input %{
    network_name: "permissioned_network_1",
    number_besu_validators: 2,
    number_besu_normal_nodes: 3,
    number_besu_boot_nodes: 1,
    join_network: false,
    managed_by: "unchain",
    consensus: "clique"
  }

  @valid_network %{
    name: "test-network-name",
    consensus: "clique"
  }

  @valid_node %{
    name: "i am a node",
    node_type: "validator"
  }

  describe "PermissionedEthereumNetworks.create_network/1" do
    setup [:create_user, :with_network_monitor]

    test "valid user input", %{user: user} do
      valid_user_input = permissioned_ethereum_network_user_input(@valid_user_input)

      {:ok, network} = PermissionedEthereumNetworks.create_network(valid_user_input, user.id)

      query =
        from network in PermissionedEthereumNetworks.Network,
          where: network.uuid == ^network.uuid,
          select: network,
          preload: [:besu_nodes, external_interfaces: :basicauth_creds]

      network = Repo.one(query)

      assert network.name == @valid_user_input.network_name

      assert Enum.count(network.besu_nodes, fn node ->
               node.node_type == "boot"
             end) == @valid_user_input.number_besu_boot_nodes

      assert Enum.count(network.besu_nodes, fn node ->
               node.node_type == "validator"
             end) == @valid_user_input.number_besu_validators

      assert Enum.count(network.besu_nodes, fn node ->
               node.node_type == "normal"
             end) == @valid_user_input.number_besu_normal_nodes

      assert Enum.count(network.external_interfaces) == 2

      assert network.config != nil

      Enum.filter(network.besu_nodes, fn node -> node.node_type == "validator" end)
      |> Enum.each(fn validator ->
        assert String.length(validator.address) > 0
        assert String.contains?(network.config["extraData"], validator.address)
      end)

      network.external_interfaces
      |> Enum.each(fn interface ->
        assert interface.url != nil
      end)

      network.external_interfaces
      |> Enum.flat_map(& &1.basicauth_creds)
      |> Enum.each(fn basicauth_cred ->
        assert basicauth_cred.username != nil
        assert basicauth_cred.password != nil
      end)
    end
  end

  describe "PermissionedEthereumNetworks.list_networks/1" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "list_networks/0 returns all ethereum networks", %{
      permissioned_network: %{id: id, uuid: uuid}
    } do
      assert PermissionedEthereumNetworks.list_networks()
             |> Enum.any?(
               &match?(%PermissionedEthereumNetworks.Network{id: ^id, uuid: ^uuid}, &1)
             )
    end
  end

  describe "PermissionedEthereumNetworks.list_nodes/1" do
    setup [
      :create_user,
      :create_permissioned_ethereum_network_with_interfaces,
      :create_permissioned_node
    ]

    test "list_nodes/0 returns all ethereum nodes", %{
      permissioned_network: %{uuid: network_uuid},
      permissioned_node: %{id: id, uuid: uuid}
    } do
      assert PermissionedEthereumNetworks.list_nodes()
             |> Enum.any?(
               &match?(
                 %PermissionedEthereumNetworks.BesuNode{
                   id: ^id,
                   uuid: ^uuid,
                   network: %{uuid: ^network_uuid}
                 },
                 &1
               )
             )
    end
  end

  describe "PermissionedEthereumNetworks.create_network_db/3" do
    setup [:create_user]

    test "works valid attrs", %{user: user} do
      {:ok, network} =
        PermissionedEthereumNetworks.create_network_db(TbgNodes.Repo, user, @valid_network)

      assert network.name == @valid_network.name
    end

    test "fails with invalid attrs", %{user: user} do
      invalid_attrs = %{blabla: "not-valid"}

      assert {:error, _} =
               PermissionedEthereumNetworks.create_network_db(TbgNodes.Repo, user, invalid_attrs)
    end
  end

  describe "PermissionedEthereumNetworks.create_node_db/3" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid attrs", %{permissioned_network: network} do
      {:ok, _node} =
        PermissionedEthereumNetworks.create_node_db(
          TbgNodes.Repo,
          network,
          @valid_node
        )
    end

    test "fails with invalid attrs", %{permissioned_network: network} do
      invalid_attrs = %{blabla: "not-valid"}

      result =
        PermissionedEthereumNetworks.create_node_db(
          TbgNodes.Repo,
          network,
          invalid_attrs
        )

      assert {:error, _} = result
    end
  end

  describe "PermissionedEthereumNetworks.create_network_config/2" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid arguments", %{permissioned_network: network} do
      node = besu_node_fixture(network)
      {:ok, dummy_config} = PermissionedEthereumNetworks.create_network_config([node], network)

      assert dummy_config != nil
    end
  end

  describe "PermissionedEthereumNetworks.add_network_config/4" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid arguments", %{user: user, permissioned_network: network} do
      config = %{test: "hello"}

      {:ok, network} =
        PermissionedEthereumNetworks.add_network_config(TbgNodes.Repo, user, network, config)

      assert network.config == config
    end

    test "fails with invalid config", %{user: user, permissioned_network: network} do
      invalid_config = %{}

      result =
        PermissionedEthereumNetworks.add_network_config(
          TbgNodes.Repo,
          user,
          network,
          invalid_config
        )

      assert {:error, _changeset} = result
    end
  end

  describe "PermissionedEthereumNetworks.deploy_network/1" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "deploy successfull", %{permissioned_network: network} do
      assert {:ok} =
               PermissionedEthereumNetworks.deploy_network(
                 network,
                 PermissionedEthereumNetworks.InfraAPIMock
               )
    end

    test "correct error return when deploy fails", %{permissioned_network: network} do
      defmodule TestInfraAPIFail do
        def deploy_network(_), do: raise("uh oh i failed")
      end

      expected_error_msg = "Creating k8s resources failed."
      result = PermissionedEthereumNetworks.deploy_network(network, TestInfraAPIFail)
      assert result == {:error, expected_error_msg}
    end

    test "correct error return when deploy times out", %{permissioned_network: network} do
      # Since the deploy_network takes 500ms but the max timeout is 250ms,
      # the function will trigger the timeout error.

      defmodule TestInfraAPITimeOut do
        def deploy_network(_) do
          :timer.sleep(500)
          {:ok}
        end
      end

      expected_error_msg = "Creating k8s resources timed-out."
      test_timeout = 250

      result =
        PermissionedEthereumNetworks.deploy_network(network, TestInfraAPITimeOut, test_timeout)

      assert result == {:error, expected_error_msg}
    end
  end

  describe "PermissionedEthereumNetworks.add_network_external_interface_db/3" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid arguments", %{permissioned_network: network} do
      {:ok, _interface} =
        PermissionedEthereumNetworks.add_network_external_interface_db(
          TbgNodes.Repo,
          network,
          "http"
        )

      query =
        from network in PermissionedEthereumNetworks.Network,
          where: network.uuid == ^network.uuid,
          select: network,
          preload: [:external_interfaces]

      network = Repo.one(query)
      assert Enum.count(network.external_interfaces) == 1
    end

    test "fails with invalid arguments", %{permissioned_network: network} do
      not_protocol = "not-http"

      {:error, _err} =
        PermissionedEthereumNetworks.add_network_external_interface_db(
          TbgNodes.Repo,
          network,
          not_protocol
        )

      query =
        from network in PermissionedEthereumNetworks.Network,
          where: network.uuid == ^network.uuid,
          select: network,
          preload: [:external_interfaces]

      network = Repo.one(query)
      assert Enum.empty?(network.external_interfaces)
    end
  end

  describe "PermissionedEthereumNetworks.add_basicauth_creds/2" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid arguments", %{permissioned_network: network} do
      {:ok, external_interface} = network_external_interface_fixture(network)
      {:ok} = PermissionedEthereumNetworks.add_basicauth_creds(TbgNodes.Repo, external_interface)

      query =
        from basicauth_cred in PermissionedEthereumNetworks.BasicauthCred,
          where: basicauth_cred.external_interface_uuid == ^external_interface.uuid,
          select: basicauth_cred

      basic_auth = TbgNodes.Repo.one(query)

      assert basic_auth.username != nil
      assert basic_auth.password != nil
    end
  end

  describe "PermissionedEthereumNetworks.deploy_network_external_interface/3" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid arguments", %{permissioned_network: network} do
      {:ok, before_deploy_interface} = network_external_interface_fixture(network)

      {:ok} =
        PermissionedEthereumNetworks.deploy_network_external_interface(
          TbgNodes.Repo,
          network,
          before_deploy_interface,
          PermissionedEthereumNetworks.InfraAPIMock
        )

      query =
        from network_interface in PermissionedEthereumNetworks.ExternalInterface,
          where: network_interface.network_id == ^network.id,
          select: network_interface

      after_deploy_interface = Repo.one(query)

      assert after_deploy_interface.url != nil
    end

    test "fails when infra_api fails to deploy", %{permissioned_network: network} do
      {:ok, before_deploy_interface} = network_external_interface_fixture(network)

      defmodule NoDeployInfraAPI do
        def deploy_external_interface(_), do: {:error, "network not deployed"}
      end

      result =
        PermissionedEthereumNetworks.deploy_network_external_interface(
          TbgNodes.Repo,
          network,
          before_deploy_interface,
          NoDeployInfraAPI
        )

      query =
        from network_interface in PermissionedEthereumNetworks.ExternalInterface,
          where: network_interface.network_id == ^network.id,
          select: network_interface

      after_deploy_interface = Repo.one(query)

      assert after_deploy_interface.url == "https://example.com"
      assert {:error, "network not deployed"} == result
    end
  end

  describe "PermissionedEthereumNetworks.get_network_for_user_by_uuid!/1" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid arguments", %{user: user, permissioned_network: network} do
      network = PermissionedEthereumNetworks.get_network_for_user_by_uuid!(user.id, network.uuid)
      assert network.besu_nodes != nil
    end

    test "uuid which doesn't exist raises error", %{user: user} do
      doesnt_exist_uuid = "123e4567-e89b-12d3-a456-426614174000"
      err_msg_regex = ~r/^expected at least one result but got none in query/

      assert_raise Ecto.NoResultsError, err_msg_regex, fn ->
        PermissionedEthereumNetworks.get_network_for_user_by_uuid!(user.id, doesnt_exist_uuid)
      end
    end

    test "user id which doesn't exist raises error", %{permissioned_network: network} do
      doesnt_exist_id = 123_123_123
      err_msg_regex = ~r/^expected at least one result but got none in query/

      assert_raise Ecto.NoResultsError, err_msg_regex, fn ->
        PermissionedEthereumNetworks.get_network_for_user_by_uuid!(doesnt_exist_id, network.uuid)
      end
    end
  end

  describe "PermissionedEthereumNetworks.get_network_by_uuid!/1" do
    setup [:create_user, :create_permissioned_ethereum_network]

    test "works with valid arguments", %{permissioned_network: network} do
      network = PermissionedEthereumNetworks.get_network_by_uuid!(network.uuid)
      assert network.besu_nodes != nil
    end

    test "uuid which doesn't exist raises error" do
      doesnt_exist_uuid = "123e4567-e89b-12d3-a456-426614174000"
      err_msg_regex = ~r/^expected at least one result but got none in query/

      assert_raise Ecto.NoResultsError, err_msg_regex, fn ->
        PermissionedEthereumNetworks.get_network_by_uuid!(doesnt_exist_uuid)
      end
    end
  end

  describe "PermissionedEthereumNetworks.delete_network/4" do
    setup [
      :create_user,
      :create_permissioned_ethereum_network_with_interfaces,
      :create_permissioned_node
    ]

    test "deletes network", %{user: user, permissioned_network: network} do
      result = PermissionedEthereumNetworks.delete_network_for_user(network.uuid, user.id)
      assert {:ok, _} = result

      network_query =
        from network in PermissionedEthereumNetworks.Network,
          where: network.uuid == ^network.uuid,
          select: network

      assert Repo.one(network_query) == nil

      node_query =
        from node in PermissionedEthereumNetworks.BesuNode,
          where: node.network_id == ^network.id,
          select: node

      assert Repo.one(node_query) == nil

      external_interface_query =
        from ei in PermissionedEthereumNetworks.ExternalInterface,
          where: ei.network_id == ^network.id,
          select: ei

      assert Repo.one(external_interface_query) == nil
    end

    test "fails if UUID is invalid", %{user: user} do
      doesnt_exist_uuid = "123e4567-e89b-12d3-a456-426614174000"

      assert_raise Ecto.NoResultsError, fn ->
        PermissionedEthereumNetworks.delete_network_for_user(doesnt_exist_uuid, user.id)
      end
    end
  end
end
