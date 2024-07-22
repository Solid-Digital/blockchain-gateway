defmodule TbgNodesWeb.Networks.PublicEthereumNetworkDetailLive do
  @moduledoc """
  The Public Live Detail Component.
  """

  use TbgNodesWeb, :live_view
  alias TbgNodes.PublicEthereumNetworks
  alias TbgNodesWeb.Router.Helpers, as: Routes

  defmodule InterfaceWithCreds do
    @moduledoc false
    defstruct [
      :id,
      :password,
      :protocol,
      :url,
      :username
    ]
  end

  @spec mount(any, map, Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, %{"current_user_id" => current_user_id} = _session, socket) do
    {
      :ok,
      assign(
        socket,
        delete_network_error: nil,
        show_delete_modal: false,
        show_actions_dropdown: false,
        current_user_id: current_user_id
      )
    }
  end

  @spec handle_params(map(), any(), Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_params(%{"uuid" => network_uuid}, _uri, socket) do
    current_user_id = Map.get(socket.assigns, :current_user_id)

    # get network and external interfaces
    network =
      PublicEthereumNetworks.get_network_with_interfaces_for_user_by_uuid!(
        current_user_id,
        network_uuid
      )

    external_interfaces =
      network.network_external_interfaces
      |> Enum.map(fn ei -> Enum.map(ei.basicauth_creds, fn cred -> {ei, cred} end) end)
      |> List.flatten()
      |> Enum.map(&get_interface_with_creds/1)

    {
      :noreply,
      assign(socket, network: network, network_external_interfaces: external_interfaces)
    }
  end

  @spec get_interface_with_creds({
          %PublicEthereumNetworks.NetworkExternalInterface{},
          %PublicEthereumNetworks.BasicauthCred{}
        }) ::
          %InterfaceWithCreds{}
  defp get_interface_with_creds({ei, cred}) do
    %InterfaceWithCreds{
      id: cred.id,
      password: cred.password,
      protocol: ei.protocol,
      url: ei.url,
      username: cred.username
    }
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.NetworksView.render("public_ethereum_network_detail.html", assigns)
  end

  @spec handle_event(String.t(), map(), Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_event("toggle_delete_modal", _, socket) do
    {:noreply, assign(socket, show_delete_modal: !socket.assigns.show_delete_modal)}
  end

  def handle_event("delete_network", %{"uuid" => uuid}, socket) do
    case PublicEthereumNetworks.delete_network(uuid, socket.assigns.current_user_id) do
      {:ok, _} ->
        # redirect to overview
        {
          :noreply,
          socket
          |> put_flash(:info, "Network deleted successfully.")
          |> push_redirect(
            to:
              Routes.live_path(
                socket,
                TbgNodesWeb.NetworkLive
              )
          )
        }

      {:error, msg} ->
        # render with error
        {:noreply, assign(socket, delete_network_error: msg)}
    end
  end
end
