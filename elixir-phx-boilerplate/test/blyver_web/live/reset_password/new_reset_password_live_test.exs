defmodule BlyverWeb.ResetPassword.NewResetPasswordLiveTest do
  use BlyverWeb.ConnCase
  import Phoenix.LiveViewTest

  setup [:create_default_user]

  test "validates and submits the reset password form", %{conn: conn, user: user} do
    user_attrs = %{
      email: user.email
    }

    invalid_attrs = %{email: ""}

    {:ok, reset_password_live, html} =
      live_isolated(conn, BlyverWeb.ResetPassword.NewResetPasswordLive)

    assert html =~ "Recover your password"

    assert reset_password_live |> element("form")

    assert reset_password_live
           |> form("form", user: invalid_attrs)
           |> render_submit() =~ "can&#39;t be blank"

    form =
      reset_password_live
      |> form("form", user: user_attrs)

    assert render_submit(form) =~ ~r/phx-trigger-action/

    conn = follow_trigger_action(form, conn)
    assert redirected_to(conn, 302) =~ Routes.pow_session_path(conn, :new)
    assert get_flash(conn, :info) =~ "We have sent you an email with a recovery link"
  end
end
