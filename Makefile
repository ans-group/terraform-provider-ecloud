testacc:
	TF_ACC=1 \
		go test -v \
		./ecloud \
		-timeout 120m \
		-run=TestAcc${TEST}

testacc-all:
	TF_ACC=1 \
		go test -v \
		./ecloud \
		-timeout 120m