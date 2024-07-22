defmodule TbgNodesWeb.RegistrationController do
  @moduledoc false
  require Logger

  use TbgNodesWeb, :controller
  use Ecto.Schema

  alias TbgNodes.Users.User
  alias TbgNodesWeb.SlackNotifier

  def step_1(conn, _params) do
    changeset = User.changeset_new_login(%User{}, %{})

    render(conn, "step_1.html", changeset: changeset)
  end

  defp next_step?(%{errors: errors}) when length(errors) > 0, do: false
  defp next_step?(%{valid?: true, errors: errors}) when errors == [], do: true

  def submit_step_1(conn, params) do
    # check if email is valid
    changeset = User.changeset_email_validation(%User{}, params)

    case next_step?(changeset) do
      false ->
        changeset = %{changeset | action: :insert}
        render(conn, "step_1.html", changeset: changeset)

      true ->
        updated_changeset = TbgNodes.Users.user_exists(changeset)

        case next_step?(updated_changeset) do
          true ->
            render(conn, "step_2.html", changeset: updated_changeset)

          false ->
            render(conn, "step_1.html", changeset: %{updated_changeset | action: :insert})
        end
    end
  end

  def submit_step_2(conn, %{"user" => user_params}) do
    conn
    |> Pow.Plug.create_user(user_params)
    |> case do
      {:ok, _user, conn} ->
        SlackNotifier.send_message("New user registration! :tada:")

        conn
        |> put_flash(:info, "Welcome!")
        |> redirect(to: Routes.redirect_path(conn, :handle_redirect))

      {:error, changeset, conn} ->
        changeset = %{changeset | action: :insert}
        render(conn, "step_2.html", changeset: changeset)
    end
  end
end
