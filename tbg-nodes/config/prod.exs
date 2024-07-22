# In this file, we load production configuration and secrets
# from environment variables. You can also hardcode secrets,
# although such is generally not recommended and you have to
# remember to add this file to your .gitignore.
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
  enabled: true

host =
  System.get_env("HOST") ||
    raise """
    environment variable HOST is missing.
    """

deployment_target =
  System.get_env("DEPLOYMENT_TARGET") ||
    raise """
    Deployment target cannot be empty.
    """

network_url_modifier =
  case deployment_target do
    "prod" ->
      "."

    "staging" ->
      "." <> "staging" <> "."
  end

config :logger, level: :info

namespace = System.get_env("NAMESPACE")

config :libcluster,
  topologies: [
    k8s: [
      strategy: Elixir.Cluster.Strategy.Kubernetes,
      config: [
        mode: :ip,
        kubernetes_node_basename: "tbg_nodes",
        kubernetes_selector: "app.kubernetes.io/name=tbg-nodes",
        kubernetes_namespace: namespace,
        polling_interval: 10_000,
        kubernetes_ip_lookup_mode: :pods
      ]
    ]
  ]

config :tbg_nodes, TbgNodes.Networks,
  network_url_templates: %{
    "ropsten" => %{
      "http" =>
        "https://ropsten-archive#{network_url_modifier}nodes.unchain.io/v0/<%= network_uuid %>",
      "websocket" =>
        "wss://ropsten-archive#{network_url_modifier}nodes.unchain.io/v0/ws/<%= network_uuid %>",
      "liveness" => "ropsten-archive.ropsten-nodes-fullsync.svc.cluster.local:8545/liveness",
      "readiness" =>
        "ropsten-archive.ropsten-nodes-fullsync.svc.cluster.local:8545/readiness?maxBlocksBehind=100"
    },
    "ropsten-archive" => %{
      "http" =>
        "https://ropsten-archive#{network_url_modifier}nodes.unchain.io/v0/<%= network_uuid %>",
      "websocket" =>
        "wss://ropsten-archive#{network_url_modifier}nodes.unchain.io/v0/ws/<%= network_uuid %>",
      "liveness" => "ropsten-archive.ropsten-nodes-fullsync.svc.cluster.local:8545/liveness",
      "readiness" =>
        "ropsten-archive.ropsten-nodes-fullsync.svc.cluster.local:8545/readiness?maxBlocksBehind=100"
    },
    "mainnet" => %{
      "http" =>
        "https://mainnet#{network_url_modifier}nodes-hz.unchain.io/v0/<%= network_uuid %>",
      "websocket" =>
        "wss://mainnet#{network_url_modifier}nodes-hz.unchain.io/v0/ws/<%= network_uuid %>",
      "liveness" => "https://mainnet#{network_url_modifier}nodes-hz.unchain.io/v0/liveness",
      "readiness" =>
        "https://mainnet#{network_url_modifier}nodes-hz.unchain.io/v0/readiness?maxBlocksBehind=5"
    },
    "mainnet-archive" => %{
      "http" =>
        "https://mainnet-archive#{network_url_modifier}nodes-hz.unchain.io/v0/<%= network_uuid %>",
      "websocket" =>
        "wss://mainnet-archive#{network_url_modifier}nodes-hz.unchain.io/v0/ws/<%= network_uuid %>",
      "liveness" =>
        "https://mainnet-archive#{network_url_modifier}nodes-hz.unchain.io/v0/liveness",
      "readiness" =>
        "https://mainnet-archive#{network_url_modifier}nodes-hz.unchain.io/v0/readiness?maxBlocksBehind=5"
    },
    "permissioned" => %{
      "http" =>
        "https://permissioned#{network_url_modifier}nodes.unchain.io/v0/<%= network_uuid %>",
      "websocket" =>
        "wss://permissioned#{network_url_modifier}nodes.unchain.io/v0/ws/<%= network_uuid %>",
      "liveness" => "<%= service_name %>.<%= namespace_name %>.svc.cluster.local:8545/liveness",
      "readiness" =>
        "<%= service_name %>.<%= namespace_name %>.svc.cluster.local:8545/readiness?maxBlocksBehind=5"
      # external-interface-9b58a79d-320e-4d8a-98c6-bd0242faa357.network-f588fba4-f3bc-4d08-a216-7a039852729b.svc.cluster.local:8545/readiness
    }
  }

database_url =
  System.get_env("DATABASE_URL") ||
    raise """
    environment variable DATABASE_URL is missing.
    For example: ecto://USER:PASS@HOST/DATABASE
    """

config :tbg_nodes,
       TbgNodes.Repo,
       ssl: true,
       url: database_url,
       prepare: :unnamed,
       pool_size: String.to_integer(System.get_env("POOL_SIZE") || "10")

secret_key_base =
  System.get_env("SECRET_KEY_BASE") ||
    raise """
    environment variable SECRET_KEY_BASE is missing.
    You can generate one by calling: mix phx.gen.secret
    """

redis_url =
  System.get_env("REDIS_URL") ||
    raise """
    environment variable REDIS_URL is missing.
    """

config :tbg_nodes,
       TbgNodesWeb.Endpoint,
       url: [
         scheme: "https",
         host: host,
         port: 443
       ],
       cache_static_manifest: "priv/static/cache_manifest.json",
       http: [
         :inet6,
         port: String.to_integer(System.get_env("PORT") || "4000")
       ],
       secret_key_base: secret_key_base,
       pubsub_server: :pubsub,
       redis_url: redis_url,
       redis_ssl: true,
       redis_socket_opts_provider: TbgNodes.Redis.ConfigDigitalOcean,
       check_origin: ["https://" <> host]

github_client_id =
  System.get_env("GITHUB_CLIENT_ID") ||
    raise """
    environment variable GITHUB_CLIENT_ID is missing.
    """

github_client_secret =
  System.get_env("GITHUB_CLIENT_SECRET") ||
    raise """
    environment variable GITHUB_CLIENT_SECRET is missing.
    """

config :tbg_nodes,
       :pow_assent,
       providers: [
         github: [
           client_id: github_client_id,
           client_secret: github_client_secret,
           strategy: Assent.Strategy.Github
         ]
       ]

config :slack, api_token: "xoxb-13655146816-1110260169584-wXuiBdzLOZ3lzVMiDwkbKv28"

mailgun_api_key =
  System.get_env("MAILGUN_API_KEY") ||
    raise """
    environment variable MAILGUN_API_KEY is missing.
    """

config :tbg_nodes,
       TbgNodesWeb.Mailer,
       adapter: Bamboo.MailgunAdapter,
       api_key: mailgun_api_key,
       domain: "mg.unchain.io"

# ## Using releases (Elixir v1.9+)
#
# If you are doing OTP releases, you need to instruct Phoenix
# to start each relevant endpoint:
#
config :tbg_nodes, TbgNodesWeb.Endpoint, server: true

#
# Then you can assemble a release by calling `mix release`.
# See `mix help release` for more information.

alerts_slack_channel = System.get_env("ALERTS_SLACK_CHANNEL") || "#alerts-staging"

config :tbg_nodes, :slack_post_message, &Slack.Web.Chat.post_message/3
config :tbg_nodes, :slack_channel, alerts_slack_channel

tbg_nodes_access_kubeconfig =
  System.get_env("TBG_NODES_ACCESS_KUBECONFIG") || "kubeconfig.local.yaml"

config :k8s,
  clusters: %{
    default: %{
      conn: tbg_nodes_access_kubeconfig,
      conn_opts: [cluster: "default", user: "default"]
    }
  }

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks,
  infra_api: TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s,
  statefulset_storage_class: "do-block-storage"

ingress_host =
  case deployment_target do
    "prod" -> "permissioned.nodes.unchain.io"
    "staging" -> "permissioned.staging.nodes.unchain.io"
  end

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s,
  deployment_target: deployment_target

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s, ingress_host: ingress_host

ingress_basicauth_host =
  case deployment_target do
    "prod" ->
      "https://tbg-nodes-auth.prod.dgo.unchain.io/permissioned"

    "staging" ->
      "https://tbg-nodes-auth.staging.dgo.unchain.io/permissioned"
  end

config :tbg_nodes, TbgNodes.PermissionedEthereumNetworks.InfraAPIK8s,
  ingress_basicauth_host: ingress_basicauth_host

rancher_cluster_id =
  System.get_env("RANCHER_CLUSTER_ID") ||
    raise """
    environment variable RANCHER_CLUSTER_ID is missing.
    """

rancher_project_id =
  System.get_env("RANCHER_PROJECT_ID") ||
    raise """
    environment variable RANCHER_PROJECT_ID is missing.
    """

config :tbg_nodes, :rancher,
  cluster_id: rancher_cluster_id,
  project_id: rancher_project_id
