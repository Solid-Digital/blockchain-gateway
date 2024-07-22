defmodule TbgNodesWeb.AdminNodesComponentLive do
  @moduledoc false

  use TbgNodesWeb, :live_view
  alias TbgNodes.PermissionedEthereumNetworks

  @spec mount(any(), map(), Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, _session, socket) do
    nodes = PermissionedEthereumNetworks.list_nodes()

    {:ok, assign(socket, nodes: nodes)}
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.AdminDashboardView.render("components/nodes_overview.html", assigns)
  end
end
