defmodule Blyver.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  def start(_type, _args) do
    children = [
      # Start the Ecto repository
      Blyver.Repo,
      # Start the Telemetry supervisor
      BlyverWeb.Telemetry,
      # Start the PubSub system
      {Phoenix.PubSub, name: Blyver.PubSub},
      # Start the Endpoint (http/https)
      BlyverWeb.Endpoint
      # Start a worker by calling: Blyver.Worker.start_link(arg)
      # {Blyver.Worker, arg}
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Blyver.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  def config_change(changed, _new, removed) do
    BlyverWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
