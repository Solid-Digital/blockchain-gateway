defmodule TbgNodesWeb.BesuNetworkDetailLiveTest do
  @moduledoc """
  The Besu Network Detail Live Test.
  """
  import Plug.Conn
  import Phoenix.ConnTest
  import Phoenix.LiveViewTest
  use TbgNodesWeb.ConnCase
  @endpoint TbgNodesWeb.Endpoint

  describe "mount/3" do
    setup [
      :create_user,
      :auth_user,
      :create_permissioned_ethereum_network_with_interfaces,
      :with_network_monitor
    ]

    test "renders page with correct network detail and info", %{
      conn: conn,
      user: user,
      permissioned_network: network
    } do
      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/besu/" <> network.uuid)

      expected_network_name = "<h2>" <> network.name <> "</h2>"
      expected_consensus = "<div class=\"clr-row info-value\">" <> network.consensus <> "</div>"

      expected_creation_date =
        "<div class=\"clr-row info-value\">" <>
          Timex.format!(network.inserted_at, "%B %-d, %Y", :strftime) <> "</div>"

      assert html =~ expected_network_name
      assert html =~ expected_consensus
      assert html =~ expected_creation_date
    end

    test "renders network interfaces correctly", %{
      conn: conn,
      user: user,
      permissioned_network: network
    } do
      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/besu/" <> network.uuid)

      expected_interface_title_http = "<h3 id=\"http\">HTTP</h3>"
      expected_interface_title_ws = "<h3 id=\"websocket\">Websocket</h3>"

      assert html =~ expected_interface_title_http
      assert html =~ expected_interface_title_ws
    end

    test "renders nodes table correctly", %{
      conn: conn,
      user: user,
      permissioned_network: network
    } do
      _ = besu_node_fixture(network)

      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/besu/" <> network.uuid)

      expected_nodes_table =
        "<tr class=\"node-row\"><td class=\"left\">i am a node</td><td class=\"left\">Validator</td>"

      assert html =~ expected_nodes_table
    end

    test "routing to an unknown network uuid returns a 400", %{conn: conn} do
      assert_error_sent 400, fn ->
        get(
          conn,
          Routes.live_path(
            conn,
            TbgNodesWeb.Networks.PermissionedBesuNetworkDetailLive,
            "ridiculous-uuid"
          )
        )
      end
    end

    test "renders curl commands correctly", %{
      conn: conn,
      user: user,
      permissioned_network: network
    } do
      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/besu/" <> network.uuid)

      http_interface =
        network.external_interfaces
        |> Enum.filter(fn interface ->
          interface.protocol == "http"
        end)
        |> List.first()

      ws_interface =
        network.external_interfaces
        |> Enum.filter(fn interface ->
          interface.protocol == "websocket"
        end)
        |> List.first()

      http_creds = http_interface.basicauth_creds
      ws_creds = ws_interface.basicauth_creds

      http_username = List.first(http_creds).username
      http_password = List.first(http_creds).password
      assert http_interface.url == "https://example.com"

      ws_username = List.first(ws_creds).username
      ws_password = List.first(ws_creds).password
      assert ws_interface.url == "wss://example.com"

      expected_http_command =
        "curl -X POST --data '{&quot;jsonrpc&quot;:&quot;2.0&quot;,&quot;method&quot;:&quot;eth_protocolVersion&quot;,&quot;params&quot;:[],&quot;id&quot;:67}' https://#{
          http_username
        }:#{http_password}@example.com"

      expected_ws_command = "wscat -c wss://#{ws_username}:#{ws_password}@example.com"

      assert html =~ expected_http_command
      assert html =~ expected_ws_command
    end
  end

  describe "delete network" do
    setup [
      :create_user,
      :auth_user,
      :create_permissioned_ethereum_network_with_interfaces,
      :with_network_monitor
    ]

    test "toggles modal", %{conn: conn, user: user, permissioned_network: network} do
      {:ok, view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/besu/" <> network.uuid)

      refute html =~ "<div class=\"delete-network-modal modal\">"

      assert render_click(
               view,
               :toggle_delete_modal
             ) =~ "<div class=\"delete-network-modal modal\">"
    end

    test "successfully deletes and redirects", %{
      conn: conn,
      user: user,
      permissioned_network: network
    } do
      {:ok, view, _html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/besu/" <> network.uuid)

      view
      |> render_click("delete_network", %{uuid: network.uuid})

      flash = assert_redirected(view, "/networks")
      assert flash["info"] == "Network deleted successfully."
    end
  end
end
