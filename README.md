# Swapdoodle HPP server

An ongoing attempt at creating a replacement Swapdoodle server.  
This may succeed, or go absolutely nowhere. I don't even know Go lol  

[PretendoNetwork/super-mario-maker](https://github.com/PretendoNetwork/super-mario-maker) was used as a base for this code. I hope that's alright!

## Status

WE BALLIN, THIS MIGHT GET SOMEWHERE!!! (read: i understand almost everything and cant wait to implement it. if only it werent finals period at school)

## Configuration

All configuration options are handled via environment variables

`.env` files are supported

| Name                                | Description                                                           | Required                                      |
|-------------------------------------|-----------------------------------------------------------------------|-----------------------------------------------|
| `PN_SD_POSTGRES_URI`               | Fully qualified URI to your Postgres server                           | Yes                                           |
| `PN_SD_HPP_SERVER_PORT`         | Port for the secure server                                            | Yes                                           |
| `PN_SD_CONFIG_S3_ENDPOINT`         | S3 server endpoint                                                    | Yes                                           |
| `PN_SD_CONFIG_S3_ACCESS_KEY`       | S3 access key ID                                                      | Yes                                           |
| `PN_SD_CONFIG_S3_ACCESS_SECRET`    | S3 secret                                                             | Yes                                           |
| `PN_SD_CONFIG_S3_BUCKET`           | S3 bucket                                                             | Yes                                           |
| `PN_SD_ACCOUNT_GRPC_HOST`          | Host name for your account server gRPC service                        | Yes                                           |
| `PN_SD_ACCOUNT_GRPC_PORT`          | Port for your account server gRPC service                             | Yes                                           |
| `PN_SD_ACCOUNT_GRPC_API_KEY`       | API key for your account server gRPC service                          | No (Assumed to be an open gRPC API)           |
