defmodule TbgNodes.UsersTest do
  use TbgNodes.DataCase

  alias TbgNodes.Users

  describe "list users with resources" do
    setup [
      :count_users,
      :create_user,
      :create_permissioned_ethereum_network,
      :create_permissioned_node,
      :create_public_ethereum_network_with_interfaces
    ]

    test "returns user with preloaded network and nodes", %{
      user_count: user_count,
      user: new_user
    } do
      users = Users.list_users_with_resources()
      assert length(users) == user_count + 1

      user = users |> Enum.find(fn user -> user.email == new_user.email end)

      assert length(user.permissioned_ethereum_networks) == 1
      assert length(user.public_ethereum_networks) == 1

      nodes = Map.get(List.first(user.permissioned_ethereum_networks), :besu_nodes)
      assert length(nodes) == 1
    end
  end

  describe "show_nps_feedback" do
    test "not if user created less than a day ago" do
      user = user_fixture()

      assert Users.show_nps_feedback(user) == false
    end

    test "if user is created more than a day ago" do
      user = create_old_user()

      assert Users.show_nps_feedback(user) == true
    end

    test "not if feedback has been dismissed in past 20 days" do
      user = create_old_user()

      Users.add_nps_feedback(user.id, 0)

      assert Users.show_nps_feedback(user) == false
    end

    test "if been dismissed more than 20 days ago" do
      user = create_old_user()
      add_old_feedback(user.id, 0, 21 * 24 * 3600)

      assert Users.show_nps_feedback(user) == true
    end

    test "not if feedback is recorded within past 100 days" do
      user = create_old_user()

      Users.add_nps_feedback(user.id, 7)

      assert Users.show_nps_feedback(user) == false
    end

    test "if feedback is recorded more than 100 days ago" do
      user = create_old_user()

      add_old_feedback(user.id, 7, 101 * 24 * 3600)

      assert Users.show_nps_feedback(user) == false
    end
  end

  defp create_old_user do
    now = DateTime.truncate(DateTime.utc_now(), :second)
    inserted_at = DateTime.add(now, -2 * 24 * 3600, :second)
    user = %Users.User{:email => "test@test.com", :inserted_at => inserted_at}
    {:ok, user} = Repo.insert(user)

    user
  end

  defp add_old_feedback(user_id, score, seconds_ago) do
    now = DateTime.truncate(DateTime.utc_now(), :second)
    days_ago = DateTime.add(now, -seconds_ago, :second)

    nps_feedback = %Users.NpsFeedback{
      :user_id => user_id,
      :score => score,
      :inserted_at => days_ago
    }

    Repo.insert(nps_feedback)
  end
end
