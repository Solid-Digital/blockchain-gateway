defmodule TbgNodesWeb.Router do
  use TbgNodesWeb, :router
  use Pow.Phoenix.Router
  use PowAssent.Phoenix.Router

  use Pow.Extension.Phoenix.Router,
    extensions: [PowResetPassword]

  import Phoenix.LiveView.Router
  import Phoenix.LiveDashboard.Router
  import Phoenix.Controller

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_live_flash
    plug :protect_from_forgery

    plug :put_root_layout, {TbgNodesWeb.LayoutView, :root}

    nonce = 16 |> :crypto.strong_rand_bytes() |> Base.url_encode64(padding: false)
    plug :assign_csp_nonce, nonce

    plug :put_secure_browser_headers, %{
      "content-security-policy" =>
        "default-src 'self' https://*.unchain.io; style-src 'self' 'unsafe-inline' *.fontawesome.com; script-src 'self' cdnjs.cloudflare.com 'nonce-#{
          nonce
        }' 'unsafe-eval'; font-src 'self' data: 'self' *.fontawesome.com; img-src 'self' data:;"
    }

    plug PowAssent.Plug.Reauthorization,
      handler: PowAssent.Phoenix.ReauthorizationPlugHandler

    plug TbgNodesWeb.ShowNpsFeedbackPlug
  end

  defp assign_csp_nonce(conn, nonce) do
    conn
    |> assign(:csp_nonce, nonce)
  end

  # sobelow_skip ["Config.CSRF"]
  pipeline :skip_csrf_protection do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_flash

    plug :put_secure_browser_headers, %{
      "content-security-policy" => "default-src 'self' https://*.github.com;"
    }
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  pipeline :protected do
    plug Pow.Plug.RequireAuthenticated,
      error_handler: TbgNodesWeb.ProtectedRoutesPlug

    plug TbgNodesWeb.LoadUserIdentitiesPlug
  end

  pipeline :admin do
    plug TbgNodesWeb.EnsureRolePlug, :admin
  end

  scope "/" do
    pipe_through :skip_csrf_protection

    pow_assent_authorization_post_callback_routes()
  end

  scope "/" do
    pipe_through :browser

    pow_extension_routes()
    pow_assent_routes()
  end

  pipeline :auth do
    plug :put_layout, {TbgNodesWeb.LayoutView, "auth.html"}
  end

  scope "/", TbgNodesWeb do
    pipe_through [:browser, :auth]

    get "/", RedirectController, :handle_redirect

    get "/login", SessionController, :new, as: :pow_session
    post "/login", SessionController, :create, as: :pow_session

    get "/signup", RegistrationController, :step_1
    post "/signup/step-1", RegistrationController, :submit_step_1
    post "/signup/step-2", RegistrationController, :submit_step_2
  end

  scope "/", TbgNodesWeb do
    pipe_through [:browser, :protected]

    # Add protected routes here
    live "/networks",
         NetworkLive,
         session: {__MODULE__, :with_session, []}

    live "/networks/new",
         NewNetworkLive,
         session: {__MODULE__, :with_session, []}

    live "/networks/besu/:uuid",
         Networks.PermissionedBesuNetworkDetailLive,
         session: {__MODULE__, :with_session, []}

    live "/networks/public-ethereum/:uuid",
         Networks.PublicEthereumNetworkDetailLive,
         session: {__MODULE__, :with_session, []}

    get "/profile/password", UserController, :password_index
    get "/profile/email", UserController, :email_index
    put "/profile/password", UserController, :change_password
    put "/profile/email", UserController, :update_email

    live "/settings",
         SettingsLive,
         session: {__MODULE__, :with_session, []}

    delete "/logout", SessionController, :delete, as: :pow_session

    def with_session(conn) do
      %{
        "current_user_id" => conn.assigns.current_user.id,
        "path" => conn.path_info
      }
    end
  end

  scope "/admin", TbgNodesWeb do
    pipe_through [:browser, :protected, :admin]

    live "/",
         AdminDashboardLive,
         session: {__MODULE__, :with_session, []}

    live_dashboard "/dashboard",
      metrics: TbgNodesWeb.Telemetry,
      csp_nonce_assign_key: :csp_nonce,
      ecto_repos: [TbgNodes.Repo]

    # metrics_history: {MyStorage, :metrics_history, []}
  end

  # Other scopes may use custom stacks.
  # scope "/api", TbgNodesWeb do
  #   pipe_through :api
  # end

  if Application.get_env(:tbg_nodes, :env) === :dev do
    forward "/sent_emails", Bamboo.SentEmailViewerPlug
  end
end
