# CDRO-CLI

## Summary

This is a lightweight command line tool to allow you to apply CDRO DSL files directly from the command line on your workstation or as a part of a CI process.

## Usage

*  `-conf` 
        Config file (default "NONE")
*  `-file` 
        Groovy or YAML file to run (default "ERROR")
*  `-password` 
        Password
*  `-type` 
        DSL or YAML files (default "groovy")
*  `-url` 
        CD/RO URL
*  `-username` 
        Username

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

## Config File

The configuration file is a `json` file with configuration information on how to connect to the CD/RO server.  An example configuration file could look as follows:

```
{
    "username": "admin",
    "password": "SuperSecretPassword",
    "url": "https://flow-web.cloudbees.com/"
}
```



