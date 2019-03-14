# rescope

rescope is a tool (Go) that lets you quickly define scopes in Burp/ZAP - mainly intended for "bug hunters" and pentesters who deal with large scopes. See [blog post](https://root4loot.com/post/introducing_rescope/) for context/background.

Simply provide a scope (file containing target identifiers) and rescope parses this to a file format that can be imported from Burp/ZAP directly.

> The latest version (0.2) takes care of several issues and introduces some welcomed changes. See [CHANGELOG](CHANGELOG.md) for details. Upgrade instructions below.

## Features 

- Identifies targets from scope (give it any structure/format)
- Set excludes in addition to includes
- Parse multiple scope-files at once (combine scopes)
- Takes care of wildcards, files, protocols, ports (saves you from messing with regex)
- Supports parsing IP ranges/CIDR (aside from domains/hosts)

## Installation

Requires [Go](https://golang.org/doc/install#install) (tested on 1.11.4).  
```
go get github.com/root4loot/rescope
```

## Compiling and running
Compiling is easy with Go.
```
go install github.com/root4loot/rescope
```
By default, Go saves binaries to `$GOPATH/bin/` which is typically `~/go/bin/` for Unix or `%USERPROFILE%\go` on Windows, unless changed.

Once compiled, you can simply create a soft link from a desired location. E.g. Unix:
```
ln -s ~/go/bin/rescope /usr/local/bin/rescope
```
Running:
```
rescope --version
```

## Upgrading
```
go get -u github.com/root4loot/rescope
```
If you run into merge conflicts, then delete the repo and install once again. Sorry for the inconvenience.

Unix:
```
// If $GOPATH is set
rm -rf $GOPATH/src/github.com/root4loot/rescope

// If $GOPATH is not set, try
rm -rf ~/go/src/github.com/root4loot.com/rescope
```
Windows:
```
// If %GOPATH% is set
rd /s /q "%GOPATH%\src\github.com\root4loot\rescope"

// If %GOPATH% is not set, try
rd /s /q "%USERPROFILE%\go\src\github.com\root4loot\rescope"
```
Then install:
```
go get github.com/root4loot/rescope
```

## Usage
```
usage: rescope [[-z|--zap | [-b]--burp] [-i|--infile "<filepath>" ...]] [-o|--outfile "<filepath>"]]] [-n|--name "<value>"] [-e|--extag "<value>"] [-s|--silent] [-h|--help] [--version]]
```

### Arguments

| Short | Long | Description   | 
| :------------- |:-------------| :-----  | 
| -h | --help     | Print help information |
| -b | --burp     | Parse to Burp Suite JSON (required option) |
| -z | --zap      | Parse to OWASP ZAP XML (required option) |
| -i | --infile   | File (scope) to be parsed (required) | 
| -o | --outfile  | File to write parsed results (required) |
| -s | --silent   | Do not print identified targets |
| -n | --name     | Name of ZAP context |
| -e | --extag    | Custom exclude tag (default: !EXCLUDE) |
|    | --version  | Print version |

-----

### Example Usage

Parse scope to Burp Suite compatible JSON:
```
rescope --burp --infile scope.txt --outfile burp.json
rescope -b -i scope.txt -o burp.json
```

Parse scope to ZAP compatible XML:
```
rescope --zap --infile scope.txt --outfile zap.context
rescope -z -i scope.txt -o zap.context
```

Parse multiple scopes to ZAP XML, set context name, silence output:
```
rescope --zap -i scope1.txt -i scope2.txt -o zap.context -n CoolScopeName --silent
```


### Setting Excludes

rescope treats identified targets as Includes by default.  
To set Excludes, specify an **!EXCLUDE** tag anywhere in the document, followed by the targets you wish to exclude.  
If this tag does not work for you, then a custom one can be set from the `--extag (-e)` parameter.


Example:

```
// include these
prod.example.com
admin.example.com

!EXCLUDE
// exclude these
dev.example.com
test.example.com
```


## Example

rescope will attempt to identify target-identifiers from the scope(s) you provide. This enables you to quickly copy/paste the scope section from various places to a file and serve this directly to rescope without having to do much edits in prior. It doesn’t matter what comes before and after the identifiers, as long as they’re there. 

Consider the following scope having both **in-scope** and **out-of-scope** targets:
```
$ cat scope.txt
In Scope:
Critical http://admin.example.com/login.aspx
Critical https://example.com:8080/upload/*
Critical *.dev.example.com and *.prod.example.com
High     10.10.10.1-2 (testing)

Out of Scope:
bgp.example.com:179
*.vendor.example.com/assets/
ftp://10.10.10.1:21
```

As you can see, most of these identifiers have leading text/whitespace and so on. But that's totally fine!  
The only thing we have to do in this case, is to specify an **!EXCLUDE** tag _before_ the "Out of Scope" identifiers.

```diff
In Scope:
Critical http://admin.example.com/login.aspx
Critical https://example.com:8080/upload/*
Critical *.dev.example.com and *.prod.example.com
High     10.10.10.1-2 (testing)

+ !EXCLUDE
Out of Scope:
bgp.example.com:179
*.vendor.example.com/assets/
ftp://10.10.10.1:21
```
Having saved this, we're ready to parse and import results to either Burp Suite or ZAP.


### Parsing to Burp Suite JSON

Parsing scope to Burp JSON is easy.

```diff
$ rescope --burp -i scope.txt -o burp.json
[-] Grabbing targets from [scope.txt]
+ http://admin.example.com/login.aspx
+ https://example.com:8080/upload/*
+ *.dev.example.com
+ *.prod.example.com
+ 10.10.10.1-2
- bgp.example.com:179
- *.vendor.example.com/assets/
- ftp://10.10.10.1:21
[-] Parsing to JSON (Burp Suite)
[✓] Done
[✓] Wrote 1732 bytes to burp.json
```
rescope will highlight Includes in Green and Excludes in Red, unless `--silent (-s)` was set.

#### Parsed results

See [importing to Burp](#to-burp-suite)

```
$ cat burp.json 
{
  "target": {
    "scope": {
      "advanced_mode": true,
      "exclude": [
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^bgp\\.example\\.com$",
          "port": "^179$",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^/assets/[\\S]*$",
          "host": "^[\\S]*\\.vendor\\.example\\.com$",
          "port": "",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^10\\.10\\.10\\.1$",
          "port": "^21$",
          "protocol": "Any"
        }
      ],
      "include": [
        {
          "enabled": true,
          "file": "^/login\\.aspx$",
          "host": "^admin\\.example\\.com$",
          "port": "^80$",
          "protocol": "http"
        },
        {
          "enabled": true,
          "file": "^/upload/[\\S]*$",
          "host": "^example\\.com$",
          "port": "^(8080|443)$",
          "protocol": "https"
        },
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^[\\S]*\\.dev\\.example\\.com$",
          "port": "",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^[\\S]*\\.prod\\.example\\.com$",
          "port": "",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^10\\.10\\.10\\.1$",
          "port": "",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^10\\.10\\.10\\.2$",
          "port": "",
          "protocol": "Any"
        }
      ]
    }
  }
}
```

### Parsing to OWASP ZAP XML

Parsing scope to ZAP XML is just as easy as with Burp.
However, there are a couple of things to note:

* ZAP requires all scopes (contexts) to have unique names. Names can be set with the `--name` or `-n` parameter. If no name is given, then rescope will ask for one.
* The default file extension for importing/exporting scopes in ZAP is **.context**. Although not required, it is advised to use this exension for the outfile, to avoid having to take unnecessary steps when importing it later. 

```diff
$ rescope --zap -i scope.txt -o zap.context -n CoolScopeName
[-] Grabbing targets from [scope.txt]
+ http://admin.example.com/login.aspx
+ https://example.com:8080/upload/*
+ *.dev.example.com
+ *.prod.example.com
+ 10.10.10.1-2
- bgp.example.com:179
- *.vendor.example.com/assets/
- ftp://10.10.10.1:21
[-] Parsing to XML (OWASP ZAP)
[✓] Done
[✓] Wrote 2154 bytes to zap.context
```

#### Parsed results

See [importing to ZAP](#to-owasp-zap)

rescope uses the default ZAP context as a template for creating new ones, meaning it'll include the default "technologies" as well. This can be easily removed from the application once the scope/context has been set.

```
$ cat zap.context | head -n 45
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<configuration>
<context>
<name>CoolScopeName</name>
<desc/>
<inscope>true</inscope>
<incregexes>^http:\/\/admin\.example\.com\/login\.aspx[\S]*$</incregexes>
<incregexes>^https:\/\/example\.com:8080\/upload\/[\S]*[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/[\S]*\.dev\.example\.com[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/[\S]*\.prod\.example\.com[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/10\.10\.10\.1[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/10\.10\.10\.2[\S]*$</incregexes>
<excregexes>^http(s)?:\/\/bgp\.example\.com:179[\S]*$</excregexes>
<excregexes>^http(s)?:\/\/[\S]*\.vendor\.example\.com\/assets\/[\S]*$</excregexes>
<excregexes>^http(s)?:\/\/ftp:\/\/10\.10\.10\.1:21[\S]*$</excregexes>
<tech>
<include>Db</include>
<include>Db.Firebird</include>
<include>Db.HypersonicSQL</include>
<include>Db.IBM DB2</include>
<include>Db.Microsoft Access</include>
<include>Db.Microsoft SQL Server</include>
<include>Db.MySQL</include>
<include>Db.Oracle</include>
<include>Db.PostgreSQL</include>
<include>Db.SAP MaxDB</include>
<include>Db.SQLite</include>
<include>Db.Sybase</include>
<include>Language</include>
<include>Language.ASP</include>
<include>Language.C</include>
<include>Language.PHP</include>
<include>Language.XML</include>
<include>OS</include>
<include>OS.Linux</include>
<include>OS.MacOS</include>
<include>OS.Windows</include>
<include>SCM</include>
<include>SCM.Git</include>
<include>SCM.SVN</include>
<include>WS</include>
<include>WS.Apache</include>
<include>WS.IIS</include>
<include>WS.Tomcat</include>
</tech>
```

## Importing

### To Burp Suite
1. Select **Target** pane
2. Select **Scope** pane
3. Click the gear (⚙︎) icon 
4. Select **Load options**
5. Select outputted JSON file from rescope

### To OWASP ZAP
**File** -> **Import Context** -> Select outputted XML file from rescope

Note: If you set `-o` filename ext to anything other than `.context` then you'll have to set 'File Format:' to 'All' (in file select).


## Disclaimer
- This is my first project in Go (and I don't consider myself a developer) so bear that in mind as far as the code goes.
- rescope may (without my knowledge) identify or parse scope-identifiers inaccurately. Therefore you should probably go over the results yourself and make sure it got what you wanted. If not then please [submit an issue](https://github.com/root4loot/rescope/issues). Alternatively you can always find me on [Twitter](https://twitter.com/root4loot).


## Author
* Daniel Antonsen (root4loot)

## License
Licensed under MIT (see file **LICENSE**)
