defmodule TbgNodesWeb.Networks.PublicEthereumNetworkStatusLive do
  @moduledoc false

  use Phoenix.LiveView
  use Phoenix.HTML

  @spec mount(any, map, Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(params, %{"current_user_id" => current_user_id} = session, socket) do
    _ =
      if connected?(socket),
        do:
          :timer.send_interval(
            TbgNodes.NetworkMonitor.get_network_monitor_loop_interval(),
            self(),
            :update
          )

    network_uuid = (is_map(params) && params["uuid"]) || session["network_uuid"]

    status = TbgNodes.PublicEthereumNetworks.get_status_for_user(network_uuid, current_user_id)

    {:ok, assign(socket, status: status, user_id: current_user_id, network_uuid: network_uuid)}
  end

  @spec handle_info(:update, Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_info(:update, socket) do
    status =
      TbgNodes.PublicEthereumNetworks.get_status_for_user(
        socket.assigns.network_uuid,
        socket.assigns.user_id
      )

    {:noreply, assign(socket, status: status)}
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.NetworksView.render("components/healthz_status.html", assigns)
  end
end
