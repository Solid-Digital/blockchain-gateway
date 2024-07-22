defmodule TbgNodesWeb.Networks.PublicEthereumNetworkCreateFormLiveComponentTest do
  @moduledoc """
  The Public Live Component Test.
  """

  use TbgNodesWeb.ConnCase
  use Phoenix.HTML
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  alias TbgNodesWeb.Networks.PublicEthereumNetworkCreateFormLiveComponent

  describe "mount/3" do
    setup [:create_user, :auth_user]

    test "public form component is visible", %{conn: conn, user: user} do
      {:ok, _view, _html} =
        conn |> put_live_session("current_user", user) |> live("/networks/new")

      rendered_component =
        render_component(PublicEthereumNetworkCreateFormLiveComponent,
          id: 1,
          current_user_id: user.id,
          network_type: "ethereum"
        )

      expected = "<div id=\"stateful-1\">"
      assert rendered_component =~ expected
    end
  end

  describe "handle_event validate" do
    setup [:create_user, :auth_user, :create_public_ethereum_network_with_interfaces]

    test "validate returns error when name input is blank", %{conn: conn, user: user} do
      network = %{
        "deployment_type" => "shared",
        "name" => "",
        "network_configuration" => "mainnet",
        "protocol" => "ethereum",
        "user_id" => user.id
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      assert view
             |> element("#stateful-1 form")
             |> render_change(%{"network" => network}) =~
               "<span class=\"tooltip-content\">All of these fields must be present: Network Name, Protocol, Network Configuration, and Deployment type </span>"
    end

    test "validate returns error when no deployment type is selected", %{conn: conn, user: user} do
      network = %{
        "deployment_type" => "",
        "name" => "new network",
        "network_configuration" => "mainnet",
        "protocol" => "ethereum",
        "user_id" => user.id
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      expected =
        "<span class=\"tooltip-content\">All of these fields must be present: Network Name, Protocol, Network Configuration, and Deployment type </span>"

      assert view
             |> element("#stateful-1 form")
             |> render_change(%{"network" => network}) =~ expected
    end

    test "validate returns error when no network configuration is selected", %{
      conn: conn,
      user: user
    } do
      network = %{
        "deployment_type" => "",
        "name" => "new network",
        "network_configuration" => "",
        "protocol" => "ethereum",
        "user_id" => user.id
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      expected =
        "<span class=\"tooltip-content\">All of these fields must be present: Network Name, Protocol, Network Configuration, and Deployment type </span>"

      assert view
             |> element("#stateful-1 form")
             |> render_change(%{"network" => network}) =~ expected
    end

    test "validate returns error when no protocol is selected", %{conn: conn, user: user} do
      network = %{
        "deployment_type" => "",
        "name" => "new network",
        "network_configuration" => "mainnet",
        "protocol" => "",
        "user_id" => user.id
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      expected =
        "<span class=\"tooltip-content\">All of these fields must be present: Network Name, Protocol, Network Configuration, and Deployment type </span>"

      assert view
             |> element("#stateful-1 form")
             |> render_change(%{"network" => network}) =~ expected
    end

    test "validate returns no errors with valid input", %{conn: conn, user: user} do
      network = %{
        "deployment_type" => "shared",
        "name" => "name",
        "network_configuration" => "mainnet",
        "protocol" => "ethereum",
        "user_id" => user.id
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      assert view
             |> element("#stateful-1 form")
             |> render_change(%{"network" => network}) =~ ""
    end
  end

  describe "archive_data" do
    test "create network with archive data appends to network_configuration", %{} do
      params = %{
        "name" => "name",
        "deployment_type" => "shared",
        "protocol" => "ethereum",
        "network_configuration" => "mainnet",
        "archive_data" => "true"
      }

      params = PublicEthereumNetworkCreateFormLiveComponent.maybe_archive_data(params)

      assert Map.get(params, "network_configuration") == "mainnet-archive"
    end
  end

  describe "handle_event create_network" do
    setup [:create_user, :auth_user, :create_public_ethereum_network_with_interfaces]

    test "create network submits when form input is valid", %{conn: conn, user: user} do
      network = %{
        name: "name",
        protocol: "ethereum",
        network_configuration: "mainnet",
        deployment_type: "shared"
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      html =
        view
        |> element("#stateful-1 form", "")
        |> render_submit(%{"network_create_form" => network})

      {:error, {:live_redirect, %{flash: _}}} = assert html
    end

    test "create network re-renders form when input is invalid", %{conn: conn, user: user} do
      attrs = %{
        name: "",
        protocol: "ethereum",
        network_configuration: "mainnet",
        deployment_type: "shared"
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")
      expected = "<div class=\"new-network main-container center p-v-m\">"

      html =
        view
        |> element("#stateful-1 form", "")
        |> render_submit(%{"network" => attrs})

      assert html =~ expected
    end
  end
end
