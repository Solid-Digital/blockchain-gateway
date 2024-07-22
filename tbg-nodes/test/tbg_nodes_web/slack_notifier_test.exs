defmodule TbgNodesWeb.SlackNotifierTest do
  use TbgNodesWeb.LibCase

  alias TbgNodesWeb.SlackNotifier

  describe "send_message/1" do
    test "send message success" do
      mock_slack_post_message = fn channel, msg, _ ->
        # Simulate post_message by sending self() a message
        send(self(), %{channel: channel, msg: msg})
        true
      end

      msg = "hello from slack"
      true = SlackNotifier.send_message(msg, %{}, mock_slack_post_message)

      assert_received %{channel: "#alerts-test", msg: "hello from slack"}
    end

    test "send message failure" do
      mock_slack_post_message = fn _, _, _ ->
        %{"ok" => false, error: "i was was supposed to fail"}
      end

      assert false == SlackNotifier.send_message("send will fail", %{}, mock_slack_post_message)
    end
  end
end
