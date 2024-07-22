defmodule Blyver.Accounts do
  @moduledoc """
  The accounts context
  """
  alias Blyver.Accounts.User
  alias Blyver.Repo
  import Ecto.Changeset, only: [add_error: 3]
  import Ecto.Query, only: [from: 2]

  @doc """
  Creates a user.
  ## Examples
      iex> create_user(%{field: value})
      {:ok, %User{}}
      iex> create_user(%{field: bad_value})
      {:error, %Ecto.Changeset{}}
  """
  def create_user(attrs \\ %{}) do
    %User{}
    |> User.changeset(attrs)
    |> Repo.insert()
  end

  @doc """
  returns an Ecto.Changeset{} for tracking user changes

  ## Examples

    iex> change_user(user)
    %Ecto.Changeset{data: %User{}}
  """
  def change_user(%User{} = user, attrs \\ %{}) do
    User.changeset(user, attrs)
  end

  @doc """
  returns an Ecto.Changeset{} with an error added if
  the email in the changeset already exists.

  It only performs this check if there are no pre-existing validation errors.
  """
  def validate_new_user(changeset) do
    email = Map.get(changeset.changes, :email)

    case changeset.valid? do
      true ->
        case user_exists?(email) do
          true -> add_error(changeset, :email, "User already exists")
          false -> changeset
        end

      false ->
        changeset
    end
  end

  defp user_exists?(email) do
    query = from u in User, where: u.email == ^email
    Repo.exists?(query)
  end
end
