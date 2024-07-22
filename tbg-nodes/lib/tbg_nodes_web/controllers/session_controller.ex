defmodule TbgNodesWeb.SessionController do
  @moduledoc false

  use TbgNodesWeb, :controller
  require Logger

  alias TbgNodes.Users.User

  def new(conn, _params) do
    #    changeset = Pow.Plug.change_user(conn)
    changeset = User.changeset_new_login(%User{}, %{})

    render(conn, "new.html", changeset: changeset)
  end

  def create(conn, %{"user" => user_params}) do
    conn
    |> Pow.Plug.authenticate_user(user_params)
    |> case do
      {:ok, conn} ->
        conn
        |> put_flash(:info, "Welcome back!")
        |> redirect(to: Routes.redirect_path(conn, :handle_redirect))

      {:error, conn} ->
        changeset = User.changeset_login_attempt(%User{}, conn.params["user"])

        conn
        |> render("new.html", changeset: changeset)
    end
  end

  def delete(conn, _params) do
    conn
    |> Pow.Plug.delete()
    |> redirect(to: Routes.redirect_path(conn, :handle_redirect))
  end
end
