defmodule TbgNodes.MixProject do
  use Mix.Project

  def project do
    [
      app: :tbg_nodes,
      version: version(),
      elixir: "~> 1.5",
      elixirc_paths: elixirc_paths(Mix.env()),
      compilers: [:phoenix, :gettext] ++ Mix.compilers(),
      start_permanent: Mix.env() == :prod,
      aliases: aliases(),
      deps: deps(),
      elixirc_options: [
        warnings_as_errors: true
      ],
      aliases: aliases(),
      dialyzer: [
        plt_file: {:no_warn, "priv/plts/dialyzer.plt"},
        ignore_warnings: ".dialyzer_ignore.exs",
        flags: [:underspecs, :error_handling, :unmatched_returns, :unknown]
      ]
    ]
  end

  # Configuration for the OTP application.
  #
  # Type `mix help compile.app` for more information.
  def application do
    [
      mod: {TbgNodes.Application, []},
      # applications: [:singleton],
      extra_applications: [
        :singleton,
        :ssl,
        :logger,
        :timex,
        :runtime_tools,
        :phoenix_pubsub_redis,
        :os_mon
      ]
    ]
  end

  def version do
    case File.read("VERSION") do
      {:ok, body} ->
        to_semver(body)

      {:error, :enoent} ->
        # If the VERSION file does not exist.
        dev_version = "0.0.1-dev"
        IO.inspect("VERSION file not found, setting to do #{dev_version}")
        dev_version

      _ ->
        raise "Application version could not be set."
    end
  end

  def to_semver(raw_version) do
    raw_version
    |> String.trim()
    |> String.trim_leading("v")
  end

  # Specifies which paths to compile per environment.
  defp elixirc_paths(:test), do: ["lib", "test/support"]
  defp elixirc_paths(_), do: ["lib"]

  # Specifies your project dependencies.
  #
  # Type `mix help deps` for examples and options.
  defp deps do
    [
      {:phoenix, "== 1.5.6"},
      {:phoenix_pubsub, "~> 2.0"},
      {:phoenix_pubsub_redis, "~> 3.0"},
      {:phoenix_html, "~> 2.14.2"},
      {:phoenix_live_reload, "~> 1.2", only: :dev},
      {:phoenix_live_view, "~> 0.14.7"},
      {:phoenix_api_toolkit, "~> 0.12.0"},
      {:phoenix_ecto, "~> 4.1.0"},
      {:ecto_sql, "~> 3.4.5"},
      {:httpoison, "~> 1.6"},
      {:postgrex, ">= 0.15.7"},
      {:gettext, "~> 0.11"},
      {:jason, "~> 1.2", override: true},
      {:plug_cowboy, "~> 2.3"},
      # authn/z & user management
      {:pow, "~> 1.0.20"},
      {:plug, "~> 1.10.4"},
      {:pow_assent, "~> 0.4.9"},
      # mailer
      {:bamboo, "~> 1.4"},
      {:certifi, "~> 2.4"},
      {:ssl_verify_fun, "~> 1.1.6"},
      {:timex, "~> 3.5"},
      {:floki, ">= 0.0.0", only: :test},
      {:redix, "~> 0.10"},
      {:credo, "~> 1.5", only: [:dev, :test], runtime: false},
      {:dialyxir, "~> 1.0.0", only: [:dev, :test], runtime: false},
      {:sobelow, "~> 0.10.6", only: [:dev, :test], runtime: false},
      {:slack, "~> 0.23.5"},
      {:castore, "~> 0.1.0"},
      {:eth, "~> 0.6.0"},
      {:k8s, "~> 0.5"},
      {:libcluster, "~> 3.2"},
      {:telemetry_metrics, "~> 0.4.0"},
      {:telemetry_poller, "~> 0.4"},
      {:phoenix_live_dashboard, "~> 0.3"},
      {:singleton, "~> 1.3.0"},
      {:ecto_psql_extras, "~> 0.2"}
    ]
  end

  # Aliases are shortcuts or tasks specific to the current project.
  # For example, to create, migrate and run the seeds file at once:
  #
  #     $ mix ecto.setup
  #
  # See the documentation for `Mix` for more info on aliases.
  defp aliases do
    [
      "ecto.setup": ["ecto.create", "ecto.migrate", "run priv/repo/seeds.exs"],
      "ecto.reset": ["ecto.drop", "ecto.setup"],
      test: ["ecto.create --quiet", "ecto.migrate", "test --exclude k8s"],
      quality: [
        "format",
        "credo --strict",
        "sobelow --config",
        "dialyzer --format=short --ignore-exit-status",
        "test --cover --exclude k8s"
      ],
      "quality.ci": [
        "format --check-formatted",
        "credo --strict",
        "sobelow --config",
        "dialyzer"
      ]
    ]
  end
end
