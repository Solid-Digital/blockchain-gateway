defmodule BlyverWeb.Pow.Mailer do
  @moduledoc false

  use Pow.Phoenix.Mailer
  use Bamboo.Mailer, otp_app: :blyver

  import Bamboo.Email

  @impl true
  def cast(%{user: user, subject: subject, text: text, html: html, assigns: _assigns}) do
    email_config = Application.get_env(:blyver, :emails)

    new_email(
      to: user.email,
      from: email_config[:from],
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
