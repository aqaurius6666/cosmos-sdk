go 1.17

module github.com/aqaurius6666/cosmos-sdk

// latest grpc doesn't work with with our modified proto compiler, so we need to enforce
// the following version across all dependencies.
replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76

// Fix upstream GHSA-h395-qcrw-5vmq vulnerability.
// TODO Remove it: https://github.com/cosmos/cosmos-sdk/issues/10409
replace github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.0

replace github.com/cosmos/cosmos-sdk/db => ./db

replace github.com/cosmos/cosmos-sdk/x/group => ./x/group

replace github.com/cosmos/cosmos-sdk/errors => ./errors

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/tendermint/tendermint v0.35.0
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
)
