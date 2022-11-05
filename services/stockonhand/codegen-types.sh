START_DIR=$PWD

cd ../codegen/pkg/codegen/templates && go generate . && cd ../../.. && go run cmd/oapi-codegen/codegen.go -o ../stockonhand -generate types -cluster stockonhand ../stockonhand/api/stockonhand.yaml
