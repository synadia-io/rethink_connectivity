# Kinda a roundabout way to do this but I'm fine with it for now
schema.go:
	@nats req '$$SRV.SCHEMA.sentiment' "" | jq '.endpoints[0].schema.request | fromjson' > request.json
	@nats req '$$SRV.SCHEMA.sentiment' "" | jq '.endpoints[0].schema.response | fromjson' > response.json
	@cue import -f -p main -l '#Request:' request.json
	@cue import -f -p main -l '#Response:' response.json
	@oapi-codegen -generate types,skip-prune -package main openapi.json
