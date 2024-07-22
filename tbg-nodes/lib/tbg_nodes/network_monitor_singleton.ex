defmodule TbgNodes.NetworkMonitorSingleton do
  @moduledoc false
  use GenServer

  def start_link(_args) do
    GenServer.start_link(__MODULE__, %{})
  end

  def init(args) do
    {:ok, _} = Singleton.start_child(TbgNodes.NetworkMonitor, [1], {:network_monitor, 1})

    {:ok, args}
  end
end
