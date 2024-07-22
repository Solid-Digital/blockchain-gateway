defmodule TbgNodesWeb.AdminNodesComponentLiveTest do
  use TbgNodesWeb.ConnCase
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  describe "mount/3 as admin" do
    setup [
      :create_admin,
      :auth_user,
      :create_permissioned_ethereum_network,
      :create_permissioned_node
    ]

    test "shows nodes total", %{conn: conn} do
      {:ok, _view, html} =
        live_isolated(conn, TbgNodesWeb.AdminNodesComponentLive, id: "nodes-overview")

      count_nodes = TbgNodes.PermissionedEthereumNetworks.list_nodes() |> length()

      expected = "<div class=\"clr-row info-value\">#{count_nodes}</div>"

      assert html =~ expected
    end

    test "shows nodes table", %{conn: conn, permissioned_node: node} do
      {:ok, _view, html} =
        live_isolated(conn, TbgNodesWeb.AdminNodesComponentLive, id: "nodes-overview")

      row_1_c1 = "<td class=\"left\">" <> node.name <> "</td>"
      row_1_c2 = "<td class=\"right\">" <> node.node_type <> "</td>"
      row_1_c3 = "<td class=\"right\">placeholder</td>"
      row_1_c4 = "<td class=\"right\">" <> node.network.name <> "</td>"
      row_1_c5 = "<td class=\"right\">" <> node.network.user.email <> "</td>"

      assert html =~ row_1_c1
      assert html =~ row_1_c2
      assert html =~ row_1_c3
      assert html =~ row_1_c4
      assert html =~ row_1_c5
    end
  end
end
