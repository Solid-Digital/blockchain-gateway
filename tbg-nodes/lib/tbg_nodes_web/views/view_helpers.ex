defmodule TbgNodesWeb.ViewHelpers do
  @moduledoc "This module contains utility functions that are used across different views, for example a 'since' function
  to format date time to a 'since .. days' string"

  def since(inserted_at) do
    unix_now = DateTime.to_unix(DateTime.utc_now())
    unix_since = DateTime.to_unix(DateTime.from_naive!(inserted_at, "Etc/UTC"))
    elapsed = unix_now - unix_since

    cond do
      elapsed < 60 ->
        "1 minute"

      elapsed < 60 * 60 ->
        Integer.to_string(trunc(elapsed / 60)) <> " minutes"

      elapsed < 60 * 60 * 24 ->
        Integer.to_string(trunc(elapsed / 60 / 60)) <> " hours"

      elapsed > 60 * 60 * 24 ->
        Integer.to_string(trunc(elapsed / 60 / 60 / 24)) <> " days"
    end
  end

  def format_date(date) do
    date
    |> Timex.format!("%B %-d, %Y", :strftime)
  end
end
