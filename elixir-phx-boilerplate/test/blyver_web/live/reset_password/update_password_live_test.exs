defmodule BlyverWeb.ResetPassword.UpdatePasswordLiveTest do
  use BlyverWeb.ConnCase
  import Phoenix.LiveViewTest

  alias PowResetPassword.Plug

  setup [:create_default_user]

  test "validates and submits the reset password form", %{conn: conn, user: user} do
    user_params = %{"email" => user.email}
    conn = Pow.Plug.put_config(conn, otp_app: :blyver)

    {:ok, %{token: token, user: _user}, conn} = Plug.create_reset_token(conn, user_params)

    user_attrs = %{
      password: "Apassw0rd$",
      password_confirmation: "Apassw0rd$"
    }

    invalid_attrs = %{
      password: "",
      password_confirmation: "Apassw0rd$"
    }

    url = Routes.pow_reset_password_reset_password_path(conn, :update, token)

    {:ok, reset_password_live, html} =
      live_isolated(conn, BlyverWeb.ResetPassword.UpdatePasswordLive, session: %{"action" => url})

    assert html =~ "Create a new password for your account"

    assert reset_password_live |> element("form")

    assert reset_password_live
           |> form("form", user: invalid_attrs)
           |> render_submit() =~ "can&#39;t be blank"

    form =
      reset_password_live
      |> form("form", user: user_attrs, _method: "put")

    assert render_submit(form) =~ ~r/phx-trigger-action/

    conn = follow_trigger_action(form, conn)
    assert redirected_to(conn, 302) =~ Routes.pow_session_path(conn, :new)
    assert get_flash(conn, :info) =~ "The password has been updated."
  end
end
