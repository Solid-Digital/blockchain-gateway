# This file is responsible for configuring your application
# and its dependencies with the aid of the Mix.Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
use Mix.Config

config :blyver,
  ecto_repos: [Blyver.Repo]

# Configures the endpoint
config :blyver, BlyverWeb.Endpoint,
  url: [host: "localhost"],
  secret_key_base: "udUoNAyRnDin0u+g4E7B0P/Wb6IfksJUbEyrcn0Ovp3dvZlSBbh0wDj1A/8yvoqZ",
  render_errors: [view: BlyverWeb.ErrorView, accepts: ~w(html json), layout: false],
  pubsub_server: Blyver.PubSub,
  live_view: [signing_salt: "hfG2z6Q6"]

# Configures Elixir's Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

config :ex_aws,
  region: "local",
  debug_requests: true,
  json_codec: Jason

config :ex_aws, :s3,
  access_key_id: "minio",
  # sobelow_skip ["Secrets"]
  secret_access_key: "minio123",
  scheme: "http://",
  host: "localhost",
  port: 9000

config :blyver, :pow,
  user: Blyver.Accounts.User,
  repo: Blyver.Repo,
  web_module: BlyverWeb,
  mailer_backend: BlyverWeb.Pow.Mailer,
  web_mailer_module: BlyverWeb,
  routes_backend: BlyverWeb.Pow.Routes,
  extensions: [PowEmailConfirmation],
  controller_callbacks: Pow.Extension.Phoenix.ControllerCallbacks,
  messages_backend: BlyverWeb.Pow.Messages

config :blyver, BlyverWeb.Pow.Mailer, adapter: Bamboo.LocalAdapter

blyver_sender_email = System.get_env("BLYVER_SENDER_EMAIL") || "no-reply@blyver.com"
support_email = System.get_env("BLYVER_SUPPORT_EMAIL") || "support@blyver.com"

config :blyver, :emails,
  from: blyver_sender_email,
  support_email: support_email

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{Mix.env()}.exs"
