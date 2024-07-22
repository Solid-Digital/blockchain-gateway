defmodule BlyverWeb.Registration.NewRegistrationLive do
  @moduledoc false
  use BlyverWeb, :registration_live_view

  alias Blyver.Accounts
  alias Blyver.Accounts.User

  @impl true
  def mount(_params, _session, socket) do
    {:ok, assign(socket, changeset: Accounts.change_user(%User{}), trigger_submit: false)}
  end

  @impl true
  def handle_event("save", %{"user" => params}, socket) do
    changeset =
      %User{}
      |> Accounts.change_user(params)
      |> Accounts.validate_new_user()
      |> Map.put(:action, :insert)

    case changeset.valid? do
      true ->
        {:noreply, assign(socket, changeset: changeset, trigger_submit: true)}

      false ->
        {:noreply, assign(socket, changeset: changeset)}
    end
  end
end
