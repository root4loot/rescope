# rescope

rescope is a tool (Go) that lets you quickly define scopes in Burp/ZAP - mainly intended for "bug hunters" and pentesters who deal with large scopes. See [blog post](https://root4loot.com/post/introducing_rescope/) for context/background.

Simply provide a scope (file containing target identifiers) and rescope parses this to a format which can be imported from Burp/ZAP.

<img src="https://root4loot.com/img/gif/rescope_min.gif" width="653" height="418">

**Features**

- Identifies targets from scope (give it any structure/format)
- Set excludes (aside from includes)
- Parse multiple scope-files at once
- Supports parsing IP-ranges/CIDR (aside from domains/hosts)

Disclaimer: This is my first project in Go (and I'm not really a programmer) so bear that in mind as far as the code goes. 

## Installation

Requires [Go](https://golang.org/doc/install#install) (tested on 1.11.4)

```
go get github.com/root4loot/rescope
```

Compiling:
```
go install github.com/root4loot/rescope
```
Compiled binaries are saved to `$GOPATH/bin/` by default.

## Usage


```
usage: rescope [[burp|zap] [-i|--infile "<value>" ...]] [-o|--outfile "<value>"]] [-n|--name "<value>"] [-e|--extag "<value>"] [-s|--silent] | [-h|--help] [--version]
```

### Arguments

| Short | Long | Description   | 
| :------------- |:-------------| :-----  | 
| -h | --help     | Print help information |
| -i | --infile   | File (scope) to be parsed (required) | 
| -o | --outfile  | File to write parsed results (required) |
| -s | --silent   | Do not print identified targets |
| -n | --name     | Name of ZAP context |
| -e | --extag    | Custom exclude tag (default: !EXCLUDE) |
|    | --version  | Print version |

-----

### Example Usage

Parse scope to Burp Suite compatible JSON
```
rescope burp -i scope.txt -o burp.json
```

Parse scope to ZAP compatible XML
```
rescope zap -i scope.txt -o zap.context
```

Parse multiple scopes to ZAP XML, set context name, silence output
```
rescope zap -i scope1.txt -i scope2.txt -o zap.context --name CoolScope  --silent
```



### Setting Excludes

rescope treats identified targets as Includes by default.
To set Excludes, specify an **!EXCLUDE** tag anywhere in the document, followed by the targets you wish to exclude. Alternatively, you can set a custom tag via the `--extag (-e)` parameter.


Example:

```
// include this
prod.example.com
admin.example.com

!EXCLUDE

// exclude this
dev.example.com
test.example.com
```



## Example

rescope will __do its best__ to identify targets from the scopes you provide. This enables you to quickly copy/paste the scope section from various places to a file and serve this directly to rescope without having to do much edits in prior. It doesn’t matter what comes before and after the target-identifiers, as long as they’re there.

Consider the following scope having both **in-scope** and **out-of-scope** targets:
```
$ cat scope.txt
In Scope:
Critical admin.example.com/login.aspx
Critical https://example.com/upload:8080
Critical *.dev.example.com and *.prod.example.com
High 192.168.0.1-2 (internal testing)

Out of Scope:
bgp.example.com:179
*.vendor.example.com
192.168.10.9
```

As you probably noticed, most identifiers have leading text/whitespace and so on but this shouldn't be a problem.
The only thing we need to do in this case is to specify a **!EXCLUDE** tag before the "out-of-scope" list.

```diff
In Scope:
Critical admin.example.com/login.aspx
Critical https://example.com/upload:8080
Critical *.dev.example.com and *.prod.example.com
High 192.168.0.1-2 (internal testing)

+ !EXCLUDE
Out of Scope:
bgp.example.com:179
*.vendor.example.com
192.168.10.9
```

Having saved this, we're ready to parse and import results to either Burp Suite or ZAP.

### Parsing to Burp Suite JSON

Parsing scope to Burp JSON is easy.

```diff
$ rescope burp --infile scope.txt --outfile burp.json
[-] Grabbing targets from [scope.txt]
+ admin.example.com/login.aspx
+ https://example.com/upload:8080
+ *.dev.example.com
+ *.prod.example.com
+ 192.168.0.1-2
- bgp.example.com:179
- *.vendor.example.com
- 192.168.10.9
[-] Parsing to JSON (Burp Suite)
[✓] Done
[✓] Wrote 1696 bytes to burp.json
```

rescope highlights Includes in Green and Excludes in Red, unless `--silent (-s)`

Important: rescope may not always parse or identify targets accurately. Therefore it's crucial that you go over the results and make sure it got what you wanted.


#### Parsed results

See [Importing to Burp](#to-burp-suite)

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
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^[\\S]*\\.vendor\\.example\\.com$",
          "port": "",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^[\\S]*$",
          "host": "^192\\.168\\.10\\.9$",
          "port": "",
          "protocol": "Any"
        }
      ],
      "include": [
        {
          "enabled": true,
          "file": "^/login\\.aspx\\/?[\\S]*$",
          "host": "^admin\\.example\\.com$",
          "port": "",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "^/upload:8080\\/?[\\S]*$",
          "host": "^example\\.com$",
          "port": "443",
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
          "file": "",
          "host": "192.168.0.1",
          "port": "",
          "protocol": "Any"
        },
        {
          "enabled": true,
          "file": "",
          "host": "192.168.0.2",
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
However, there are a few things to note:

* ZAP requires all scopes (contexts) to have unique names. Names can be set with the `--name` or `-n` parameter. If no name is given, then rescope will ask for one.
* Although not required, it is advised to use `.context` as the file extension for the outfile to avoid taking unnecessary steps when importing it later.

```diff
$ rescope zap --name CoolScope --infile example.txt --outfile zap.context
[-] Grabbing targets from [example.txt]
+ admin.example.com/login.aspx
+ https://example.com/upload:8080
+ *.dev.example.com
+ *.prod.example.com
+ 192.168.0.1-2
- bgp.example.com:179
- *.vendor.example.com
- 192.168.10.9
[-] Parsing to XML (OWASP ZAP)
[✓] Done
[✓] Wrote 2130 bytes to zap.context
```

#### Parsed results
See [Importing to ZAP](#to-owasp-zap)

Nope: rescope uses the default ZAP context as a template for creating new scopes, meaning it'll include the standard "technologies" as well. This can be easily removed within the application once the scope/context is set.

```
$ cat zap.context 
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<configuration>
<context>
<name>MyScope</name>
<desc/>
<inscope>true</inscope>
<incregexes>^http(s)?:\/\/admin\.example\.com\/login\.aspx[\S]*$</incregexes>
<incregexes>^https:\/\/example\.com\/upload:8080[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/[\S]*\.dev\.example\.com[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/[\S]*\.prod\.example\.com[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/192\.168\.0\.1[\S]*$</incregexes>
<incregexes>^http(s)?:\/\/192\.168\.0\.2[\S]*$</incregexes>
<excregexes>^http(s)?:\/\/bgp\.example\.com:179[\S]*$</excregexes>
<excregexes>^http(s)?:\/\/[\S]*\.vendor\.example\.com[\S]*$</excregexes>
<excregexes>^http(s)?:\/\/192\.168\.10\.9[\S]*$</excregexes>
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
<urlparser>
<class>org.zaproxy.zap.model.StandardParameterParser</class>
<config>{"kvps":"&amp;","kvs":"=","struct":[]}</config>
</urlparser>
<postparser>
<class>org.zaproxy.zap.model.StandardParameterParser</class>
<config>{"kvps":"&amp;","kvs":"=","struct":[]}</config>
</postparser>
<authentication>
<type>0</type>
</authentication>
<forceduser>-1</forceduser>
<session>
<type>0</type>
</session>
<authorization>
<type>0</type>
<basic>
<header/>
<body/>
<logic>AND</logic>
<code>-1</code>
</basic>
</authorization>
</context>
</configuration>
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

## Author

* Daniel Antonsen (root4loot)

## License

Licensed under MIT (see file **LICENSE**)
