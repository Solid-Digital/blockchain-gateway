defmodule TbgNodesWeb.NPSFeedbackFlashLiveTest do
  use TbgNodesWeb.ConnCase
  import Phoenix.LiveViewTest
  @endpoint TbgNodesWeb.Endpoint

  alias TbgNodes.Users
  alias TbgNodesWeb.NPSFeedbackFlashLive

  describe "mount/3" do
    setup [:create_user]

    test "feedback flash is visible", %{conn: conn, user: user} do
      {:ok, _view, html} =
        live_isolated(conn, NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      assert html =~ "<div id=\"nps-feedback-flash\">"
    end

    test "form is visible", %{conn: conn, user: user} do
      {:ok, _view, html} =
        live_isolated(conn, NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      assert html =~ "<div id=\"nps-form\">"
    end
  end

  describe "handle_event select_nps" do
    setup [:create_user]

    test "detractor scores show thank you message welcome criticism", %{conn: conn, user: user} do
      {:ok, view, _html} =
        live_isolated(conn, NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      assert render_click(view, :select_nps, %{nps: "1"}) =~
               "Thanks for your feedback. We highly value all ideas and suggestions from our customers, whether they&apos;re positive or critical."
    end

    test "neutral scores show thank you message appreciation", %{conn: conn, user: user} do
      {:ok, view, _html} =
        live_isolated(conn, NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      assert render_click(view, :select_nps, %{nps: "7"}) =~
               "Thanks for your feedback. Our goal is to create the best possible product, and your thoughts, ideas, and suggestions play a major role."
    end

    test "promotor scores show thank you message fan", %{conn: conn, user: user} do
      {:ok, view, _html} =
        live_isolated(conn, NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      assert render_click(view, :select_nps, %{nps: "9"}) =~
               "Thanks for your feedback. It&apos;s great to hear that you&apos;re a fan. Your feedback helps us discover new opportunities to improve our product."
    end

    test "valid nps scores get stored in db", %{conn: conn, user: user} do
      {:ok, view, _html} =
        live_isolated(conn, NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      render_click(view, :select_nps, %{nps: "7"})

      feedback = Users.get_nps_feedback_by_user_id(user.id)
      assert feedback.score == 7
    end

    test "invalid nps scores do not get processed", %{conn: conn, user: user} do
      {:ok, view, _html} =
        live_isolated(conn, TbgNodesWeb.NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      res = render_click(view, :select_nps, %{nps: "1000"})

      assert res =~ "<div class=\"feedback-text\">Thanks for your feedback.</div>"

      assert_raise Ecto.NoResultsError, fn ->
        Users.get_nps_feedback_by_user_id(user.id)
      end
    end
  end

  describe "handle close_feedback" do
    setup [:create_user]

    test "close makes feedback container invisible", %{conn: conn, user: user} do
      {:ok, view, _html} =
        live_isolated(conn, TbgNodesWeb.NPSFeedbackFlashLive, session: %{"user_id" => user.id})

      assert render_click(view, :close_feedback) == ""
      feedback = Users.get_nps_feedback_by_user_id(user.id)
      assert feedback.score == 0
    end
  end
end
