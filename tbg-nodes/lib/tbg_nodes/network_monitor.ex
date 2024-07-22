defmodule TbgNodes.NetworkMonitor do
  @moduledoc false

  use GenServer

  alias TbgNodes.NetworkMonitor.Status
  require Logger

  @spec get_cached_status(
          %TbgNodes.PublicEthereumNetworks.Network{}
          | %TbgNodes.PermissionedEthereumNetworks.Network{}
          | %TbgNodes.PermissionedEthereumNetworks.BesuNode{}
        ) :: Status.t()
  def get_cached_status(network_or_node) do
    GenServer.call({:global, {:network_monitor, 1}}, {:get_cached_status, network_or_node})
  end

  @spec start_link(any) :: GenServer.on_start()
  def start_link(_args) do
    GenServer.start_link(__MODULE__, %{})
  end

  @impl true
  @spec init(any()) :: {:ok, atom | :ets.tid()}
  def init(_args) do
    # Schedule work to be performed at some point
    schedule_work(1000)
    {:ok, TbgNodes.NetworkMonitorInternal.init()}
  end

  @spec get_network_monitor_loop_interval :: number()
  def get_network_monitor_loop_interval do
    Application.get_env(:tbg_nodes, TbgNodes.NetworkMonitor)[:loop_interval]
  end

  @spec schedule_work(non_neg_integer()) :: reference()
  defp schedule_work(time) do
    Process.send_after(self(), :work, time)
  end

  @impl true
  @spec handle_info(:work, any()) :: {:noreply, any()}
  def handle_info(:work, state) do
    TbgNodes.NetworkMonitorInternal.cache_all_statuses()

    schedule_work(get_network_monitor_loop_interval())
    {:noreply, state}
  end

  def handle_info(msg, state) do
    Logger.info("#{__MODULE__} - unhandled event #{inspect(msg)}")
    {:noreply, state}
  end

  @impl true
  @spec handle_call(
          {:get_cached_status,
           %TbgNodes.PublicEthereumNetworks.Network{}
           | %TbgNodes.PermissionedEthereumNetworks.Network{}
           | %TbgNodes.PermissionedEthereumNetworks.BesuNode{}},
          GenServer.from(),
          term()
        ) :: {:reply, Status.t(), any()}
  def handle_call({:get_cached_status, network_or_node}, _from, state) do
    {:reply, TbgNodes.NetworkMonitorInternal.get_cached_status(network_or_node), state}
  end
end
