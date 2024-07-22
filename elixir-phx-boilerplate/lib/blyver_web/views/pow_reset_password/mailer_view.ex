defmodule BlyverWeb.PowResetPassword.MailerView do
  use BlyverWeb, :mailer_view

  def subject(:reset_password, _assigns), do: "Recover your password"
end
