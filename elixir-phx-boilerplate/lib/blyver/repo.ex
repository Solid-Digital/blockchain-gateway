defmodule Blyver.Repo do
  use Ecto.Repo,
    otp_app: :blyver,
    adapter: Ecto.Adapters.Postgres
end
