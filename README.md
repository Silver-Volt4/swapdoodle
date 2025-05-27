# ARCHIVED - THIS REPOSITORY WILL NOT ACCEPT CONTRIBUTIONS

Development has been moved to [SwapdoodleRevival/swapdoodle](https://github.com/SwapdoodleRevival/swapdoodle), which is a proper merge of [PretendoNetwork/swapdoodle](https://github.com/PretendoNetwork/swapdoodle) and will be (hopefully) merged there at some point. The code here is identical to [SwapdoodleRevival/swapdoodle - commit 8a1f568](https://github.com/SwapdoodleRevival/swapdoodle/commit/8a1f5687e4bf5f2ee0ae1e925fb91e240fbce376). No new changes will be made here.

# Swapdoodle HPP server

A sorta-working replacement server for Swapdoodle

## Status

The server should be considered alpha-quality at the moment. It has worked in our private testing, but broader public testing will be needed.

## Credits

[PretendoNetwork/super-mario-maker](https://github.com/PretendoNetwork/super-mario-maker), it was used as a base for this code. I hope that's alright, I don't know Go much...  
[milesthecreator.bsky.social](https://bsky.app/profile/milesthecreator.bsky.social), who supplied us with unencrypted Swapdoodle network dumps that have significantly helped development.  

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
