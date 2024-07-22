defmodule TbgNodesWeb.SignInMethodComponentLiveTest do
  @moduledoc false
  use TbgNodesWeb.ConnCase
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  describe "github user" do
    setup [:create_github_user]

    test "shows github user type correctly", %{user: user} do
      assert render_component(TbgNodesWeb.SignInMethodComponentLive, user: user) =~ "Github"
    end
  end

  describe "email user" do
    setup [:create_user]

    test "shows email user type correctly", %{user: user} do
      html = render_component(TbgNodesWeb.SignInMethodComponentLive, user: user)
      assert html =~ "Email and password"
      assert html =~ "Email address: " <> user.email
    end
  end
end
