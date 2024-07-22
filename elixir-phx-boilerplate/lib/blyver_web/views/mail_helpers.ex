defmodule BlyverWeb.MailHelpers do
  @moduledoc """
  Convenience methods for emails.
  """

  @doc """
  Retrieves the support email from configuration
  """
  def support_email do
    Application.get_env(:blyver, :emails)[:support_email]
  end

  @doc """
  Returns the support_email formatted as html mailto:.
  """
  def mail_to_support do
    email = support_email()
    "mailto: #{email}"
  end
end
