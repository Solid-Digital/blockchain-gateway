defmodule TbgNodesWeb.ShowNpsFeedbackPlug do
  @moduledoc false

  import Plug.Conn

  def init(config), do: config

  def call(conn, _opts) do
    # If user is not logged in, don't bother checking
    if Kernel.is_nil(conn.assigns.current_user) do
      conn
    else
      assign(conn, :show_nps_feedback, show_nps_form(conn))
    end
  end

  defp show_nps_form(conn) do
    # If show_nps_feedback does not exist yet, run check
    if Enum.member?(conn.assigns, :show_nps_feedback) == false do
      TbgNodes.Users.show_nps_feedback(conn.assigns.current_user)
      # Return false if check has already run
    else
      false
    end
  end
end
