ExUnit.start()
Ecto.Adapters.SQL.Sandbox.mode(Blyver.Repo, :manual)

defmodule Blyver.TestHelpers do
  def valid_user_attrs do
    %{
      email: "user-#{:rand.uniform(1_000_000)}@test.com",
      password: "SomeValidPa$$word",
      confirm_password: "SomeValidPa$$word"
    }
  end

  def user_fixture(attrs \\ %{}) do
    {:ok, user} =
      valid_user_attrs()
      |> Enum.into(attrs)
      |> Pow.Operations.create(otp_app: :blyver)

    user
  end

  def create_default_user(%{}) do
    {:ok, user: user_fixture()}
  end
end
