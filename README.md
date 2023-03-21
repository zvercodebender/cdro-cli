# CDRO-CLI

[![Github All Releases](https://img.shields.io/github/downloads/zvercodebender/cdro-cli/total.svg)]()
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/69125833c2ab49a8b79c97c31284419a)](https://app.codacy.com/gh/zvercodebender/cdro-cli/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

## Summary

This is a lightweight command line tool to allow you to apply CDRO DSL files directly from the command line on your workstation or as a part of a CI process.

## Usage

* `--conf` 
        Config file (default "NONE")
* `--file` 
        Groovy or YAML file to run (default "ERROR")
* `--password` 
        Password
* `--type` 
        DSL or YAML files (default "groovy")
* `--url` 
        CD/RO URL
* `--username` 
        Username
* `--test`
        Don't apply
* `--verbose` 
        Show extra output
* `--value <key>=<value>`
        Pass values into the script to replace placeholders

## Special File Tags

In addition to command line options you can use custom tags in your scripts as follows:

* `@@IncludeFile: <filename>@@` This will insert another file at this point in your main file.  An example could look like the following:

*test.groovy*

```
println "Hello Earth"
println "Hello Mars"
@@IncludeFile: cruft/start.groovy@@
```

*start.groovy*
```
println "Hello Moon"
```

The resulting script would be as follows:

```
println "Hello Earth"
println "Hello Mars"
println "Hello Moon"
```

* `@@IncludeValue: <key>@@` This will be replace by `values` passed in a command line arguments with the `--values <key>=<value>` command line switch.  An example of this is as follows:

```
println "Hello World"
println "Hello Mars"
println "Hello @@IncludeValue: message@@"
```

user the command line switch `--value message=Bob`.  The resulting script will be as follows:

```
println "Hello World"
println "Hello Mars"
println "Hello Bob"
```

## Config File

The configuration file is a `json` file with configuration information on how to connect to the CD/RO server.  An example configuration file could look as follows:

```
{
    "username": "admin",
    "password": "SuperSecretPassword",
    "url": "https://flow-web.cloudbees.com/"
}
```



