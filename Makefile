export GO111MODULE=off

swagger:
	swagger generate spec -o ./swagger.yaml  --scan-models
	
