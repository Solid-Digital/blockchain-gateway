defmodule TbgNodesWeb.RouterTest do
  use TbgNodesWeb.ConnCase

  describe "browser pipeline" do
    setup [:create_user, :auth_user]

    test "uses correct content-security-policy", %{conn: conn} do
      conn = get(conn, "/")
      nonce = conn.assigns[:csp_nonce]

      csp_tuple =
        {"content-security-policy",
         "default-src 'self' https://*.unchain.io; style-src 'self' 'unsafe-inline' *.fontawesome.com; script-src 'self' cdnjs.cloudflare.com 'nonce-#{
           nonce
         }' 'unsafe-eval'; font-src 'self' data: 'self' *.fontawesome.com; img-src 'self' data:;"}

      assert csp_tuple in conn.resp_headers
    end
  end
end
