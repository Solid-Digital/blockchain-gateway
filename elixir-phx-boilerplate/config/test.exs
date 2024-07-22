use Mix.Config

# Configure your database
#
# The MIX_TEST_PARTITION environment variable can be used
# to provide built-in test partitioning in CI environment.
# Run `mix help test` for more information.
config :blyver, Blyver.Repo,
  username: "postgres",
  password: "postgres",
  database: "blyver_test#{System.get_env("MIX_TEST_PARTITION")}",
  hostname: "localhost",
  pool: Ecto.Adapters.SQL.Sandbox

config :blyver, BlyverWeb.Pow.Mailer, adapter: Bamboo.TestAdapter

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :blyver, BlyverWeb.Endpoint,
  http: [port: 4002],
  server: false,
  secret_key_base: "sGEdLQ4FeJ0KhH0L67pAAsWl6RluS9zUwtCbahwFbSr97RAxYO+uPkbKeHoJ051l"

# Print only warnings and errors during test
config :logger, level: :warn
