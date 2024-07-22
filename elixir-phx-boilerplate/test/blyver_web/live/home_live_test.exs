defmodule BlyverWeb.PageLiveTest do
  use BlyverWeb.ConnCase

  import Phoenix.LiveViewTest

  test "disconnected and connected render", %{conn: conn} do
    {:ok, page_live, disconnected_html} = live(conn, "/")
    assert disconnected_html =~ "What if investing in property was more equal?"
    assert render(page_live) =~ "What if investing in property was more equal?"
  end
end
