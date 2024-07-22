defmodule TbgNodesWeb.ShowNpsFeedbackPlugTests do
  use TbgNodesWeb.ConnCase

  import Plug.Conn

  alias TbgNodesWeb.ShowNpsFeedbackPlug

  describe "authed user call/2" do
    setup [:create_user, :auth_user]

    test "result is saved in session", %{conn: conn} do
      conn = ShowNpsFeedbackPlug.call(conn, %{})
      assert Map.has_key?(conn.assigns, :show_nps_feedback) == true
    end
  end

  describe "not authed user call/2" do
    test "returns false for not authed user", %{conn: conn} do
      conn = assign(conn, :current_user, nil)
      conn = ShowNpsFeedbackPlug.call(conn, %{})
      assert Map.has_key?(conn.assigns, :show_nps_feedback) == false
    end
  end
end
