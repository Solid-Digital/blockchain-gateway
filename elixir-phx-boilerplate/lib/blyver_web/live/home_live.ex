defmodule BlyverWeb.HomeLive do
  @moduledoc false
  use BlyverWeb, :live_view

  @impl true
  def mount(_params, _session, socket) do
    {:ok, socket}
  end
end
