package jsonutils

import "github.com/ghodss/yaml"

func ParseYAML(str string) (JSONObject, error) {
	jsonBytes, err := yaml.YAMLToJSON([]byte(str))
	if err != nil {
		return nil, err
	}
	return Parse(jsonBytes)
}

func yamlString(obj JSONObject) string {
	yamlBytes, _ := yaml.JSONToYAML([]byte(obj.String()))
	return string(yamlBytes)
}

func (th *JSONValue) YAMLString() string {
	return yamlString(th)
}

func (th *JSONString) YAMLString() string {
	return yamlString(th)
}

func (th *JSONInt) YAMLString() string {
	return yamlString(th)
}

func (th *JSONFloat) YAMLString() string {
	return yamlString(th)
}

func (th *JSONBool) YAMLString() string {
	return yamlString(th)
}

func (th *JSONArray) YAMLString() string {
	return yamlString(th)
}

func (th *JSONDict) YAMLString() string {
	return yamlString(th)
}
