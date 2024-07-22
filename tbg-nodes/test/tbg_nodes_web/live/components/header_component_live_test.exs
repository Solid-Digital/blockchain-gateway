defmodule TbgNodesWeb.HeaderComponentLiveTest do
  use TbgNodesWeb.ConnCase
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  describe "mount/3" do
    setup [:create_user, :auth_user]

    test "header shows with menu options", %{conn: conn, user: user} do
      {:ok, _view, html} =
        live_isolated(
          conn,
          TbgNodesWeb.HeaderComponentLive,
          id: "header",
          session: %{
            "current_user_id" => user.id,
            "path" => ["/"]
          }
        )

      header_bar =
        "<header class=\"header\"><div class=\"branding\"><img class=\"logo\" src=\"/images/unchain_logo_darkbg.png\" alt=\"Unchain\"/></div>"

      assert html =~ header_bar

      networks_nav_item = "<a class=\"nav-link nav-text \" href=\"/networks\">Networks</a>"
      assert html =~ networks_nav_item

      documentation_nav_item =
        "<a href=\"https://support.unchain.io\" class=\"nav-link nav-text\" target=\"_blank\" rel=\"noreferrer\">\n\t\t\tDocumentation\n\t\t\t<i class=\"fas fa-external-link-alt\"></i></a>"

      assert html =~ documentation_nav_item

      dropdown_button_item =
        "<button id=\"header-user-dropdown-toggle\" @click=\"isOpen = !isOpen\" class=\"nav-text dropdown-toggle\">\n#{
          user.email
        }"

      dropdown_button_item2 =
        "\t\t\t\t<i class=\"fas fa-angle-down\"></i></button><div class=\"dropdown-menu\">"

      assert html =~ dropdown_button_item
      assert html =~ dropdown_button_item2
    end

    test "header menu option shows active", %{conn: conn, user: user} do
      {:ok, _view, html} =
        live_isolated(
          conn,
          TbgNodesWeb.HeaderComponentLive,
          id: "header",
          session: %{
            "current_user_id" => user.id,
            "path" => ["networks"]
          }
        )

      networks_nav_item = "<a class=\"nav-link nav-text active\" href=\"/networks\">Networks</a>"
      assert html =~ networks_nav_item
    end
  end

  describe "mount/3 unauthorized" do
    setup [:create_user]

    test "header does not display", %{conn: conn} do
      {:error, {:redirect, %{flash: _flash, to: address}}} = live(conn, "/networks")

      assert address == "/login"
    end
  end

  describe "mount/3 as admin" do
    setup [:create_admin, :auth_user]

    test "header shows admin menu options", %{conn: conn, user: user} do
      {:ok, _view, html} =
        live_isolated(
          conn,
          TbgNodesWeb.HeaderComponentLive,
          id: "header",
          session: %{
            "current_user_id" => user.id,
            "path" => ["/"]
          }
        )

      phoenix_dashboard_nav_item =
        "<a class=\"nav-link nav-text \" href=\"/admin/dashboard\">Phoenix Live Dashboard</a>"

      assert html =~ phoenix_dashboard_nav_item

      admin_nav_item = "<a class=\"nav-link nav-text \" href=\"/admin\">Admin</a>"
      assert html =~ admin_nav_item
    end
  end
end
