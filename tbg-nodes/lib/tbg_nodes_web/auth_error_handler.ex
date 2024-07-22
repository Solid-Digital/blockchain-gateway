defmodule TbgNodesWeb.AuthErrorHandler do
  @moduledoc false

  use TbgNodesWeb, :controller
  alias Plug.Conn

  @spec call(Conn.t(), atom()) :: Conn.t()
  def call(conn, :not_authenticated) do
    conn
    |> put_flash(:error, "You've to be authenticated first")
    |> redirect(to: Routes.pow_session_path(conn, :new))
  end

  def call(conn, :already_authenticated) do
    conn
    |> put_flash(:error, "You're already authenticated")
    |> redirect(to: Routes.redirect_path(conn, :handle_redirect))
  end
end
