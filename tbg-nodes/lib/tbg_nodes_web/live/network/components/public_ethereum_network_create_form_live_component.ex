defmodule TbgNodesWeb.Networks.PublicEthereumNetworkCreateFormLiveComponent do
  @moduledoc """
  The Public Live Form Component.
  """
  use Phoenix.LiveComponent
  use Phoenix.HTML

  alias TbgNodes.PublicEthereumNetworks
  alias TbgNodesWeb.Router.Helpers, as: Routes

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    Phoenix.View.render(
      TbgNodesWeb.Networks.NewNetworkView,
      "components/public_ethereum_network_create_form.html",
      assigns
    )
  end

  @spec mount(map(), map(), any()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(
        _params,
        %{"current_user_id" => current_user_id, "network_type" => network_type},
        socket
      ) do
    socket =
      socket
      |> assign(:current_user_id, current_user_id)
      |> assign(:network_type, network_type)

    {:ok, socket}
  end

  @spec update(map(), map()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def update(%{current_user_id: current_user_id, network_type: network_type, id: id}, socket) do
    changeset =
      PublicEthereumNetworks.NetworkCreateForm.changeset(
        %PublicEthereumNetworks.NetworkCreateForm{},
        %{}
      )

    {
      :ok,
      assign(
        socket,
        current_user_id: current_user_id,
        id: id,
        changeset: changeset,
        network_type: network_type
      )
    }
  end

  @spec preload(list(map())) :: [map()]
  def preload(list_of_assigns) do
    Enum.map(
      list_of_assigns,
      fn %{id: id, network_type: network_type, current_user_id: current_user_id} ->
        %{id: id, network_type: "#{network_type}", current_user_id: current_user_id}
      end
    )
  end

  @spec handle_event(String.t(), map(), any()) :: {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_event("validate", %{"network_create_form" => params}, socket) do
    changeset =
      %PublicEthereumNetworks.NetworkCreateForm{}
      |> PublicEthereumNetworks.NetworkCreateForm.changeset(params)
      |> Map.put(:action, :insert)

    {:noreply, assign(socket, changeset: changeset)}
  end

  def handle_event("create_network", %{"network_create_form" => params}, socket) do
    params = maybe_archive_data(params)

    case PublicEthereumNetworks.create_network_with_interfaces(
           socket.assigns.current_user_id,
           params
         ) do
      {:ok, network} ->
        {
          :noreply,
          socket
          |> put_flash(:info, "Network created successfully.")
          |> push_redirect(
            to:
              Routes.live_path(
                socket,
                TbgNodesWeb.Networks.PublicEthereumNetworkDetailLive,
                network.uuid
              )
          )
        }

      {:error, %Ecto.Changeset{} = changeset} ->
        {:noreply, assign(socket, changeset: changeset)}
    end
  end

  @spec maybe_archive_data(map()) :: map()
  def maybe_archive_data(params) do
    cond do
      Map.get(params, "archive_data") == nil ->
        params

      String.to_existing_atom(Map.get(params, "archive_data")) == true ->
        # append _archive to network_configuration in case of archive data selected
        network_configuration = Map.get(params, "network_configuration") <> "-archive"
        Map.put(params, "network_configuration", network_configuration)

      true ->
        params
    end
  end
end
