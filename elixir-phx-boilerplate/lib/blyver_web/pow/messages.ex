defmodule BlyverWeb.Pow.Messages do
  @moduledoc false
  use Pow.Phoenix.Messages

  use Pow.Extension.Phoenix.Messages,
    extensions: [PowEmailConfirmation, PowResetPassword]

  import BlyverWeb.Gettext

  def invalid_credentials(_conn), do: gettext("Invalid email and password combination.")

  def pow_reset_password_email_has_been_sent(_conn),
    do:
      gettext(
        "We have sent you an email with a recovery link. Follow the link to recover your password."
      )

  def pow_reset_password_maybe_email_has_been_sent(_conn),
    do:
      gettext(
        "We have sent you an email with a recovery link. Follow the link to recover your password."
      )
end
