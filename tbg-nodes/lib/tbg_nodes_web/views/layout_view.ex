defmodule TbgNodesWeb.LayoutView do
  @moduledoc false

  use TbgNodesWeb, :view

  def version do
    Application.spec(:tbg_nodes, :vsn) |> to_string()
  end

  # This function is used to determine whether a path is active
  # If the @conn.path_info matches with the path of the link, i.e. Routes.network_path(@conn, :index)
  # returns 'active', otherwise it returns nil
  def active_class(path_info, link_path) do
    current_path = "/" <> List.first(path_info)

    if link_path == current_path do
      "active"
    else
      nil
    end
  end
end
