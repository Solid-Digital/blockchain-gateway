defmodule BlyverWeb.PowEmailConfirmation.MailerView do
  use BlyverWeb, :mailer_view

  def subject(:email_confirmation, _assigns), do: "Successful Registration!"
end
