package rules

import data.lib
import data.lib.testing

test_CUSTOM_3 {
	test_cases := [{
		"want_msgs": ["spec.externalIPs"],
		"fixture": {
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
		},
	}]

	testing.evaluate_test_cases("CUSTOM-3", test_cases)
}
