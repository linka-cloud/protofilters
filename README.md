# Proto Filters

[![Go Reference](https://pkg.go.dev/badge/go.linka.cloud/protofilters.svg)](https://pkg.go.dev/go.linka.cloud/protofilters)

Proto filters provides a simple way to filter protobuf message based on field filter conditions.

**Project status: *alpha***

Not all planned features are completed.
The API, spec, status and other user facing objects are subject to change.
We do not support backward-compatibility for the alpha releases.


## Overview

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

## Usage

Download:

```bash
go get go.linka.cloud/protofilters
```

Basic example:

```go
package main

import (
	"log"

	"google.golang.org/protobuf/types/known/wrapperspb"

	pf "go.linka.cloud/protofilters"
	"go.linka.cloud/protofilters/matcher"
	test "go.linka.cloud/protofilters/tests/pb"
)

func main() {
	m := &test.Test{
		BoolField:      true,
		BoolValueField: wrapperspb.Bool(false),
	}
	ok, err := matcher.MatchFilters(m, &pf.FieldFilter{Field: "bool_field", Filter: pf.True()})
	if err != nil {
		log.Fatalln(err)
	}
	if !ok {
		log.Fatalln("should be true")
	}
	ok, err = matcher.MatchFilters(m, &pf.FieldFilter{Field: "bool_value_field", Filter: pf.False()})
	if err != nil {
		log.Fatalln(err)
	}
	if !ok {
		log.Fatalln("should be true")
	}
}

```

## TODOs

- [ ] support **and/or** conditions
- [ ] support more languages
