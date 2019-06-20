# open-callopticum
An analyzer for Asterisk's Call Detail Records

This is an after-hours hobby for me so do not expect me to have fast progress or even to finish something … :) .

The long term goal of this repository is to have an analyzer for Asterisk Call Detail Records.

## Current features

Currently, there are two main features I concentrate on.

### Report Generation

Reports are a great way to obtain concise statistics of your PBX. For example, you can count the number of calls and
produce daily email reports to prevent telephone fraud.

**Note** In some jurisdictions there are privacy laws protecting the data of your employees or customers. Be sure
to comply with the law of your jurisdiction when generating reports. For example, in Europe, GDPR prohibits reports
that concern phoning behaviour of particular employees. Ask your lawyer to be on the safe side.

To generate reports, you need to write two files, a report definition and a template definition. Then a report
can be generated like this:

````
$cd cmd/report
go run report.go -template ~/reportTemplate.gohtml -definition ~/reportDefinition.json /var/log/asterisk/cdrcsv/Master.csv
````

Per default, this command escapes special HTML characters. If you do not want this, pass the `-plain` argument.

First, you need a JSON file containing a definition what you want to count. See [reportDefinition](mockdata/reportDefinition.json)
for an example. The allowed column names can be found in in [parse.go](cdrcsv/parse.go).

Second, you need the actual report template as text file or html file, for example [reportTemplate](mockdata/reportTemplate.tmpl)

The syntax is Go's template engine syntax and with the `{{index .Stats "stat_name"}}` you can access the statistics
defined in your JSON file. See the linked file for more examples on how to use template files. Additionally, you can 
make use of the [Go Template Specification](https://blog.gopheracademy.com/advent-2017/using-go-templates/).
### CDR Pseudonymization

In order to test the other features of Open Callopticum I had to write a tool that replaces personal data in CDR records
by pseudonymified data. After pseudonymization all calls are the same relative to each other, only with other data. This
means that, for example, you can still see that a particular number has called a lot of times, but the number is not the
real number.

You can execute the tool with these commands:

````
$ cd cmd/pseudonymify
$ go run ./pseudo.go -contacts ../mockdata/persons.csv -years 1 -days -20 -minutes 5 -hours -3 -contexts hq,production,support,door,hr /tmp/cdr*.csv
````

The last arguments the csv files to by pseudonymified. The tool replaces all personal data with the data found in the
`persons.csv` provided with the `-contact` argument. It also shifts all calls by the time you specify with the arguments
`-years`, `-days`, `-minutes` and `-hours`. The `-context` parameter lets you define a list of context to be used to
replace the contexts in the real csv.
