defmodule BlyverWeb.ConnCase do
  @moduledoc """
  This module defines the test case to be used by
  tests that require setting up a connection.

  Such tests rely on `Phoenix.ConnTest` and also
  import other functionality to make it easier
  to build common data structures and query the data layer.

  Finally, if the test case interacts with the database,
  we enable the SQL sandbox, so changes done to the database
  are reverted at the end of every test. If you are using
  PostgreSQL, you can even run database tests asynchronously
  by setting `use BlyverWeb.ConnCase, async: true`, although
  this option is not recommended for other databases.
  """

  use ExUnit.CaseTemplate

  using do
    quote do
      # Import conveniences for testing with connections
      import Plug.Conn
      import Phoenix.ConnTest
      import BlyverWeb.ConnCase
      import Blyver.TestHelpers

      alias BlyverWeb.Router.Helpers, as: Routes

      # The default endpoint for testing
      @endpoint BlyverWeb.Endpoint
    end
  end

  setup tags do
    :ok = Ecto.Adapters.SQL.Sandbox.checkout(Blyver.Repo)

    unless tags[:async] do
      Ecto.Adapters.SQL.Sandbox.mode(Blyver.Repo, {:shared, self()})
    end

    conn =
      Phoenix.ConnTest.build_conn()
      |> Map.put(
        :secret_key_base,
        "sGEdLQ4FeJ0KhH0L67pAAsWl6RluS9zUwtCbahwFbSr97RAxYO+uPkbKeHoJ051l"
      )

    {:ok, conn: conn}
  end
end
