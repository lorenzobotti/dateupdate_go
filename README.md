# dateupdate_go

Every non-Ubuntu based GNU/Linux distro I try always has some problem with the system clock. Even when I set it manually, it ends up shifting a little bit over a few weeks. I suspect this is the fault of my eleven year old laptop, but still.  
So I made this tool to automatically set the correct(ish) time and date automatically. It queries the [WorldTimeAPI](http://worldtimeapi.org/) and generates a valid `date` command you can run in your shell.

## Usage
### Flags
* `--help`, `-h`  
    Show a help screen
    
* `--command`, `-c`  
    Format the output as a date command. Example: `date --set "20021227 21:23"`
    
* `--sudo`, `-s`  
    Format the output as a date command with sudo. Example: `sudo date --set "20021227 21:23"`

* `--location`, `-l`  
    Set the location of the timezone in the call to the API. Example: `Europe`

* `--area`, `-a`  
    Set the area of the timezone in the call to the API. Example: `Rome`

### Examples
```
$ dateupdate
20210730 19:31

$ dateupdate --command
date --set "20210730 19:31"

$ dateupdate --command --sudo
sudo date --set "20210730 19:31"

$ dateupdate --command --sudo --area Europe --location London
sudo date --set "20210730 18:31"

$ # run the output of the program as a command
$ dateupdate --command --sudo --area Europe --location London | sh
ven 30 lug 2021, 18:31:00, CEST
```

## Installation
### Compiling the executable
You must have `go` in your path to compile the application, though you probably don't need the latest version. After compilation, move the newly created binary to your `/usr/bin` folder (or any other bin folder you may prefer) in order to be able to call it from the command line.
```
$ go build .
$ sudo mv dateupdate /usr/bin/
```
### Adding a cron job
You can set up a cron job to automatically correct your clock every once in a while. Since `date` requires root privileges to modify the clock, you're going to have to open a root shell, open `crontab`, and set up your cron job.
```
$ su root
# crontab -e
```

Add the following line to your crontab
```
0 0 * * * dateupdate --command --area Europe --location Rome | sh
```
This will run the program every day at midnight as root (remember to replace Europe and Rome with your timezone's area and location)