defmodule TbgNodes.Users do
  @moduledoc false
  require Ecto.Query

  import Ecto.Query, only: [from: 2]

  alias TbgNodes.{Repo, Users.NpsFeedback, Users.User}

  @type t :: %User{}
  @day_seconds 24 * 3600

  def get_user_by_email!(email, repo \\ TbgNodes.Repo), do: repo.get_by!(User, email: email)

  def list_users, do: Repo.all(User)

  def list_users_with_resources do
    query =
      from u in User,
        preload: [
          :public_ethereum_networks,
          permissioned_ethereum_networks: [:besu_nodes]
        ]

    Repo.all(query)
  end

  def get_user_by_id(user_id), do: Repo.get!(User, user_id)

  def update_email(user, params) do
    Repo.get!(User, user.id)
    |> User.changeset_update_email(params)
    |> Repo.update()
  end

  def change_password(user, params) do
    Repo.get!(User, user.id)
    |> User.changeset_change_password(params)
    |> Repo.update()
  end

  @spec create_admin(map()) :: {:ok, t()} | {:error, Ecto.Changeset.t()}
  def create_admin(params) do
    %User{}
    |> User.changeset(params)
    |> User.changeset_role(%{role: "admin"})
    |> Repo.insert()
  end

  def set_admin_role(user_email, repo \\ TbgNodes.Repo) do
    user_email
    |> get_user_by_email!(repo)
    |> User.changeset_role(%{role: "admin"})
    |> repo.update()
  end

  def user_exists(user_changeset) do
    user = Ecto.Changeset.apply_changes(user_changeset)
    query = Ecto.Query.from(u in User, where: u.email == ^user.email)

    Repo.exists?(query)
    |> case do
      true ->
        Ecto.Changeset.add_error(user_changeset, :email, "user already exists")

      false ->
        user_changeset
    end
  end

  def add_nps_feedback(user_id, score) do
    %NpsFeedback{}
    |> NpsFeedback.changeset(%{:user_id => user_id, :score => score})
    |> Repo.insert()
  end

  def get_nps_feedback_by_user_id(user_id) do
    Repo.get_by!(NpsFeedback, user_id: user_id)
  end

  @doc """
  Business logic on whether or not to request user for NPS feedback.
  Returns boolean
  """
  def show_nps_feedback(user) do
    # Retrieve most recent NPS Feedback record
    q =
      from n in NpsFeedback,
        where: n.user_id == ^user.id,
        order_by: [
          desc: n.id
        ],
        limit: 1

    last_user_feedback = Repo.one(q)

    now = DateTime.utc_now()

    cond do
      # user has not been around for more than a day
      DateTime.diff(now, user.inserted_at, :second) < @day_seconds ->
        false

      # user has not submitted or dismissed nps feedback before
      is_nil(last_user_feedback) ->
        true

      # user has not answered an NPS form in the past 100 days
      DateTime.diff(
        last_user_feedback.inserted_at,
        now,
        :second
      ) > 100 * @day_seconds ->
        true

      # user has not dismissed an NPS form in the past 20 days
      last_user_feedback.score == 0 and
          DateTime.diff(
            now,
            last_user_feedback.inserted_at,
            :second
          ) > 20 * @day_seconds ->
        true

      # fallback clause, if none of the above apply, return false
      true ->
        false
    end
  end
end
