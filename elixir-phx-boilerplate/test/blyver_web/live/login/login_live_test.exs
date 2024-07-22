defmodule BlyverWeb.Login.LoginLiveTest do
  use BlyverWeb.ConnCase
  import Phoenix.LiveViewTest

  setup [:create_default_user]

  test "validates and submits the login form", %{conn: conn, user: user} do
    user_attrs = %{
      email: user.email,
      password: "SomeValidPa$$word"
    }

    invalid_attrs = %{email: "test@email.com"}
    does_not_exist_attrs = %{email: "does_not_exist@email.com", password: "SomeValidPa$$word"}

    {:ok, login_live, html} = live_isolated(conn, BlyverWeb.Login.LoginLive)

    assert html =~ "Login to your Blyver account"

    assert login_live |> element("form")

    assert login_live
           |> form("form", user: invalid_attrs)
           |> render_submit() =~ "can&#39;t be blank"

    form =
      login_live
      |> form("form", user: does_not_exist_attrs)

    assert render_submit(form) =~ ~r/phx-trigger-action/
    conn = follow_trigger_action(form, conn)

    assert html_response(conn, 200) =~ "Invalid email and password combination"

    form =
      login_live
      |> form("form", user: user_attrs)

    assert render_submit(form) =~ ~r/phx-trigger-action/

    conn = follow_trigger_action(form, conn)
    assert redirected_to(conn, 302) =~ Routes.dashboard_path(conn, :index)
  end
end
