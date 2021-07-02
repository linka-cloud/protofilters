# Proto Filters

[![Go Reference](https://pkg.go.dev/badge/go.linka.cloud/protofilters.svg)](https://pkg.go.dev/go.linka.cloud/protofilters)

Proto filters provides a simple way to filter protobuf message based on field filter conditions.


## Usage

The two message filtering types available follow the same pattern as `google.protobuf.FieldMask`:

```proto
message FieldsFilter {
  // filters is a map of <field path, Filter>
  map<string, Filter> filters = 1;
}
message FieldFilter {
    string field = 1;
    Filter filter = 2;
}
```

The message's Field is selected by its path and compared against the provided typed filter:

```proto
message Filter {
  oneof match {
    StringFilter string = 1;
    NumberFilter number = 2;
    BoolFilter bool = 3;
    NullFilter null = 4;
    TimeFilter time = 5;
    DurationFilter duration = 6;
  }
  // not negates the match result
  bool not = 7;
}
```

The typed filters are the following:

```proto
message StringFilter {
  message In {
    repeated string values = 1;
  }
  oneof condition {
    string equals = 1;
    string regex = 2;
    In in = 3;
  }
  bool case_insensitive = 4;
}

message NumberFilter {
  message In {
    repeated double values = 1;
  }
  oneof condition {
    double equals = 1;
    double sup = 2;
    double inf = 3;
    In in = 4;
  }
}

message NullFilter {}

message BoolFilter {
  bool equals = 1;
}

message TimeFilter {
  oneof condition {
    google.protobuf.Timestamp equals = 1;
    google.protobuf.Timestamp before = 2;
    google.protobuf.Timestamp after = 3;
  }
}

message DurationFilter {
  oneof condition {
    google.protobuf.Duration equals = 1;
    google.protobuf.Duration sup = 2;
    google.protobuf.Duration inf = 3;
  }
}
```
