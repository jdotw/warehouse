START_DIR=$PWD

cd ../codegen/pkg/codegen/templates && go generate . && cd ../../.. && go run cmd/oapi-codegen/codegen.go -o ../stock -generate types -cluster stock ../stock/api/stock.yaml
