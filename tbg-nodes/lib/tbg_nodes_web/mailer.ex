defmodule TbgNodesWeb.Mailer do
  @moduledoc false

  use Pow.Phoenix.Mailer
  use Bamboo.Mailer, otp_app: :tbg_nodes

  import Bamboo.Email
  require Logger

  @impl true
  def cast(%{user: user, subject: subject, text: text, html: html, assigns: _assigns}) do
    new_email(
      to: user.email,
      from: "no-reply@unchain.io",
      subject: subject,
      html_body: html,
      text_body: text
    )
  end

  @impl true
  def process(email) do
    deliver_now(email)
  end
end
