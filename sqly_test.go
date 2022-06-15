package sqly

import "testing"

func Test(t *testing.T) {
    config := Config{
        DatabasesFile:    "config/databases.xml",
        ScriptsGlobFiles: "config/sql/*.xml",
        Environment:      "debug",
        SourceDecryptFunc: func(source string) string {
            return source
        },
    }
    err := Init(config)
    if err != nil {
        panic(err)
    }
}
