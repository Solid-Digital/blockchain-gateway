defmodule PermissionedEthereumNetworkCreateFormLiveComponentTest do
  @moduledoc """
  The Permissioned Live Component Test.
  """

  use TbgNodesWeb.ConnCase
  use Phoenix.HTML

  import Phoenix.LiveViewTest
  import Ecto.Query

  @endpoint TbgNodesWeb.Endpoint

  @user_input %{
    :network_name => "test network",
    :number_besu_validators => 1,
    :number_besu_normal_nodes => 1,
    :number_besu_boot_nodes => 1,
    :deployment_option => "cloud",
    :consensus => "IBFT"
  }

  describe "mount/3" do
    setup [:create_user, :auth_user]

    test "permissioned form component is visible", %{conn: conn, user: user} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      html =
        view
        |> element("#besu")
        |> render_click()

      expected = "<div data-phx-component=\"2\" id=\"stateful-besu\">"
      assert html =~ expected
    end
  end

  describe "handle_change_nodes" do
    setup [:create_user, :auth_user]

    test "inreases normal nodes", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      expected = "<span id=\"normal-nodes\" class=\"node-numbers\">2</span>"

      assert view
             |> element("#increase-nodes")
             |> render_click() =~ expected
    end

    test "decreases normal nodes", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      expected = "<span id=\"normal-nodes\" class=\"node-numbers\">2</span>"

      for _ <- 0..1 do
        view
        |> element("#increase-nodes")
        |> render_click()
      end

      assert view
             |> element("#decrease-nodes")
             |> render_click() =~ expected
    end

    test "normal nodes cant be increased higher than 5", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      for _ <- 0..4 do
        view
        |> element("#increase-nodes")
        |> render_click()
      end

      expected = "<span id=\"normal-nodes\" class=\"node-numbers\">5</span>"

      assert view
             |> element("#increase-nodes")
             |> render_click() =~ expected
    end

    test "normal nodes cant be decreased below 1", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      expected = "<span id=\"normal-nodes\" class=\"node-numbers\">1</span>"

      assert view
             |> element("#decrease-nodes")
             |> render_click() =~ expected
    end

    test "increases validator nodes", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      expected = "<span id=\"validator-nodes\" class=\"node-numbers\">2</span>"

      assert view
             |> element("#increase-validators")
             |> render_click() =~ expected
    end

    test "decreases validator nodes", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      expected = "<span id=\"validator-nodes\" class=\"node-numbers\">2</span>"

      for _ <- 0..1 do
        view
        |> element("#increase-validators")
        |> render_click()
      end

      assert view
             |> element("#decrease-validators")
             |> render_click() =~ expected
    end

    test "validator nodes cant be increased higher than 5", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      for _ <- 0..4 do
        view
        |> element("#increase-validators")
        |> render_click()
      end

      expected = "<span id=\"validator-nodes\" class=\"node-numbers\">5</span>"

      assert view
             |> element("#increase-validators")
             |> render_click() =~ expected
    end

    test "validator nodes cant be decreased below 1", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      expected = "<span id=\"validator-nodes\" class=\"node-numbers\">1</span>"

      assert view
             |> element("#decrease-validators")
             |> render_click() =~ expected
    end
  end

  describe "handle_event validate" do
    setup [:create_user, :auth_user]

    test "invalidated changeset shows tooltip on submit", %{conn: conn, user: user} do
      user_input = %{
        "network_name" => "",
        "deployment_option" => "cloud",
        "consensus" => "IBFT",
        "number_besu_validators" => 1,
        "number_besu_normal_nodes" => 1
      }

      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      expected =
        "<span class=\"tooltip-content\">Complete the form to create a new network.</span>"

      view
      |> element("#besu")
      |> render_click()

      assert view
             |> element("#stateful-besu form")
             |> render_change(%{"network_user_input" => user_input}) =~ expected
    end
  end

  describe "handle_event create_network" do
    setup [:create_user, :auth_user]

    test "create network shows loading state on submit", %{user: user, conn: conn} do
      {:ok, view, _html} = conn |> put_live_session("current_user", user) |> live("/networks/new")

      view
      |> element("#besu")
      |> render_click()

      expected_html =
        "<button class=\"btn btn-spinner\" type=\"submit\" disabled=\"disabled\"><span class=\"spinner spinner-inline\"></span></button>"

      html =
        view
        |> element("#stateful-besu form")
        |> render_submit(%{"network_user_input" => @user_input})

      # allow time for create_network to create permissioned network (otherwise test errors occur related to
      # DB connections).
      :timer.sleep(500)

      name = @user_input[:network_name]

      _network =
        from(network in TbgNodes.PermissionedEthereumNetworks.Network,
          where: network.name == ^name and network.user_id == ^user.id,
          select: network
        )
        |> TbgNodes.Repo.one!()

      assert html =~ expected_html
    end
  end
end
