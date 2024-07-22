defmodule Blyver.Utils.S3Test do
  use Blyver.DataCase

  @default_bucket "test-bucket"
  @test_file_path "./test/blyver/utils/object_store/object.txt"

  alias Blyver.Utils.S3

  describe "list_objects/1" do
    setup [:create_or_clear_bucket, :put_test_file]

    test "returns list of objects successfully", %{} do
      {:ok, body} = S3.list_objects(@default_bucket)
      assert length(body.contents) == 1
    end

    test "returns error when bucket does not exists", %{} do
      {:error, result} = S3.list_objects("random-bucket")
      assert result == "could not find bucket or directory"
    end
  end

  describe "put_object/2" do
    setup [:create_or_clear_bucket]

    test "puts object in bucket successfully", %{} do
      {:ok, result} = S3.put_object(@default_bucket, "object.txt", object_fixture())
      assert String.length(result) >= 34
    end
  end

  describe "get_object/2" do
    setup [:create_or_clear_bucket, :put_test_file]

    test "gets object from bucket successfully", %{} do
      {:ok, result} = S3.get_object(@default_bucket, "object.txt")
      assert result == "test-object\n"
    end
  end

  def object_fixture do
    File.read!(@test_file_path)
  end

  def create_or_clear_bucket(%{}) do
    case ExAws.S3.list_buckets() |> ExAws.request!() do
      %{body: %{buckets: [], owner: _owner}, headers: _headers, status_code: 200} ->
        put_default_bucket()

      %{body: %{buckets: _buckets, owner: _owner}, headers: _headers, status_code: 200} ->
        clear_bucket()
    end
  end

  def put_default_bucket do
    ExAws.S3.put_bucket(@default_bucket, "local") |> ExAws.request!()
    :ok
  end

  def clear_bucket do
    stream =
      ExAws.S3.list_objects(@default_bucket, prefix: "/")
      |> ExAws.stream!()
      |> Stream.map(& &1.key)

    ExAws.S3.delete_all_objects(@default_bucket, stream) |> ExAws.request()
    ExAws.S3.delete_bucket(@default_bucket) |> ExAws.request!()
    put_default_bucket()
  end

  def put_test_file(%{}) do
    ExAws.S3.put_object(@default_bucket, "object.txt", object_fixture()) |> ExAws.request!()
    :ok
  end
end
