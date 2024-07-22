import Config

besu_image = System.get_env("BESU_IMAGE") || "hyperledger/besu:1.5.4"

config :tbg_nodes, TbgNodes.Networks, besu_image: besu_image

test_user_email = System.get_env("LIVE_TESTER_USER_EMAIL") || "bot@unchain.io"

test_user_password =
  System.get_env("LIVE_TESTER_USER_PASSWORD") ||
    :crypto.strong_rand_bytes(12)
    |> Base.encode32()
    |> binary_part(0, 12)
    |> String.downcase()

config :tbg_nodes, TbgNodes.LiveTester,
  loop_interval: 3600 * 1000,
  test_user_email: test_user_email,
  test_user_password: test_user_password,
  enabled: false

config :tbg_nodes, TbgNodes.Networks,
  network_url_templates: %{
    "ropsten" => %{
      "http" => "https://<%= network_name %>.infura.io/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "websocket" => "https://<%= network_name %>.infura.io/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "liveness" => "",
      "readiness" => ""
    },
    "mainnet" => %{
      "http" => "https://<%= network_name %>.infura.io/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "websocket" => "wss://<%= network_name %>.infura.io/ws/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "liveness" => "",
      "readiness" => ""
    },
    "ropsten-archive" => %{
      "http" => "https://ropsten.infura.io/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "websocket" => "wss://ropsten.infura.io/ws/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "liveness" => "",
      "readiness" => ""
    },
    "mainnet-archive" => %{
      "http" => "https://mainnet.infura.io/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "websocket" => "wss://mainnet.infura.io/ws/v3/bf6f955a7c314b15ae2e6643f7c1c5c6",
      "liveness" => "",
      "readiness" => ""
    },
    "permissioned" => %{
      "http" => "http://permissioned.nodes.localhost/v0/<%= network_uuid %>",
      "websocket" => "ws://permissioned.nodes.localhost/v0/ws/<%= network_uuid %>",
      "liveness" =>
        "http://localhost:8001/api/v1/namespaces/<%= namespace_name %>/services/<%= service_name %>:http/proxy/liveness",
      "readiness" =>
        "http://localhost:8001/api/v1/namespaces/<%= namespace_name %>/services/<%= service_name %>:http/proxy/readiness?maxBlocksBehind=5"
    }
  }

config :tbg_nodes, TbgNodes.NetworkMonitor,
  loop_interval: 1 * 1000,
  enabled: false

postgres_host = System.get_env("POSTGRES_HOST") || "localhost"

config :tbg_nodes, :telemetry, false

# Configure your database
config :tbg_nodes, TbgNodes.Repo,
  username: "postgres",
  password: "postgres",
  database: "tbg_nodes_test",
  hostname: postgres_host,
  pool: Ecto.Adapters.SQL.Sandbox,
  timeout: 300_000

redis_url = System.get_env("REDIS_URL") || "redis://localhost:6379"

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :tbg_nodes, TbgNodesWeb.Endpoint,
  http: [port: 4002],
  server: true,
  redis_url: redis_url

# Print only warnings and errors during test
config :logger,
       :console,
       level: :warn

config :tbg_nodes, TbgNodesWeb.Mailer, adapter: Bamboo.TestAdapter

tbg_nodes_access_kubeconfig =
  System.get_env("TBG_NODES_ACCESS_KUBECONFIG") || "kubeconfig.local.yaml"

config :k8s,
  clusters: %{
    default: %{
      conn: tbg_nodes_access_kubeconfig,
      conn_opts: [cluster: "default", user: "default"]
    }
  }

infra_api =
  case System.get_env("INFRA_API") || "mock" do
    "mock" -> TbgNodes.PermissionedEthereumNetworks.InfraAPIMock
    "k8s" -> TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s
  end

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks, infra_api: infra_api

# This storage class is specific to k3s: https://rancher.com/docs/k3s/latest/en/storage/
config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s,
  statefulset_storage_class: "local-path"

config :pow, Pow.Ecto.Schema.Password, iterations: 1

config :tbg_nodes, :pow_assent,
  providers: [
    github: [
      client_id: "8c76e5bc0d925d53b905",
      client_secret: "3f7d5dddebf067f03fa87257b379a7b6fa032b3a",
      strategy: Assent.Strategy.Github
    ]
  ]

config :tbg_nodes, :env, :test

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s, deployment_target: "local"

config :slack, api_token: "xoxb-13655146816-1110260169584-wXuiBdzLOZ3lzVMiDwkbKv28"
config :tbg_nodes, :slack_post_message, &TbgNodesWeb.SlackWebChatMock.post_message/3

config :tbg_nodes, :slack_channel, "#alerts-test"

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s,
  ingress_host: "permissioned.nodes.test.localhost"

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s,
  ingress_basicauth_host:
    "http://tbg-nodes-auth-test-local.tbg-nodes-local.svc.cluster.local:8080/permissioned"

config :tbg_nodes, :rancher,
  cluster_id: "",
  project_id: ""
