defmodule TbgNodes.PermissionedEthereumNetworks.InfraAPIMock do
  @moduledoc false

  alias TbgNodes.PermissionedEthereumNetworks

  @behaviour PermissionedEthereumNetworks.InfraAPI

  @impl PermissionedEthereumNetworks.InfraAPI
  def delete_network(%PermissionedEthereumNetworks.Network{} = _network) do
    {:ok}
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def deploy_network(%PermissionedEthereumNetworks.Network{} = _network) do
    {:ok}
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def deploy_external_interface(
        %PermissionedEthereumNetworks.ExternalInterface{} = external_interface
      ) do
    url_base = "this.is.a.fake.url/#{external_interface.uuid}"

    url_protocol =
      case external_interface.protocol do
        "http" -> "https"
        "websocket" -> "wss"
      end

    full_url = "#{url_protocol}://#{url_base}/#{external_interface.uuid}/"
    {:ok, {:url, full_url}}
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def get_liveness_url(%PermissionedEthereumNetworks.BesuNode{} = node) do
    TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s.get_liveness_url(node)
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def get_liveness_url(%PermissionedEthereumNetworks.Network{} = network) do
    PermissionedEthereumNetworks.InfraAPIK8s.get_liveness_url(network)
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def get_readiness_url(%PermissionedEthereumNetworks.BesuNode{} = node) do
    PermissionedEthereumNetworks.InfraAPIK8s.get_readiness_url(node)
  end

  @impl PermissionedEthereumNetworks.InfraAPI
  def get_readiness_url(%PermissionedEthereumNetworks.Network{} = network) do
    PermissionedEthereumNetworks.InfraAPIK8s.get_readiness_url(network)
  end
end
