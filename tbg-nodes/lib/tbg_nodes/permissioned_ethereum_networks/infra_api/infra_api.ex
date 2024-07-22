defmodule TbgNodes.PermissionedEthereumNetworks.InfraAPI do
  @moduledoc """
  This module defines the behaviours for InfraAPI modules. All
  calls to interface with infrastructure create/modified/monitor/destroy
  permissioned nodes related resources should be handled by an InfraAPI module.
  """

  alias TbgNodes.PermissionedEthereumNetworks

  @spec get_infra_api! :: TbgNodes.PermissionedEthereumNetworks.InfraAPI
  def get_infra_api! do
    Application.get_env(:tbg_nodes, TbgNodes.PermissionedEthereumNetworks)[:infra_api] ||
      raise "No config value found for :infra_api."
  end

  @callback delete_network(%PermissionedEthereumNetworks.Network{}) ::
              {:ok} | {:error, String.t()}
  @callback deploy_network(%PermissionedEthereumNetworks.Network{}) ::
              {:ok} | {:error, String.t()}

  @callback deploy_external_interface(%PermissionedEthereumNetworks.ExternalInterface{}) ::
              {:ok, {:url, String.t()}} | {:error, String.t()}

  @callback get_liveness_url(%PermissionedEthereumNetworks.BesuNode{}) ::
              {:ok, String.t()} | {:error, String.t()}
  @callback get_readiness_url(%PermissionedEthereumNetworks.BesuNode{}) ::
              {:ok, String.t()} | {:error, String.t()}
  @callback get_liveness_url(%PermissionedEthereumNetworks.Network{}) ::
              {:ok, String.t()} | {:error, String.t()}
  @callback get_readiness_url(%PermissionedEthereumNetworks.Network{}) ::
              {:ok, String.t()} | {:error, String.t()}
end
