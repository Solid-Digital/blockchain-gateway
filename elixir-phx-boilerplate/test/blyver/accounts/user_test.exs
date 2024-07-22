defmodule Blyver.Accounts.UserTest do
  use Blyver.DataCase

  alias Blyver.Accounts.User

  test "changeset/2 pow changeset validation is setup correctly" do
    invalid_email_attrs = %{email: "testexample.com"}
    no_confirm_password_attrs = %{email: "test@example.com", password: "Thepassword$"}
    no_password_attrs = %{email: "test@example.com"}

    attribs = [invalid_email_attrs, no_confirm_password_attrs, no_password_attrs]
    refute_all(attribs)
  end

  test "changeset/2 validates password for presence of upper case, lower case and special character" do
    valid_attrs = %{
      email: "test@example.com",
      password: "Thepassword$",
      confirm_password: "Thepassword$"
    }

    no_upper_case_attrs = %{
      email: "test@example.com",
      password: "thepassword$",
      confirm_password: "thepassword$"
    }

    no_lower_case_attrs = %{
      email: "test@example.com",
      password: "THEPASSWORD$",
      confirm_password: "THEPASSWORD$"
    }

    no_special_xter_attrs = %{
      email: "test@example.com",
      password: "Thepassword",
      confirm_password: "Thepassword"
    }

    attribs = [no_upper_case_attrs, no_lower_case_attrs, no_special_xter_attrs]
    refute_all(attribs)
    assert User.changeset(%User{}, valid_attrs).valid?
  end

  test "login_changeset/2 validates the presense of both email and password in the changeset" do
    valid_attrs = %{
      email: "test@example.com",
      password: "thePassword"
    }

    invalid_attrs = %{
      email: "",
      password: ""
    }

    assert User.login_changeset(%User{}, valid_attrs).valid?
    refute User.login_changeset(%User{}, invalid_attrs).valid?
  end

  test "new_reset_password_changeset/2 validates the presense of email in the changeset" do
    valid_attrs = %{email: "test@example.com"}

    invalid_attrs = %{email: ""}

    assert User.new_reset_password_changeset(%User{}, valid_attrs).valid?
    refute User.new_reset_password_changeset(%User{}, invalid_attrs).valid?
  end

  test "submit_reset_password_changeset/2 performs password validation" do
    valid_attrs = %{
      password: "Thepassword$",
      password_confirmation: "Thepassword$"
    }

    invalid_attrs = %{
      password: "thepassword$",
      password_confirmation: "thepassword"
    }

    assert User.submit_reset_password_changeset(%User{}, valid_attrs).valid?
    refute User.submit_reset_password_changeset(%User{}, invalid_attrs).valid?
  end

  defp refute_all(attribs) do
    for attrs <- attribs do
      refute User.changeset(%User{}, attrs).valid?
    end
  end
end
