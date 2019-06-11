# open-callopticum
An analyzer for Asterisk's Call Detail Records

This is an after-hours hobby for me so do not expect me to have fast progress or even to finish something … :) .

The long term goal of this repository is to have an analyzer for Asterisk Call Detail Records.

## Current state

The work is very, very early in progress and currently one can only pseudonymify real cdrs. I implemented this feature at first to have reliable test data that does not reveal private data.

Here is a sample call to pseudonymify a bunch of cdr files:

````
$ cd cmd
$ go run ./pseudo.go -contacts ../mockdata/persons.csv -years 1 -days -20 -minutes 5 -hours -3 -contexts hq,production,support,door,hr /tmp/cdr*.csv
````

The above command pseudonymifies all csv files beginning with `cdr` in the `/tmp`-directory of your machine. The contacts paramener points to a csv file where the program can find pseudo contacts needed to replace the real names and phone numbers. Of course you can create your own csv contact file with the same format. Beware that neither the cdr files nor the pseudo contacts file have csv headers.

A description of the other command line parameters can be found by executing:

````
$ go run ./pseudo.go -h
````
