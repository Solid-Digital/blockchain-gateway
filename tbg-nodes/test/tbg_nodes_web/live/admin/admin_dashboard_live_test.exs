defmodule TbgNodesWeb.AdminDashboardLiveTest do
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
      assert html_response(conn, 200) =~ "Admin Dashboard"
    end

    test "can switch between users and nodes", %{conn: conn, user: user} do
      {:ok, view, html} =
        conn
        |> put_live_session("current_user", user)
        |> live("/admin")

      # starts on users overview
      assert html =~ "Total users"

      assert render_click(view, "select_tab", %{"active_tab" => "nodes"}) =~ "Total nodes"
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
