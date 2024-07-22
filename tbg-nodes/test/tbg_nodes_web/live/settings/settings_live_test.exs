defmodule TbgNodesWeb.SettingsLiveTest do
  use TbgNodesWeb.ConnCase
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  #  alias TbgNodesWeb.SettingsLive

  describe "mount/3" do
    setup [:create_user, :auth_user]

    test "renders statics", %{conn: conn} do
      conn = get(conn, "/settings")
      assert html_response(conn, 200) =~ "<h1>Settings</h1>"
    end

    test "renders live", %{conn: conn} do
      {:ok, _view, html} = live(conn, "/settings")
      assert html =~ "<h1>Settings</h1>"
    end

    test "shows email", %{conn: conn, user: user} do
      {:ok, _view, html} = live(conn, "/settings")
      assert html =~ "<div>" <> user.email <> "</div>"
    end

    test "shows since", %{conn: conn, user: user} do
      {:ok, _view, html} = live(conn, "/settings")
      assert html =~ "Since " <> TbgNodesWeb.ViewHelpers.since(user.inserted_at)
    end
  end
end
