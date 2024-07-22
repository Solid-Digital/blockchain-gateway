defmodule TbgNodesWeb.AdminDashboardControllerTest do
  @moduledoc """
  """

  use TbgNodesWeb.ConnCase
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  describe "index for admin" do
    setup [:create_admin, :auth_user, :create_user]

    test "admin nav link visible for admin", %{conn: conn} do
      conn = get(conn, Routes.live_path(conn, TbgNodesWeb.NetworkLive))
      assert html_response(conn, 200) =~ "Admin"
    end

    test "admin user can render admin dashboard", %{conn: conn} do
      conn = get(conn, Routes.live_path(conn, TbgNodesWeb.AdminDashboardLive))
      assert html_response(conn, 200) =~ "Users"
    end

    test "shows users total", %{conn: conn, user: user} do
      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/admin")

      user_count = TbgNodes.Users.list_users() |> length()

      expected = "<div class=\"clr-row info-value\">#{user_count}</div>"

      assert html =~ expected
    end

    test "shows node and network count", %{conn: conn, user: user} do
      create_public_ethereum_network_with_interfaces(%{user: user})
      {:ok, [permissioned_network: network]} = create_permissioned_ethereum_network(%{user: user})
      create_permissioned_node(%{permissioned_network: network})

      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/admin")

      expected =
        ~s(<tr><td class=\"left\">#{user.email}</td><td class=\"right perm-networks-col\">1</td><td class=\"right perm-nodes-col\">1</td><td class=\"right pub-networks-col\">1</td><td class=\"right\">#{
          TbgNodesWeb.ViewHelpers.format_date(user.inserted_at)
        }</td></tr>)

      assert html =~ expected
    end

    test "shows users table", %{conn: conn, user: user} do
      user1 = user_fixture()

      {:ok, _view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/admin")

      row_1_c1 = "<td class=\"left\">" <> user.email <> "</td>"
      row_2_c1 = "<td class=\"left\">" <> user1.email <> "</td>"

      assert html =~ row_1_c1
      assert html =~ row_2_c1
    end
  end

  describe "index for normal user" do
    setup [:create_user, :auth_user]

    test "admin nav link invisible for user", %{conn: conn} do
      conn = get(conn, Routes.live_path(conn, TbgNodesWeb.NetworkLive))
      assert html_response(conn, 200) != "Admin"
    end

    test "normal user cannot render admin dashboard", %{conn: conn} do
      conn = get(conn, Routes.live_path(conn, TbgNodesWeb.AdminDashboardLive))
      assert html_response(conn, 302)
    end
  end
end
