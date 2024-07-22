defmodule TbgNodes.LiveTestHelpers do
  @moduledoc false

  import Ecto.Repo

  require Logger

  @spec check_until(any, keyword) :: any
  def check_until(fun, opts \\ []) do
    opts = Keyword.merge([max_iterations: 120, interval: 1000, label: ""], opts)

    Enum.reduce_while(1..opts[:max_iterations], false, fn i, _ ->
      Logger.info("#{opts[:label]} iteration #{i}")
      :timer.sleep(opts[:interval])

      if fun.() do
        {:halt, true}
      else
        {:cont, false}
      end
    end)
  end

  @spec user_fixture(map()) :: %TbgNodes.Users.User{}
  def user_fixture(attrs \\ %{}) do
    attrs =
      attrs
      |> Enum.into(%{
        email: "user#{:rand.uniform(1_000_000_000_000_000)}@test.com",
        password: "supersecret_password"
      })

    {:ok, user} =
      case attrs |> Pow.Operations.create(otp_app: :tbg_nodes) do
        {:ok, user} ->
          {:ok, user}

        {:error, %Ecto.Changeset{errors: [email: {_, [constraint: :unique, constraint_name: _]}]}} ->
          {:ok, _} =
            TbgNodes.Users.get_user_by_email!(attrs[:email])
            |> Pow.Ecto.Schema.Changeset.new_password_changeset(attrs, otp_app: :tbg_nodes)
            |> Pow.Ecto.Context.do_update(otp_app: :tbg_nodes)
      end

    user
  end

  def permissioned_ethereum_network_user_input(user_input) do
    %TbgNodes.PermissionedEthereumNetworks.NetworkUserInput{}
    |> TbgNodes.PermissionedEthereumNetworks.NetworkUserInput.changeset(user_input)
    |> Ecto.Changeset.apply_changes()
  end
end
