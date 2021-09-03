package lib

has_field(obj, field) {
	_ := obj[field]
}

normalize_to_array(resource) = out_array {
	is_array(resource)
	out_array = resource
} else = out_array {
	out_array = [resource]
}

pick(k, obj1, _) = obj1[k]

pick(k, obj1, obj2) = obj2[k] {
	not has_field(obj1, k)
}

merge_objects(a, b) = c {
	keys := {k | some k; _ = a[k]} | {k | some k; _ = b[k]}
	c := {k: v | k := keys[_]; v := pick(k, b, a)}
}