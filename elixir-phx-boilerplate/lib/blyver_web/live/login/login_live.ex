defmodule BlyverWeb.Login.LoginLive do
  @moduledoc false
  use BlyverWeb, :login_live_view
  alias Blyver.Accounts.User

  @impl true
  def mount(_params, _session, socket) do
    {:ok, assign(socket, changeset: User.login_changeset(%User{}, %{}), trigger_submit: false)}
  end

  @impl true
  def handle_event("save", %{"user" => params}, socket) do
    changeset =
      %User{}
      |> User.login_changeset(params)
      |> Map.put(:action, :insert)

    case changeset.valid? do
      true ->
        {:noreply, assign(socket, changeset: changeset, trigger_submit: true)}

      false ->
        {:noreply, assign(socket, changeset: changeset)}
    end
  end
end
