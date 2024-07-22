defmodule TbgNodesWeb.RedirectController do
  @moduledoc false
  use TbgNodesWeb, :controller
  alias Pow.Plug

  def handle_redirect(conn, _params) do
    if is_user_logged_in(conn) do
      conn
      |> Phoenix.Controller.redirect(to: Routes.live_path(conn, TbgNodesWeb.NetworkLive))
    else
      conn
      |> Phoenix.Controller.redirect(to: Routes.pow_session_path(conn, :new))
    end
  end

  defp is_user_logged_in(conn) do
    Plug.current_user(conn)
  end
end
