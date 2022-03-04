# Wasabi Bug Sample

This project is a proof of concept showing that supplying a UserMetadata value
that contains two consecutive spaces causes a PutObject operation to fail when
uploading files to Wasabi. These uploads do not fail on Amazon S3 or Minio 
server.

This POC demonstrates APTrust preservation-services 
[bug #57](https://github.com/APTrust/preservation-services/issues/57).

## Usage

1. Set the environment vars WASABI_ACCESS_KEY and WASABI_SECRET_KEY to a 
keypair that is valid for your Wasabi bucket.

2. Choose a Wasabi bucket for upload. In the example command, we'll call it
my-bucket.

3. Run the following command:

```
go run main.go my-bucket
```

If you get an error about missing environment values, try one of the following
commands. The first passes current env vars into the process created by 
`go run`:

```
WASABI_ACCESS_KEY=$WASABI_ACCESS_KEY WASABI_SECRET_KEY=$WASABI_SECRET_KEY go run main.go my-bucket
```

This option sets the literal key pair values

```
WASABI_ACCESS_KEY=<access_key> WASABI_SECRET_KEY=<secret_key> go run main.go my-bucket
```

## What Happens

The program uploads the same text file twice. The first upload sets the header
`x-amz-meta-custom-data` to value `Metadata values with single spaces are OK`.
This upload succeeds.

The second upload sets the same header to 
`Metadata values with two consecutive  spaces cause upload to fail`. This 
upload fails.

Here's the output:

```
Upload with goodPutOptions SUCCEEDED
Upload with badPutOptions FAILED with error: The request signature we calculated does not match the signature you provided. Check your key and signing method.
```

It seems any file that includes a UserMetadata value with multiple consecutive
spaces returns this signature mismatch error. Again, these files upload without
error to Amazon S3 and to Minio servers.
