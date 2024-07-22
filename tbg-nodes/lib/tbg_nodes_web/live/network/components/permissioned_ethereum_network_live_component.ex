defmodule PermissionedEthereumNetworkCreateFormLiveComponent do
  @moduledoc """
  The Permissioned Live Form Component.
  """
  import Ecto.Changeset
  use Phoenix.LiveComponent
  use Phoenix.HTML

  alias TbgNodes.PermissionedEthereumNetworks
  alias TbgNodesWeb.Router.Helpers, as: Routes

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    Phoenix.View.render(
      TbgNodesWeb.Networks.NewNetworkView,
      "components/permissioned_ethereum_network_create_form.html",
      assigns
    )
  end

  @spec mount(Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(socket) do
    initial_values = %{
      :number_besu_validators => 1,
      :number_besu_normal_nodes => 1,
      :number_besu_boot_nodes => 1,
      :deployment_option => "cloud",
      :consensus => "IBFT"
    }

    changeset =
      %PermissionedEthereumNetworks.NetworkUserInput{}
      |> PermissionedEthereumNetworks.NetworkUserInput.changeset(initial_values)

    {
      :ok,
      assign(
        socket,
        changeset: changeset,
        creating_network: false
      )
    }
  end

  @spec handle_event(String.t(), any(), Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_event("increase_besu_nodes", _, socket) do
    nodes = socket.assigns.changeset.changes.number_besu_normal_nodes

    changeset =
      PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
        :number_besu_normal_nodes => nodes + 1
      })

    changeset =
      if valid_node_change?(changeset) == false do
        PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
          :number_besu_normal_nodes => nodes
        })
      else
        changeset
      end

    {:noreply, assign(socket, changeset: changeset)}
  end

  def handle_event("decrease_besu_nodes", _, socket) do
    nodes = socket.assigns.changeset.changes.number_besu_normal_nodes

    changeset =
      PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
        :number_besu_normal_nodes => nodes - 1
      })

    changeset =
      if valid_node_change?(changeset) == false do
        PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
          :number_besu_normal_nodes => nodes
        })
      else
        changeset
      end

    {:noreply, assign(socket, changeset: changeset)}
  end

  def handle_event("increase_besu_validators", _, socket) do
    validators = socket.assigns.changeset.changes.number_besu_validators

    changeset =
      PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
        :number_besu_validators => validators + 1
      })

    changeset =
      if valid_validator_change?(changeset) == false do
        PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
          :number_besu_validators => validators
        })
      else
        changeset
      end

    {:noreply, assign(socket, changeset: changeset)}
  end

  def handle_event("decrease_besu_validators", _, socket) do
    validators = socket.assigns.changeset.changes.number_besu_validators

    changeset =
      PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
        :number_besu_validators => validators - 1
      })

    changeset =
      if valid_validator_change?(changeset) == false do
        PermissionedEthereumNetworks.NetworkUserInput.changeset(socket.assigns.changeset, %{
          :number_besu_validators => validators
        })
      else
        changeset
      end

    {:noreply, assign(socket, changeset: changeset)}
  end

  def handle_event("validate", %{"network_user_input" => params}, socket) do
    changeset =
      %PermissionedEthereumNetworks.NetworkUserInput{}
      |> PermissionedEthereumNetworks.NetworkUserInput.changeset(params)

    {:noreply, assign(socket, changeset: %{changeset | action: :insert})}
  end

  def handle_event("create_network", %{"network_user_input" => params}, socket) do
    changeset =
      %PermissionedEthereumNetworks.NetworkUserInput{}
      |> PermissionedEthereumNetworks.NetworkUserInput.changeset(params)

    if changeset.valid? do
      send(self(), {:create_permissioned_besu_network, changeset})
      {:noreply, assign(socket, creating_network: true, changeset: changeset)}
    else
      {:noreply, assign(socket, changeset: changeset)}
    end
  end

  # This function is called by the new_network_live parent liveview.
  @spec create_network(Ecto.Changeset.t(), %Phoenix.LiveView.Socket{}) ::
          {:noreply, %Phoenix.LiveView.Socket{}}
  def create_network(%{valid?: true} = changeset, socket) do
    user_id = socket.assigns.current_user_id
    user_input = apply_changes(changeset)

    case PermissionedEthereumNetworks.create_network(user_input, user_id) do
      {:ok, network} ->
        {:noreply,
         socket
         |> put_flash(:info, "Network created successfully.")
         |> push_redirect(
           to:
             Routes.live_path(
               socket,
               TbgNodesWeb.Networks.PermissionedBesuNetworkDetailLive,
               network.uuid
             )
         )}

      {:error, _err_msg} ->
        {:noreply, assign(socket, network_type: "besu", changeset: changeset)}
    end
  end

  @spec valid_node_change?(Ecto.Changeset.t()) :: bool
  defp valid_node_change?(%Ecto.Changeset{} = changeset) do
    case changeset.errors[:number_besu_normal_nodes] do
      {_, [validation: :number, kind: :less_than_or_equal_to, number: _]} -> false
      {_, [validation: :number, kind: :greater_than_or_equal_to, number: _]} -> false
      _ -> true
    end
  end

  @spec valid_validator_change?(Ecto.Changeset.t()) :: bool
  defp valid_validator_change?(%Ecto.Changeset{} = changeset) do
    case changeset.errors[:number_besu_validators] do
      {_, [validation: :number, kind: :less_than_or_equal_to, number: _]} -> false
      {_, [validation: :number, kind: :greater_than_or_equal_to, number: _]} -> false
      _ -> true
    end
  end
end
