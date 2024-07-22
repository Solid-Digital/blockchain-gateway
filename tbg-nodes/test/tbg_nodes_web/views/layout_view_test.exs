defmodule TbgNodesWeb.LayoutViewTest do
  use TbgNodesWeb.ConnCase, async: true

  alias TbgNodesWeb.LayoutView

  test "active_class returns active when conn and path match" do
    path_info = ["resource", "resource-id-12345"]
    link_path = "/resource"

    assert LayoutView.active_class(path_info, link_path) == "active"
  end

  test "active_class returns nil when conn and path dont match" do
    path_info = ["resource", "resource-id-12345"]
    link_path = "/other-resource"

    assert is_nil(LayoutView.active_class(path_info, link_path))
  end
end
