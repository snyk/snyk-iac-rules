package rules

import data.lib

is_kind(kind) {
	input.kind == kind
}

services[service] {
	is_kind("Service")
	service = input
}

is_clusterIP_service_type(service) {
	not lib.has_field(service.spec, "type")
} else {
	service.spec.type == "ClusterIP"
}

service_has_externalIPs(service) {
	lib.has_field(service.spec, "externalIPs")
	externalIPs := service.spec.externalIPs
	not is_null(externalIPs)
	count(externalIPs) > 0
}

deny[msg] {
	is_kind("Service")
	services[service]
	is_clusterIP_service_type(service)
	service_has_externalIPs(service)

	msg := {
		"publicId": "CUSTOM-3",
		"title": "Test",
		"subType": "",
		"severity": "critical",
		"issue": "",
		"impact": "",
		"remediation": "",
		"msg": sprintf("spec.externalIPs", []),
		"references": [],
	}
}
