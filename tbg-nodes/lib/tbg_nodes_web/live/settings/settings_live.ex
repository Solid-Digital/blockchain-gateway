defmodule TbgNodesWeb.SettingsLive do
  use TbgNodesWeb, :live_view

  @moduledoc "Overview page of settings"

  @spec mount(any, map, Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, %{"current_user_id" => current_user_id} = _session, socket) do
    user = TbgNodes.Users.get_user_by_id(current_user_id)

    socket =
      socket
      |> assign(:current_user, user)
      |> assign(:current_user_id, current_user_id)

    {:ok, socket}
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.SettingsView.render("index.html", assigns)
  end
end
