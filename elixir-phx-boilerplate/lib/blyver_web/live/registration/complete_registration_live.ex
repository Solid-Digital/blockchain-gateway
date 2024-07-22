defmodule BlyverWeb.Registration.CompleteRegistrationLive do
  @moduledoc false
  use BlyverWeb, :registration_live_view

  def mount(_params, _session, socket) do
    {:ok, socket}
  end
end
