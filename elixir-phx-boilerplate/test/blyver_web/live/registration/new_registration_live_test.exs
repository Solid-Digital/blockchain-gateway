defmodule BlyverWeb.Registration.NewRegistrationLiveTest do
  use BlyverWeb.ConnCase
  import Phoenix.LiveViewTest

  test "validates and submits registration form", %{conn: conn} do
    user_attrs = %{
      email: "test@email.com",
      password: "ThePassword$",
      password_confirmation: "ThePassword$"
    }

    invalid_attrs = %{email: "test@email.com"}

    {:ok, registration_live, _html} = live(conn, Routes.new_registration_path(conn, :index))

    assert registration_live |> element("form")

    assert registration_live
           |> form("form", user: invalid_attrs)
           |> render_submit() =~ "can&#39;t be blank"

    form =
      registration_live
      |> form("form", user: user_attrs)

    assert render_submit(form) =~ ~r/phx-trigger-action/

    conn = follow_trigger_action(form, conn)
    assert redirected_to(conn, 302) =~ Routes.complete_registration_path(conn, :index)
  end
end
