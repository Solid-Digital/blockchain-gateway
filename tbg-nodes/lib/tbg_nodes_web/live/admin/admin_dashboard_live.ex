defmodule TbgNodesWeb.AdminDashboardLive do
  @moduledoc false

  use TbgNodesWeb, :live_view

  @spec mount(any(), map(), Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, %{"current_user_id" => current_user_id} = _session, socket) do
    {:ok, assign(socket, current_user_id: current_user_id, active_tab: "users")}
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.AdminDashboardView.render("index.html", assigns)
  end

  def handle_event("select_tab", %{"active_tab" => active_tab}, socket) do
    {:noreply, assign(socket, active_tab: active_tab)}
  end
end
