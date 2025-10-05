build-scripts:
	cd ./scripts/export-grafan-based-alerts && go build -o export-grafan-based-alerts .
	cd ./scripts/repo-init && go build -o repo-init . 
