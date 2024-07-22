defmodule TbgNodesWeb.HeaderComponentLive do
  @moduledoc false

  use TbgNodesWeb, :live_view

  #  user Phoenix.HTML
  alias TbgNodes.Users

  @spec mount(any(), map(), Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, %{"current_user_id" => current_user_id, "path" => path} = _session, socket) do
    current_user = Users.get_user_by_id(current_user_id)

    {
      :ok,
      socket
      |> assign(:path, path)
      |> assign(:current_user_id, current_user_id)
      |> assign(:current_user, current_user)
    }
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.LayoutView.render("components/header_live.html", assigns)
  end
end
