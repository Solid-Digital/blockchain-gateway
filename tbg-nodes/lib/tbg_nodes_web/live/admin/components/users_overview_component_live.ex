defmodule TbgNodesWeb.AdminUsersComponentLive do
  @moduledoc false

  use TbgNodesWeb, :live_view

  alias TbgNodes.Users

  @spec mount(any(), map(), Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, _session, socket) do
    users = Users.list_users_with_resources()

    {:ok, assign(socket, users: users)}
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.AdminDashboardView.render("components/users_overview.html", assigns)
  end
end
