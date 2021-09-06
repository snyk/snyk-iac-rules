#!/bin/bash
Describe 'go run main.go parse ./fixtures/custom-rules/rules/CUSTOM-1/fixtures/test.tf'
   It 'returns passing test status'
      When call go run main.go parse ./fixtures/custom-rules/rules/CUSTOM-1/fixtures/test.tf
      The status should be success
      The output should include '{
	"resource": {
		"aws_security_group": {
			"denied": {
				"description": "Allow SSH inbound from anywhere",
				"ingress": {
					"cidr_blocks": [
						"0.0.0.0/0"
					],
					"from_port": 22,
					"protocol": "tcp",
					"to_port": 22
				},
				"name": "allow_ssh",
				"vpc_id": "${aws_vpc.main.id}"
			}
		}
	}
}'
   End
End

Describe 'go run main.go parse ./fixtures/custom-rules/rules/CUSTOM-3/fixtures/test.yaml --format yaml'
   It 'returns passing test status'
      When call go run main.go parse ./fixtures/custom-rules/rules/CUSTOM-3/fixtures/test.yaml --format yaml
      The status should be success
      The output should include '{
	"apiVersion": "v1",
	"kind": "Service",
	"metadata": {
		"name": "using-external-ip-default"
	},
	"spec": {
		"externalIPs": [
			"1.1.1.1"
		],
		"ports": [
			{
				"name": "http",
				"port": 80,
				"protocol": "TCP",
				"targetPort": 8080
			}
		],
		"selector": {
			"app": "MyApp"
		}
	}
}'
   End
End
