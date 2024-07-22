defmodule BlyverWeb.Pow.Routes do
  @moduledoc false

  use Pow.Phoenix.Routes
  alias BlyverWeb.Router.Helpers, as: Routes

  def after_registration_path(conn), do: Routes.complete_registration_path(conn, :index)

  def after_sign_in_path(conn), do: Routes.dashboard_path(conn, :index)
end
