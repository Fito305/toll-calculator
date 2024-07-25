obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	# @go build -o bin/receiver data_receiver/main.go 
	## If you are going to build this with main.go 
	# it's going to build this main.go file singlely. 
	# So basically what is oging to happen is if you
	# have other files that main depends on, it's not
	# going to work. So you need to build the whole folder.   
	@./bin/receiver

.PHONY: obu
