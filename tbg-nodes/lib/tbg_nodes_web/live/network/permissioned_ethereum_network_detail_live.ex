defmodule TbgNodesWeb.Networks.PermissionedBesuNetworkDetailLive do
  @moduledoc false

  use TbgNodesWeb, :live_view
  use Phoenix.HTML
  alias TbgNodes.PermissionedEthereumNetworks
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
        show_actions_dropdown: false,
        delete_network_error: nil,
        show_delete_modal: false,
        current_user_id: current_user_id
      )
    }
  end

  @spec handle_params(map(), any(), Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_params(%{"uuid" => uuid}, _uri, socket) do
    current_user_id = Map.get(socket.assigns, :current_user_id)

    # get network and external interfaces
    network = PermissionedEthereumNetworks.get_network_for_user_by_uuid!(current_user_id, uuid)

    external_interfaces =
      network.external_interfaces
      |> Enum.map(fn ei -> Enum.map(ei.basicauth_creds, fn cred -> {ei, cred} end) end)
      |> List.flatten()
      |> Enum.map(&get_interface_with_creds/1)

    {
      :noreply,
      assign(socket, network: network, external_interfaces: external_interfaces)
    }
  end

  @spec handle_event(String.t(), map(), Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_event("toggle_delete_modal", _, socket) do
    {:noreply, assign(socket, show_delete_modal: !socket.assigns.show_delete_modal)}
  end

  def handle_event("delete_network", %{"uuid" => uuid}, socket) do
    case PermissionedEthereumNetworks.delete_network_for_user(
           uuid,
           socket.assigns.current_user_id
         ) do
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

  @spec get_interface_with_creds({
          %PermissionedEthereumNetworks.ExternalInterface{},
          %PermissionedEthereumNetworks.BasicauthCred{}
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
    TbgNodesWeb.NetworksView.render("permissioned_ethereum_network_detail.html", assigns)
  end
end
