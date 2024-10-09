package container

import (
    "testing"
)

func TestDockerVersion(t *testing.T) {
    result := ApiVersionRegex.Match("4.9.4-rhel")
    fmt.Print(result)
}
