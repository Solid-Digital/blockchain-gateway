defmodule TbgNodes.Redis.ConfigDigitalOcean do
  @moduledoc """
  Custom configs for connecting to redis on DigitalOcean
  """

  def socket_opts do
    [
      customize_hostname_check: [
        match_fun: :public_key.pkix_verify_hostname_match_fun(:https)
      ]
    ]
  end
end
