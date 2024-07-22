defmodule TbgNodes.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  def start(_type, _args) do
    {:ok, _, _} = TbgNodes.Release.migrate()

    redis_socket_opts_provider =
      Application.get_env(:tbg_nodes, TbgNodesWeb.Endpoint)[:redis_socket_opts_provider] ||
        TbgNodes.Redis.ConfigDefault

    # List all child processes to be supervised
    children =
      [
        # Start the Ecto repository
        TbgNodes.Repo,
        TbgNodesWeb.Telemetry,
        # Start the endpoint when the application starts
        TbgNodesWeb.Endpoint
      ] ++
        network_monitor() ++
        live_tester() ++
        pubsub(redis_socket_opts_provider.socket_opts) ++
        redix(redis_socket_opts_provider.socket_opts) ++
        libcluster()

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: TbgNodes.Supervisor]
    Supervisor.start_link(children, opts)
  end

  def network_monitor do
    if Application.get_env(:tbg_nodes, TbgNodes.NetworkMonitor)[:enabled] do
      [TbgNodes.NetworkMonitorSingleton]
    else
      []
    end
  end

  def live_tester do
    if Application.get_env(:tbg_nodes, TbgNodes.LiveTester)[:enabled] do
      [TbgNodes.LiveTesterSingleton]
    else
      []
    end
  end

  defp pubsub(socket_opts) do
    [
      {
        Phoenix.PubSub,
        [
          name: :pubsub,
          adapter: Phoenix.PubSub.Redis,
          url:
            Application.get_env(:tbg_nodes, TbgNodesWeb.Endpoint)[:redis_url] ||
              "redis://localhost:6379",
          ssl: Application.get_env(:tbg_nodes, TbgNodesWeb.Endpoint)[:redis_ssl] || false,
          node_name: System.get_env("NODE") || "node@127.0.0.1",
          socket_opts: socket_opts
        ]
      }
    ]
  end

  defp redix(socket_opts) do
    [
      {
        Redix,
        {
          Application.get_env(:tbg_nodes, TbgNodesWeb.Endpoint)[:redis_url] ||
            "redis://localhost:6379",
          [
            name: :redix,
            ssl: Application.get_env(:tbg_nodes, TbgNodesWeb.Endpoint)[:redis_ssl] || false,
            socket_opts: socket_opts
          ]
        }
      }
    ]
  end

  defp libcluster do
    libcluster_topologies = Application.get_env(:libcluster, :topologies)

    if libcluster_topologies != nil do
      [{Cluster.Supervisor, [libcluster_topologies, [name: TbgNodes.ClusterSupervisor]]}]
    else
      []
    end
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  def config_change(changed, _new, removed) do
    TbgNodesWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
