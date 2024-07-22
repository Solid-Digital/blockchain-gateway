defmodule TbgNodesWeb.LibCase do
  @moduledoc """
  This module defines the test case to be used by
  tests that are util or standalone actions.

  If the test case interacts with the database,
  we enable the SQL sandbox, so changes done to the database
  are reverted at the end of every test. If you are using
  PostgreSQL, you can even run database tests asynchronously
  by setting `use TbgNodesWeb.ConnCase, async: true`, although
  this option is not recommendded for other databases.
  """

  use ExUnit.CaseTemplate

  using do
    quote do
      import TbgNodes.TestHelpers

      # The default endpoint for testing
      @endpoint TbgNodesWeb.Endpoint
    end
  end

  setup tags do
    :ok = Ecto.Adapters.SQL.Sandbox.checkout(TbgNodes.Repo)

    unless tags[:async] do
      Ecto.Adapters.SQL.Sandbox.mode(TbgNodes.Repo, {:shared, self()})
    end

    :ok
  end
end
