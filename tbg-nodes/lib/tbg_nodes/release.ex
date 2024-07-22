defmodule TbgNodes.Release do
  @moduledoc false
  alias TbgNodes.Repo
  require Logger

  @app :tbg_nodes

  def migrate do
    if ensure_loaded(@app) do
      {:ok, _} = Application.ensure_all_started(:ssl)
      {:ok, _, _} = Ecto.Migrator.with_repo(Repo, &Ecto.Migrator.run(&1, :up, all: true))
    else
      {:error, :not_loaded}
    end
  end

  def ensure_loaded(app) do
    case Application.load(app) do
      :ok ->
        true

      {:error, {:already_loaded, ^app}} ->
        true

      {:error, _} ->
        false
    end
  end

  def rollback(version) do
    :ok = Application.load(@app)

    {:ok, _} = Application.ensure_all_started(:ssl)
    {:ok, _, _} = Ecto.Migrator.with_repo(Repo, &Ecto.Migrator.run(&1, :down, to: version))
  end

  def rollback_once do
    :ok = Application.load(@app)

    {:ok, _} = Application.ensure_all_started(:ssl)
    {:ok, _, _} = Ecto.Migrator.with_repo(Repo, &Ecto.Migrator.run(&1, :down, step: 1))
  end

  def set_admin_role(user_email) do
    _ = Application.load(@app)

    [:ssl, :postgrex, :ecto]
    |> Enum.each(fn app ->
      {:ok, _} = Application.ensure_all_started(app)
    end)

    res =
      Ecto.Migrator.with_repo(Repo, fn repo ->
        TbgNodes.Users.set_admin_role(user_email, repo)
      end)

    Logger.info("User #{user_email} is now admin")

    res
  end
end
