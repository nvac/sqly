package sqly

import (
    "encoding/xml"
    "errors"
    "fmt"
    "github.com/jmoiron/sqlx"
    "io/ioutil"
    "path/filepath"
    "time"
)

func getScriptByName(name string) (*scriptsCacheValue, error) {
    scriptsValue, ok := scriptsCache.Load(name)
    if !ok {
        return nil, errors.New("not found script")
    }

    script := scriptsValue.(scriptsCacheValue)

    return &script, nil
}

func getDatabaseByName(databaseName string) (*databasesCacheValue, error) {
    databasesCacheKey := getDatabasesCacheKey(databaseName)
    databasesValue, ok := databasesCache.Load(databasesCacheKey)
    if !ok {
        return nil, errors.New("not found database")
    }

    database := databasesValue.(databasesCacheValue)
    if database.ping {
        err := database.db.Ping()
        if err != nil {
            return nil, err
        }

        database.ping = true
    }

    return &database, nil
}

func loadDatabasesFile() error {
    if globalConfig.DatabasesFile == "" {
        return errors.New("miss DatabasesFile")
    }
    file, err := ioutil.ReadFile(globalConfig.DatabasesFile)
    if err != nil {
        return err
    }

    data := new(databases)
    err = xml.Unmarshal(file, data)
    if err != nil {
        return err
    }

    if data == nil || data.Databases == nil {
        return errors.New("no available database")
    }

    var cacheDatabase []string
    for _, database := range data.Databases {
        if globalConfig.Environment != "" && database.Environment != globalConfig.Environment {
            continue
        }

        source := database.Source
        if globalConfig.SourceDecryptFunc != nil {
            source = globalConfig.SourceDecryptFunc(database.Source)
        }

        db, err := sqlx.Open(database.Driver, source)
        if err != nil {
            return err
        }

        if database.MaxOpenConns != nil {
            db.SetMaxOpenConns(*database.MaxOpenConns)
        }

        if database.MaxIdleConns != nil {
            db.SetMaxIdleConns(*database.MaxIdleConns)
        }

        if database.ConnMaxLifetime != nil {
            seconds := time.Duration(*database.ConnMaxLifetime)
            db.SetConnMaxLifetime(seconds * time.Second)
        }

        if database.ConnMaxIdleTime != nil {
            seconds := time.Duration(*database.ConnMaxIdleTime)
            db.SetConnMaxIdleTime(seconds * time.Second)
        }

        databasesCacheKey := getDatabasesCacheKey(database.Name)
        value := databasesCacheValue{
            db:   db,
            name: database.Name,
            ping: false,
        }
        cacheDatabase = append(cacheDatabase, databasesCacheKey)
        databasesCache.Store(databasesCacheKey, value)
    }

    lintDatabasesCache(cacheDatabase)

    return nil
}

func loadScriptsGlobFile(path string) error {
    fileContent, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }

    data := new(scripts)
    err = xml.Unmarshal(fileContent, data)
    if err != nil {
        return err
    }

    for _, script := range data.Scripts {
        key := script.Name
        value := scriptsCacheValue{
            name:    script.Name,
            content: script.Content,
            path:    path,
        }

        if loadedValue, loaded := scriptsCache.LoadOrStore(key, value); loaded {
            loadedValue := loadedValue.(scriptsCacheValue)
            message := fmt.Sprintf("The duplicate script name('%s') in the current file('%s') was found in another file('%s')",
                loadedValue.name, path, loadedValue.path)
            return errors.New(message)
        }
    }
    return nil
}

func loadScriptsGlobFiles() error {
    files, _ := filepath.Glob(globalConfig.ScriptsGlobFiles)
    for _, file := range files {
        if err := loadScriptsGlobFile(file); err != nil {
            return err
        }
    }
    return nil
}

func contains(values []string, target string) bool {
    for _, value := range values {
        if value == target {
            return true
        }
    }
    return false
}

func connect(databaseName string) (*DB, error) {
    if !globalInit {
        return nil, errors.New("sqly has not been initialized, please use 'sqly.Init(&sqly.Config{})' to initialize")
    }

    databaseCache, err := getDatabaseByName(databaseName)
    if err != nil {
        return nil, err
    }

    if databaseCache.ping {
        err := databaseCache.db.Ping()
        if err != nil {
            return nil, err
        }
        databaseCache.ping = true
    }

    return &DB{DB: databaseCache.db}, nil
}
