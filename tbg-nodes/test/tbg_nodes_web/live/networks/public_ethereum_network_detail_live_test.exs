defmodule TbgNodesWeb.PublicEthereumNetworkDetailLiveTest do
  import Plug.Conn
  import Phoenix.ConnTest
  import Phoenix.LiveViewTest
  use TbgNodesWeb.ConnCase
  @endpoint TbgNodesWeb.Endpoint

  describe "handle_params" do
    setup [
      :create_user,
      :auth_user,
      :create_public_ethereum_network_with_interfaces,
      :with_network_monitor
    ]

    test "renders page with correct network detail and info", %{
      conn: conn,
      user: user,
      public_network: network
    } do
      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/public-ethereum/" <> network.uuid)

      expected_network_name = "<h2>" <> network.name <> "</h2>"

      expected_creation_date =
        "<div class=\"clr-row info-value\">" <>
          Timex.format!(network.inserted_at, "%B %-d, %Y", :strftime) <> "</div>"

      assert html =~ expected_network_name
      assert html =~ expected_creation_date
    end

    test "renders network interfaces correctly", %{
      conn: conn,
      user: user,
      public_network: network
    } do
      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/public-ethereum/" <> network.uuid)

      expected_interface_title_http = "<h3 id=\"http\">HTTP</h3>"
      expected_interface_title_ws = "<h3 id=\"websocket\">Websocket</h3>"

      assert html =~ expected_interface_title_http
      assert html =~ expected_interface_title_ws
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
  end

  describe "delete network" do
    setup [
      :create_user,
      :auth_user,
      :create_public_ethereum_network_with_interfaces,
      :with_network_monitor
    ]

    test "toggles modal", %{conn: conn, user: user, public_network: network} do
      {:ok, view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/public-ethereum/" <> network.uuid)

      refute html =~ "<div class=\"delete-network-modal modal\">"

      assert render_click(
               view,
               :toggle_delete_modal
             ) =~ "<div class=\"delete-network-modal modal\">"
    end

    test "successfully deletes and redirects", %{conn: conn, user: user, public_network: network} do
      {:ok, view, _html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/networks/public-ethereum/" <> network.uuid)

      view
      |> render_click("delete_network", %{uuid: network.uuid})

      flash = assert_redirected(view, "/networks")
      assert flash["info"] == "Network deleted successfully."
    end
  end
end
