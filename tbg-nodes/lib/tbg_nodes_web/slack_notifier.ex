defmodule TbgNodesWeb.SlackNotifier do
  @moduledoc false

  require Logger

  def send_message(
        message,
        params \\ %{},
        post_message \\ Application.get_env(:tbg_nodes, :slack_post_message)
      ) do
    channel = Application.get_env(:tbg_nodes, :slack_channel)

    case post_message.(channel, message, params) do
      error_msg = %{"ok" => false} ->
        Logger.warn(
          "Could not post message '#{message}' params '#{params |> inspect}' to slack channel #{
            channel
          } due to '#{error_msg["error"]}'"
        )

        false

      _ ->
        true
    end
  end
end

defmodule TbgNodesWeb.SlackWebChatMock do
  @moduledoc false

  def post_message(_, _, _) do
    %{ok: true}
  end
end
