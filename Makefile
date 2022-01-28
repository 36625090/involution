run:
	cd src && go run example/main.go --profile dev --log.level trace \
		--log.path ../logs --app example --http.path=/example --ui --http.cors --log.console \
		--newrelic --newrelic.key 0f02f3f4e1f5e995e6e365c6fb08d310f12048fe --http.trace

run-local:
	cd src && go run example/main.go --profile dev --log.level trace \
		--config ../config/config.hcl \
		--log.path ../logs --app example --http.path=/example --ui --http.cors --log.console \
		--newrelic --newrelic.key 0f02f3f4e1f5e995e6e365c6fb08d310f12048fe --http.port 8081
