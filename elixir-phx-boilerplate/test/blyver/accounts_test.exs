defmodule Blyver.AccountsTest do
  use Blyver.DataCase
  alias Blyver.Accounts
  alias Blyver.Accounts.User

  describe "user" do
    test "change_user/1 returns an Ecto changeset" do
      user = %User{}
      attrs = %{email: "test@user.com"}
      assert %Ecto.Changeset{} = Accounts.change_user(user, attrs)
    end

    test "validate_new_user/1 returns a changeset with errors if the email exists" do
      user = user_fixture()
      attrs = %{valid_user_attrs() | email: user.email}
      changeset = User.changeset(%User{}, attrs)
      changeset = Accounts.validate_new_user(changeset)

      refute changeset.valid?
      assert [email: {"User already exists", []}] = changeset.errors
    end

    test "validate_new_user/1 returns a changeset without errors if the email doesn't exist" do
      attrs = valid_user_attrs()
      changeset = User.changeset(%User{}, attrs)

      assert Accounts.validate_new_user(changeset).valid?
    end
  end
end
