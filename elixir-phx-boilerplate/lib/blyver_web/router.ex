defmodule BlyverWeb.Router do
  use BlyverWeb, :router

  use Pow.Extension.Phoenix.Router,
    extensions: [PowEmailConfirmation, PowResetPassword]

  use Pow.Phoenix.Router

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_live_flash
    plug :put_root_layout, {BlyverWeb.LayoutView, :root}
    plug :protect_from_forgery
    nonce = 16 |> :crypto.strong_rand_bytes() |> Base.url_encode64(padding: false)
    plug :assign_csp_nonce, nonce

    plug :put_secure_browser_headers, %{
      "content-security-policy" =>
        "default-src 'self' https://*.blyver.com; style-src 'self' 'unsafe-inline' *.fontawesome.com; script-src 'self' 'unsafe-eval' cdnjs.cloudflare.com 'nonce-#{
          nonce
        }'; font-src 'self' *.fontawesome.com data: 'self'; img-src 'self' data:;"
    }
  end

  pipeline :pow_email_layout do
    plug :put_pow_mailer_layout, {BlyverWeb.LayoutView, "email.html"}
  end

  pipeline :no_layout do
    plug :put_layout, false
  end

  pipeline :protected do
    plug Pow.Plug.RequireAuthenticated,
      error_handler: Pow.Phoenix.PlugErrorHandler
  end

  defp assign_csp_nonce(conn, nonce) do
    conn
    |> assign(:csp_nonce, nonce)
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/registration", BlyverWeb do
    pipe_through :browser

    live "/new", Registration.NewRegistrationLive, :index
    live "/complete", Registration.CompleteRegistrationLive, :index
  end

  scope "/", BlyverWeb do
    pipe_through :browser

    live "/", HomeLive, :index
  end

  scope "/" do
    pipe_through [:browser, :no_layout, :pow_email_layout]

    pow_routes()
    pow_extension_routes()
  end

  scope "/", BlyverWeb do
    pipe_through [:browser, :protected]

    live "/dashboard", Dashboard.DashboardLive, :index
  end

  # Other scopes may use custom stacks.
  # scope "/api", BlyverWeb do
  #   pipe_through :api
  # end

  # Enables LiveDashboard only for development
  #
  # If you want to use the LiveDashboard in production, you should put
  # it behind authentication and allow only admins to access it.
  # If your application does not have an admins-only section yet,
  # you can use Plug.BasicAuth to set up some basic authentication
  # as long as you are also using SSL (which you should anyway).
  if Mix.env() in [:dev, :test] do
    import Phoenix.LiveDashboard.Router

    scope "/" do
      pipe_through :browser
      live_dashboard "/live-dashboard", metrics: BlyverWeb.Telemetry
    end
  end

  if Mix.env() == :dev do
    forward "/sent_emails", Bamboo.SentEmailViewerPlug
  end

  defp put_pow_mailer_layout(conn, layout), do: put_private(conn, :pow_mailer_layout, layout)
end
