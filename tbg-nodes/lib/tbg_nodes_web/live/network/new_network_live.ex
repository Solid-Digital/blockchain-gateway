defmodule TbgNodesWeb.NewNetworkLive do
  @moduledoc """
  The New Network context.
  """

  use TbgNodesWeb, :live_view

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    Phoenix.View.render(
      TbgNodesWeb.Networks.NewNetworkView,
      "components/new_network.html",
      assigns
    )
  end

  @spec mount(any(), map(), Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, %{"current_user_id" => current_user_id} = _session, socket) do
    socket =
      socket
      |> assign(:network_type, "ethereum")
      |> assign(:current_user_id, current_user_id)

    {:ok, socket}
  end

  @spec handle_event(String.t(), map(), Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_event("select_network_type", %{"network_type" => network_type}, socket) do
    if network_type == "besu" do
      {:noreply, assign(socket, network_type: network_type)}
    else
      {:noreply, assign(socket, network_type: network_type)}
    end
  end

  # The submit form in the permissioned besu network form sends a msg to self
  # that is received on the parent liveview with this handle_info.
  # In order to keep the logic contained in the permissioned besu form components
  # this function calls that component.
  def handle_info({:create_permissioned_besu_network, %{valid?: true} = changeset}, socket) do
    PermissionedEthereumNetworkCreateFormLiveComponent.create_network(changeset, socket)
  end
end
