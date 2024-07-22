defmodule TbgNodes.Repo do
  use Ecto.Repo,
    otp_app: :tbg_nodes,
    adapter: Ecto.Adapters.Postgres
end
