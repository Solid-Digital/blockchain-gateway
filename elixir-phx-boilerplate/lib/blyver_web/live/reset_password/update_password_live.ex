defmodule BlyverWeb.ResetPassword.UpdatePasswordLive do
  @moduledoc false
  use BlyverWeb, :reset_password_live_view

  alias Blyver.Accounts.User

  @impl true
  def mount(_params, %{"action" => action}, socket) do
    {:ok,
     assign(socket,
       action: action,
       changeset: User.submit_reset_password_changeset(%User{}, %{}),
       trigger_submit: false
     )}
  end

  @impl true
  def handle_event("save", %{"user" => params}, socket) do
    changeset =
      %User{}
      |> User.submit_reset_password_changeset(params)
      |> Map.put(:action, :insert)

    case changeset.valid? do
      true ->
        {:noreply, assign(socket, changeset: changeset, trigger_submit: true)}

      false ->
        {:noreply, assign(socket, changeset: changeset)}
    end
  end
end
