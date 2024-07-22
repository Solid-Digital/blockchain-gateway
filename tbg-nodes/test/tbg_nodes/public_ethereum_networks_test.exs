defmodule TbgNodes.PublicEthereumNetworksTest do
  use TbgNodes.DataCase

  alias TbgNodes.PublicEthereumNetworks

  describe "public ethereum networks" do
    setup [:create_user, :create_public_ethereum_network_with_interfaces]

    @valid_attrs %{
      name: "alpha-network",
      deployment_type: "shared",
      network_configuration: "ropsten",
      user_id: 1
    }
    @update_attrs %{
      name: "beta-network",
      deployment_type: "shared",
      network_configuration: "mainnet",
      user_id: 1
    }
    @invalid_attrs %{name: nil, deployment_type: nil, network_configuration: nil}

    test "list_ethereum_networks/0 returns all ethereum networks", %{public_network: %{id: id}} do
      assert [%PublicEthereumNetworks.Network{id: ^id}] = PublicEthereumNetworks.list_networks()
    end

    test "get_ethereum_network!/1 returns the network with given id", %{
      public_network: %{id: id} = network
    } do
      assert %PublicEthereumNetworks.Network{id: ^id} =
               PublicEthereumNetworks.get_network_by_uuid!(network.uuid)
    end

    test "create_ethereum_network/1 with valid data creates an ethereum network" do
      assert {:ok, %PublicEthereumNetworks.Network{} = network} =
               PublicEthereumNetworks.create_network(@valid_attrs)

      assert network.name == "alpha-network"
    end

    test "create_ethereum_network/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = PublicEthereumNetworks.create_network(@invalid_attrs)
    end

    test "update_ethereum_network/2 with valid data updates the network", %{
      public_network: network
    } do
      assert {:ok, %PublicEthereumNetworks.Network{} = network} =
               PublicEthereumNetworks.update_network(network, @update_attrs)

      assert network.name == "beta-network"
      assert network.network_configuration == "mainnet"
    end

    test "update_ethereum_network/2 with invalid data returns error changeset", %{
      public_network: %{id: id} = network
    } do
      assert {:error, %Ecto.Changeset{}} =
               PublicEthereumNetworks.update_network(network, @invalid_attrs)

      assert %PublicEthereumNetworks.Network{id: ^id} =
               PublicEthereumNetworks.get_network_by_uuid!(network.uuid)
    end

    test "delete_ethereum_network/1 deletes the network", %{public_network: network, user: user} do
      assert {:ok, %PublicEthereumNetworks.Network{}} =
               PublicEthereumNetworks.delete_network(network.uuid, user.id)

      assert_raise Ecto.NoResultsError, fn ->
        PublicEthereumNetworks.get_network_by_uuid!(network.uuid)
      end
    end

    test "change_ethereum_network/1 returns a network changeset", %{public_network: network} do
      assert %Ecto.Changeset{} = PublicEthereumNetworks.change_network(network)
    end

    test "network is invalid if only one of the fields is set" do
      networks = [
        %PublicEthereumNetworks.Network{name: "name"},
        %PublicEthereumNetworks.Network{network_configuration: "network_configuration"},
        %PublicEthereumNetworks.Network{deployment_type: "deployment_type"}
      ]

      for network <- networks do
        refute TbgNodes.PublicEthereumNetworks.Network.changeset_new(network, %{}).valid?
      end
    end

    test "network is valid if all of the fields are set" do
      network = %PublicEthereumNetworks.Network{
        name: "name",
        network_configuration: "network_configuration",
        deployment_type: "deployment_type"
      }

      assert TbgNodes.PublicEthereumNetworks.Network.changeset_new(network, %{}).valid?
    end
  end
end
