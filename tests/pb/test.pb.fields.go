package test

var TestFields = struct {
	StringField          string
	NumberField          string
	BoolField            string
	EnumField            string
	MessageField         string
	RepeatedStringField  string
	RepeatedMessageField string
	NumberValueField     string
	StringValueField     string
	BoolValueField       string
	TimeValueField       string
	DurationValueField   string
	OptionalStringField  string
	OptionalNumberField  string
	OptionalBoolField    string
	OptionalEnumField    string
}{
	StringField:          "string_field",
	NumberField:          "number_field",
	BoolField:            "bool_field",
	EnumField:            "enum_field",
	MessageField:         "message_field",
	RepeatedStringField:  "repeated_string_field",
	RepeatedMessageField: "repeated_message_field",
	NumberValueField:     "number_value_field",
	StringValueField:     "string_value_field",
	BoolValueField:       "bool_value_field",
	TimeValueField:       "time_value_field",
	DurationValueField:   "duration_value_field",
	OptionalStringField:  "optional_string_field",
	OptionalNumberField:  "optional_number_field",
	OptionalBoolField:    "optional_bool_field",
	OptionalEnumField:    "optional_enum_field",
}

var SchemaLessFields = struct {
	Data string
}{
	Data: "data",
}
