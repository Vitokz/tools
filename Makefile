build-scripts:
	cd ./scripts/export-grafan-based-alerts && go build -o export-grafan-based-alerts .
	cd ./scripts/cursor-rules-installer && go build -o cursor-rules-installer . 
