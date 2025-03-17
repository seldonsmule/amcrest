# amcrest
Example of how to use the Amcrest Camera CGI API

```
CLI for Amcrest Cameras
Usage amcrest -cmd [a command, see below]

  -cmd string
    	Command to run (default "help")
  -conffile string
    	config file name (default ".amcrest.conf")
  -debug
    	If true, do debug magic
  -passwd string
    	Amcrest password (default "notset")
  -url string
    	URL (IP) of the camera (default "http://localhost")
  -userid string
    	Amcrest userid (default "notset")

cmds:
       setconf - Setup Conf file
             -userid Camera userid
             -passwd Camera password
             -conffile name of conffile (.amcrest.conf default)

       readconf - Display conf info

       gettime - Displays camera's current time
       settime - Sets the cameras current time to our system time
       getinfo - Gets Camera System Info
```
