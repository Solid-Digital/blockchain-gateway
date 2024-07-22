defmodule BlyverWeb.Dashboard.DashboardLive do
  @moduledoc false
  use BlyverWeb, :live_view

  def mount(_params, _session, socket) do
    {:ok, socket}
  end
end
