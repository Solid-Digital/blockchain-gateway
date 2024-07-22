defmodule TbgNodesWeb.NPSFeedbackFlashLive do
  use Phoenix.LiveView
  use Phoenix.HTML

  alias TbgNodes.Users

  @moduledoc """
  Show a flash asking the active user to give an NPS score.
  When the user has given a score, it will be stored and a thank you message
  will be shown.
  """

  @spec mount(any(), map(), Phoenix.LiveView.Socket.t()) :: {:ok, Phoenix.LiveView.Socket.t()}
  def mount(_params, %{"user_id" => user_id} = _session, socket) do
    {:ok,
     socket
     |> assign(:visible, true)
     |> assign(:form_state, "open")
     |> assign(:user_id, user_id)}
  end

  @spec render(map()) :: Phoenix.LiveView.Rendered.t()
  def render(assigns) do
    TbgNodesWeb.LayoutView.render("components/nps_feedback_flash_live.html", assigns)
  end

  @spec handle_event(String.t(), map(), Phoenix.LiveView.Socket.t()) ::
          {:noreply, Phoenix.LiveView.Socket.t()}
  def handle_event("select_nps", %{"nps" => nps_string}, socket) do
    nps = String.to_integer(nps_string)

    # render thank you message
    msg =
      cond do
        nps <= 6 ->
          "Thanks for your feedback. We highly value all ideas and suggestions from our customers, whether they're positive or critical."

        nps >= 7 && nps <= 8 ->
          "Thanks for your feedback. Our goal is to create the best possible product, and your thoughts, ideas, and suggestions play a major role."

        nps >= 9 ->
          "Thanks for your feedback. It's great to hear that you're a fan. Your feedback helps us discover new opportunities to improve our product."

        nps < 1 || nps > 10 ->
          "That's not a score we are expecting."
      end

    # store NPS response
    case Users.add_nps_feedback(socket.assigns.user_id, nps) do
      {:error, _} ->
        {:noreply,
         socket
         |> assign(:form_state, "submitted")
         |> assign(:show_nps_feedback, false)
         |> assign(:msg, "Thanks for your feedback.")}

      {:ok, _} ->
        {:noreply,
         socket
         |> assign(:form_state, "submitted")
         |> assign(:show_nps_feedback, false)
         |> assign(:msg, msg)}
    end
  end

  def handle_event("close_feedback", _value, socket) do
    # in order to register dismissing of feedback a score of 0 is registered
    # dismiss feedback flash
    case Users.add_nps_feedback(socket.assigns.user_id, 0) do
      {:error, _} ->
        {:noreply,
         socket
         |> assign(:visible, false)}

      {:ok, _} ->
        {:noreply,
         socket
         |> assign(:visible, false)}
    end
  end
end
