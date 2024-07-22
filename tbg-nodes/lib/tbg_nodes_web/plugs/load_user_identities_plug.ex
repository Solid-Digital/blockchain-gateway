defmodule TbgNodesWeb.LoadUserIdentitiesPlug do
  @moduledoc false

  alias Plug.Conn

  @doc false
  @spec init(any()) :: any()
  def init(config), do: config

  @doc false
  @spec call(Conn.t(), atom() | binary() | [atom()] | [binary()]) :: Conn.t()
  def call(conn, _opts) do
    config = Pow.Plug.fetch_config(conn)

    case Pow.Plug.current_user(conn, config) do
      nil ->
        conn

      %{user_identities: %Ecto.Association.NotLoaded{}} = user ->
        user = TbgNodes.Repo.preload(user, :user_identities)
        config = Pow.Plug.fetch_config(conn)
        plug = Pow.Plug.get_plug(config)

        conn
        |> plug.do_create(user, config)

      _ ->
        conn
    end
  end
end
