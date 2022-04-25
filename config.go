package sqly

type Config struct {
    DatabasesFile     string
    ScriptsGlobFiles  string
    Environment       string
    SourceDecryptFunc func(source string) string
}
