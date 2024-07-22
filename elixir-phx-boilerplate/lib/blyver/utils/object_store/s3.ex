defmodule Blyver.Utils.S3 do
  @moduledoc """
  This module provides a S3 implementation of the Objects behaviour and can be used to store and retrieve objects from
  an S3-compatible object store. The implementation is based on ExAws and expects configuration fields. Please refer
  to the S3 documentation to getting started.
  """
  @behaviour Blyver.Utils.Objects

  @impl Blyver.Utils.Objects
  def list_objects(bucket, directory \\ "/") do
    case ExAws.S3.list_objects(bucket, prefix: directory)
         |> ExAws.request() do
      {:ok, %{body: body, headers: _headers, status_code: 200}} ->
        {:ok, body}

      {:error, {:http_error, 404, _result}} ->
        {:error, "could not find bucket or directory"}

      {:error, {_error, _status_code, result}} ->
        {:error, result.body}
    end
  end

  @impl Blyver.Utils.Objects
  def get_object(bucket, filepath) do
    case ExAws.S3.get_object(bucket, filepath) |> ExAws.request() do
      {:ok, %{body: body, headers: _headers, status_code: 200} = _result} ->
        {:ok, body}

      {:error, {_error, _status_code, result}} ->
        {:error, result.body}
    end
  end

  @impl Blyver.Utils.Objects
  def put_object(bucket, filepath, object) do
    case ExAws.S3.put_object(bucket, filepath, object)
         |> ExAws.request() do
      {:ok, %{body: _body, headers: resp_headers, status_code: 200} = _result} ->
        # The ETag is created by S3 and can be a md5 checksum of the whole file or the md5 hash
        # of a concatenation of all file chunks.
        {"ETag", etag} = List.keyfind(resp_headers, "ETag", 0)
        {:ok, etag}

      {:ok, %{body: _body, headers: _headers, status_code: _status_code} = result} ->
        {:error, result}
    end
  end
end
