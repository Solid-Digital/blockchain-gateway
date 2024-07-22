defmodule TbgNodesWeb.RegistrationControllerTest do
  use TbgNodesWeb.ConnCase

  alias TbgNodes.Users

  @moduletag :RegistrationController
  describe "RegistrationController" do
    test "step_1 loads", %{conn: conn} do
      conn = get(conn, Routes.registration_path(conn, :step_1))
      assert html_response(conn, 200) =~ "Signup with your email or through your Github account."
    end

    test "submit_step_1: warning shown with empty email address", %{conn: conn} do
      test_params = %{
        email: ""
      }

      conn = post(conn, Routes.registration_path(conn, :submit_step_1, test_params))
      assert html_response(conn, 200) =~ "can&#39;t be blank"
    end

    test "submit_step_1: proceed to step_2 with valid email", %{conn: conn} do
      test_params = %{
        email: "user@domain.com"
      }

      conn = post(conn, Routes.registration_path(conn, :submit_step_1, test_params))
      assert html_response(conn, 200) =~ "Secure your account with a well crafted password"
    end

    test "submit_step_2: redirects to / with valid email & password", %{conn: conn} do
      test_params = %{
        user: %{
          email: "user@domain.com",
          password: "1234567890!A"
        }
      }

      conn = post(conn, Routes.registration_path(conn, :submit_step_2, test_params))
      assert redirected_to(conn) == Routes.redirect_path(conn, :handle_redirect)
    end

    test "submit_step_2: invalid password doesn't succeed", %{conn: conn} do
      test_params = %{
        user: %{
          email: "user@domain.com",
          password: "ABC"
        }
      }

      conn = post(conn, Routes.registration_path(conn, :submit_step_2, test_params))
      assert html_response(conn, 200) =~ "should be at least 12 character(s)"
    end

    test "submit_step_1: 'user already exists' when second user with same email tries to signup",
         %{conn: conn} do
      email = "user@domain.com"

      conn1_params = %{
        user: %{
          email: email,
          password: "1234567890!A"
        }
      }

      conn2_params = %{
        email: email
      }

      _conn1 = post(conn, Routes.registration_path(conn, :submit_step_2, conn1_params))
      conn2 = post(conn, Routes.registration_path(conn, :submit_step_1, conn2_params))
      assert html_response(conn2, 200) =~ "user already exists"
    end
  end

  describe "admin registration" do
    test "non-admin users register as regular users", %{conn: conn} do
      test_params = %{
        user: %{
          email: "regular@user.io",
          password: "1234567890!A"
        }
      }

      _conn = post(conn, Routes.registration_path(conn, :submit_step_2, test_params))

      user = Users.get_user_by_email!(test_params[:user][:email])

      assert user.role == "user"
    end
  end
end
