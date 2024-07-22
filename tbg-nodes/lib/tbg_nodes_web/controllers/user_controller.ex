defmodule TbgNodesWeb.UserController do
  @moduledoc false

  use TbgNodesWeb, :controller
  alias TbgNodes.Users
  alias TbgNodes.Users.User

  def password_index(conn, _params) do
    user = Pow.Plug.current_user(conn)
    changeset = User.changeset_edit(user)

    render(conn, "change-password.html", changeset: changeset)
  end

  def email_index(conn, _params) do
    user = Pow.Plug.current_user(conn)
    changeset = User.changeset_edit(user)

    render(conn, "update-email.html", changeset: changeset)
  end

  def update_email(conn, params) do
    user = Pow.Plug.current_user(conn)

    case Users.update_email(user, params) do
      {:ok, _} ->
        conn
        |> put_flash(:info, "Email updated successfully")
        |> Pow.Plug.delete()
        |> redirect(to: Routes.redirect_path(conn, :handle_redirect))

      {:error, %Ecto.Changeset{} = changeset} ->
        render(conn, "update-email.html", changeset: changeset)
    end
  end

  def change_password(conn, params) do
    user = Pow.Plug.current_user(conn)

    case Users.change_password(user, params) do
      {:ok, _} ->
        conn
        |> put_flash(:info, "Password changed successfuly")
        |> redirect(to: Routes.live_path(conn, TbgNodesWeb.NetworkLive))

      {:error, %Ecto.Changeset{} = changeset} ->
        render(conn, "change-password.html", changeset: changeset)
    end
  end
end
