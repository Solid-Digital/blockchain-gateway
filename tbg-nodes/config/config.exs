# This file is responsible for configuring your application
# and its dependencies with the aid of the Mix.Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
import Config

config :tbg_nodes,
  ecto_repos: [TbgNodes.Repo]

config :tbg_nodes, TbgNodes.NetworkMonitor,
  loop_interval: 30 * 1000,
  enabled: false,
  query_status_fn: &TbgNodes.NetworkMonitor.query_status/2

# Configures the endpoint
config :tbg_nodes, TbgNodesWeb.Endpoint,
  url: [host: "localhost"],
  secret_key_base: "xLARAipbiw2HtlG+4VSckwIBn7BxIyrTvSZ2XB6BdM5HkvW82xsPycrD8fFQpieg",
  render_errors: [view: TbgNodesWeb.ErrorView, accepts: ~w(html json)],
  pubsub_server: :pubsub,
  live_view: [
    signing_salt: "Xs69GyMtDHkbAGgTDY6HFU6rqekzCjC9"
  ]

# Configures Elixir's Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id],
  level: :debug

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

config :tbg_nodes, :telemetry, true

config :tbg_nodes, :pow,
  user: TbgNodes.Users.User,
  repo: TbgNodes.Repo,
  extensions: [PowResetPassword],
  controller_callbacks: Pow.Extension.Phoenix.ControllerCallbacks,
  web_module: TbgNodesWeb,
  mailer_backend: TbgNodesWeb.Mailer,
  cache_store_backend: TbgNodesWeb.PowRedisCache

config :tbg_nodes, TbgNodesWeb.Mailer, adapter: Bamboo.LocalAdapter

config :logger,
  backends: [:console]

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
with true <- File.exists?("./config/#{Mix.env()}.exs") do
  import_config "#{Mix.env()}.exs"
end
