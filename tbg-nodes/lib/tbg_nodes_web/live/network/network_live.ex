defmodule TbgNodesWeb.NetworkLive do
  @moduledoc false

  use TbgNodesWeb, :live_view

  alias TbgNodes.PermissionedEthereumNetworks
  alias TbgNodes.PublicEthereumNetworks

  @spec mount(any, map, Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, %{"current_user_id" => current_user_id} = _session, socket) do
    public_ethereum_networks =
      PublicEthereumNetworks.list_networks_with_interfaces_for_user(current_user_id)

    permissioned_ethereum_networks =
      PermissionedEthereumNetworks.list_networks_for_user(current_user_id)

    {
      :ok,
      assign(
        socket,
        current_user_id: current_user_id,
        public_ethereum_networks: public_ethereum_networks,
        permissioned_ethereum_networks: permissioned_ethereum_networks
      )
    }
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.NetworksView.render("index.html", assigns)
  end
end
