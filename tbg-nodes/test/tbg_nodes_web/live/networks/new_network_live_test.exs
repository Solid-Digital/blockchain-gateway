defmodule TbgNodesWeb.NewNetworkLiveTest do
  @moduledoc """
  The New Network Live Test.
  """

  use TbgNodesWeb.ConnCase
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  alias TbgNodes.PermissionedEthereumNetworks

  @create_attrs %{
    name: "name",
    network_configuration: "mainnet",
    deployment_type: "shared"
  }

  describe "mount/3" do
    setup [:create_user, :auth_user]

    test "new network live is visible", %{conn: conn, user: user} do
      {:ok, _view, html} = conn |> put_live_session("current_user", user) |> live("/networks/new")
      expected = "<div class=\"new-network main-container center p-v-m\">"

      assert html =~ expected
    end
  end

  describe "handle_event select_network_type" do
    setup [:create_user, :auth_user, :create_public_ethereum_network_with_interfaces]

    test "select network type selects ethereum network type", %{conn: conn, user: user} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")
      expected = "<div id=\"ethereum-card\" class=\"card clickable card-img active\">"

      assert render_click(view, :select_network_type, %{"network_type" => "ethereum"}) =~ expected
    end
  end

  describe "renders PublicEthereumNetworkFormLiveComponent" do
    setup [:create_user, :auth_user]

    test "public live form component is visible", %{conn: conn, user: user} do
      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> Phoenix.LiveViewTest.live("/networks/new")

      expected = "<div data-phx-component=\"1\" id=\"stateful-1\">"

      assert html =~ expected
    end

    test "create network redirects when data is valid", %{conn: conn, user: user} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      html =
        view
        |> element("#stateful-1 form", "")
        |> render_submit(%{"network_create_form" => @create_attrs})

      {:error, {:live_redirect, %{flash: _}}} = assert html
    end
  end

  describe "renders PermissionedEthereumNetworkLiveComponent" do
    setup [:create_user, :auth_user]

    test "permissioned live form component is visible", %{conn: conn, user: user} do
      {:ok, view, _html} =
        conn
        |> put_live_session("current_user", user)
        |> Phoenix.LiveViewTest.live("/networks/new")

      expected = "<div data-phx-component=\"2\" id=\"stateful-besu\">"

      assert view
             |> element("#besu")
             |> render_click() =~ expected
    end
  end

  describe "handle_info create_permissioned_besu_network" do
    setup [:create_user, :auth_user, :with_network_monitor]

    test "redirects after successful network creation", %{conn: conn, user: user} do
      {:ok, view, static_html} =
        conn |> put_live_session("current_user", user) |> live("/networks/new")

      assert static_html =~ "<h2>Create network</h2>"

      user_input = %{
        :network_name => "test network",
        :number_besu_validators => 1,
        :number_besu_normal_nodes => 1,
        :number_besu_boot_nodes => 1,
        :deployment_option => "cloud",
        :consensus => "IBFT"
      }

      changeset =
        PermissionedEthereumNetworks.NetworkUserInput.changeset(
          %PermissionedEthereumNetworks.NetworkUserInput{},
          user_input
        )

      send(view.pid, {:create_permissioned_besu_network, changeset})

      %{proxy: {ref, topic, _}} = view
      # wait for up to 200ms for message for network creation redirect
      timeout = 200

      receive do
        {^ref, {:redirect, ^topic, %{to: to}}} ->
          # Assert match on regex pattern that matches the route /networks/besu/:uuid
          assert String.match?(
                   to,
                   ~r/\/networks\/besu\/[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}/
                 )

          # Assert  that the detail page reflects the created network
          {:ok, _view, detail_html} = conn |> put_live_session("current_user", user) |> live(to)
          assert detail_html =~ "<h2>" <> user_input[:network_name] <> "</h2>"
      after
        timeout ->
          raise "timed out"
      end
    end
  end
end
