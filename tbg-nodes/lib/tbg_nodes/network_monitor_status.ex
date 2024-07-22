defmodule TbgNodes.NetworkMonitor.Status do
  @moduledoc false

  @type t :: :down | :live | :ready | :fetching_status | :not_available

  defmacro down, do: quote(do: :down)
  defmacro live, do: quote(do: :live)
  defmacro ready, do: quote(do: :ready)
  defmacro fetching_status, do: quote(do: :fetching_status)
  defmacro not_available, do: quote(do: :not_available)

  @spec to_string(t()) :: String.t()
  def to_string(status) do
    case status do
      :down -> "Down"
      :live -> "Live"
      :ready -> "Ready"
      :fetching_status -> "Fetching status"
      :not_available -> "Not available"
    end
  end
end
