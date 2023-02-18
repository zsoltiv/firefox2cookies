package main

import (
    "fmt"
    "log"
    "os"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func usage() {
    fmt.Println("Usage: " + os.Args[0] + " /path/to/firefox-profile/cookies.sql [output_name]")
}

func stringBool(expr bool) string {
    if expr {
        return "TRUE"
    } else {
        return "FALSE"
    }
}

func stringInt(expr int) string {
    if expr == 1 {
        return "TRUE"
    } else {
        return "FALSE"
    }
}

func rowToLine(row *sql.Rows) (string, error) {
    var (
        host, path, subdomains, isSecure, name, value string
        expiry, secure int
    )
    
    if err := row.Scan(&host, &path, &name, &value, &secure, &expiry); err != nil {
        return "", err
    }
    isSecure = stringInt(secure)
    subdomains = stringBool(host[0] == '.')

    const format string = "%s\t%s\t%s\t%s\t%d\t%s\t%s\n"

    if subdomains == "TRUE" {
        return fmt.Sprintf(format, host[1:], subdomains, path, isSecure, expiry, name, value), nil
    } else {
        return fmt.Sprintf(format, host, subdomains, path, isSecure, expiry, name, value), nil
    }
}

func main() {
    if len(os.Args) < 2 {
        usage()
        return
    }

    db, err := sql.Open("sqlite3", os.Args[1])
    if err != nil {
        log.Fatal(err.Error())
        return
    }
    defer db.Close()
    err = db.Ping()
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    var txt string
    if len(os.Args) == 3 {
        txt = os.Args[2]
    } else {
        txt = "cookies.txt"
    }

    file, err := os.Create(txt)
    if err != nil {
        log.Fatal(err.Error())
        return
    }
    defer file.Close()

    rows, err := db.Query("SELECT host, path, name, value, isSecure, expiry FROM moz_cookies")
    if err != nil {
        log.Fatal(err.Error())
        return
    }
    defer rows.Close()

    for rows.Next() {
        line, err := rowToLine(rows)
        if err != nil {
            log.Fatal(err.Error())
        }
        file.WriteString(line)
    }
}
