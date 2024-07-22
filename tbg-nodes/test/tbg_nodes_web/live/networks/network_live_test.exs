defmodule TbgNodesWeb.NetworkControllerTest do
  use TbgNodesWeb.ConnCase
  use Phoenix.HTML

  @endpoint TbgNodesWeb.Endpoint

  describe "when user is logged in and has no networks" do
    setup [:create_user, :auth_user]

    test "index shows no networks", %{conn: conn} do
      conn = get(conn, "/networks")
      response = html_response(conn, 200)

      assert response =~ "Setup your first network"
    end
  end

  describe "when user is logged in and has public ethereum networks" do
    setup [
      :create_user,
      :auth_user,
      :create_public_ethereum_network_with_interfaces,
      :with_network_monitor
    ]

    test "index lists all public networks", %{conn: conn, public_network: network} do
      conn = get(conn, Routes.live_path(conn, TbgNodesWeb.NetworkLive))
      response = html_response(conn, 200)

      assert response =~ network.name
    end

    test "show network renders detail page", %{conn: conn, public_network: network} do
      conn =
        get(
          conn,
          Routes.live_path(
            conn,
            TbgNodesWeb.Networks.PublicEthereumNetworkDetailLive,
            network.uuid
          )
        )

      assert html_response(conn, 200) =~ network.name
    end

    test "routing to an unknown network uuid returns a 400", %{conn: conn} do
      assert_error_sent 400, fn ->
        get(
          conn,
          Routes.live_path(
            conn,
            TbgNodesWeb.Networks.PublicEthereumNetworkDetailLive,
            "ridiculous-uuid"
          )
        )
      end
    end
  end

  describe "when user is logged in and has permissioned networks" do
    setup [
      :create_user,
      :auth_user,
      :create_permissioned_ethereum_network_with_interfaces,
      :with_network_monitor
    ]

    test "index lists all permissioned networks", %{conn: conn, permissioned_network: network} do
      conn = get(conn, Routes.live_path(conn, TbgNodesWeb.NetworkLive))
      response = html_response(conn, 200)

      assert response =~ network.name
    end

    test "show network renders detail page", %{conn: conn, permissioned_network: network} do
      conn =
        get(
          conn,
          Routes.live_path(
            conn,
            TbgNodesWeb.Networks.PermissionedBesuNetworkDetailLive,
            network.uuid
          )
        )

      assert html_response(conn, 200) =~ network.name
    end

    test "routing to an unknown network uuid returns a 400", %{conn: conn} do
      assert_error_sent 400, fn ->
        get(
          conn,
          Routes.live_path(
            conn,
            TbgNodesWeb.Networks.PublicEthereumNetworkDetailLive,
            "ridiculous-uuid"
          )
        )
      end
    end
  end
end
