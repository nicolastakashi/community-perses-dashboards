apply-dashboards:
	@echo "Applying dashboards"
	@go run main.go | percli apply -f -