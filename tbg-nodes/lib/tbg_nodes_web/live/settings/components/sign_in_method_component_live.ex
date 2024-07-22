defmodule TbgNodesWeb.SignInMethodComponentLive do
  @moduledoc false
  use Phoenix.LiveComponent

  @spec mount(any(), map(), Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, _session, socket) do
    {:ok, socket}
  end

  @spec update(map(), map()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def update(assigns, socket) do
    user = assigns.user

    provider =
      if is_list(user.user_identities) do
        get_provider(assigns.user.user_identities)
      else
        "email"
      end

    {
      :ok,
      socket
      |> assign(:provider, provider)
      |> assign(:user, assigns.user)
    }
  end

  @spec get_provider(map()) :: String.t()
  defp get_provider(user_identities) do
    if Enum.find(user_identities, fn ui -> ui.provider == "github" end) do
      "github"
    else
      "email"
    end
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.SettingsView.render("components/sign-in-method.html", assigns)
  end
end
