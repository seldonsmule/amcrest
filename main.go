package main

import (

  "fmt"
  "flag"
  "time"
  "os"
  "io/ioutil"
  "encoding/base64"
  "github.com/seldonsmule/simpleconffile"
//  "github.com/seldonsmule/restapi"
  "net/http"
  "github.com/icholy/digest"

  "github.com/seldonsmule/logmsg"

)

type Configuration struct {

  Userid string
  Passwd string
  Filename string
  Encrypted bool

}

const COMPILE_IN_KEY = "example key 1234"

var gMyconf Configuration

func Dump(){


  fmt.Println("Userid=" + gMyconf.Userid)
  fmt.Println("Passwd=" + gMyconf.Passwd)
  fmt.Println("Filename=" + gMyconf.Filename)
  fmt.Printf("Encrypted=%t\n" , gMyconf.Encrypted)

}

func help(){

  fmt.Println("CLI for Amcrest Cameras");
      
  fmt.Println("Usage amcrest -cmd [a command, see below]")
  fmt.Println()
  flag.PrintDefaults()
  fmt.Println()
  fmt.Println("cmds:")
  fmt.Println("       setconf - Setup Conf file")
  fmt.Println("             -userid Camera userid")
  fmt.Println("             -passwd Camera password")
  fmt.Println("             -conffile name of conffile (.amcrest.conf default)")
  fmt.Println()
  fmt.Println("       readconf - Display conf info")
  fmt.Println()
  fmt.Println("       gettime - Displays camera's current time")
  fmt.Println("       settime - Sets the cameras current time to our system time")
  fmt.Println("       getinfo - Gets Camera System Info")


}

func readconf(confFile string){

  simple := simpleconffile.New(COMPILE_IN_KEY, confFile)

  if(!simple.ReadConf(&gMyconf)){
    fmt.Println("Error reading conf file: ", confFile)
    os.Exit(3)
  }

  if(gMyconf.Encrypted){
    gMyconf.Userid = simple.DecryptString(gMyconf.Userid)
    gMyconf.Passwd = simple.DecryptString(gMyconf.Passwd)
  }

}

func send(url string, userid string, passwd string) (string, bool) {

  client := &http.Client{
		Transport: &digest.Transport{
			Username: userid,
			Password: passwd,
		},
	}

  res, err := client.Get(url)
    if err != nil {
      logmsg.Print(logmsg.Error, "Error", err)
      fmt.Println("Error getting to server at URL:", url)
      return "failed", false
    }

  defer res.Body.Close()

  switch res.StatusCode {

    case 200:
    case 201:

    default:
      logmsg.Print(logmsg.Error,"HTTP Response Status:", res.StatusCode, http.StatusText(res.StatusCode))
      logmsg.Print(logmsg.Error,url)

      return "failed", false

  }

  fmt.Printf("Return status code: %d\n", res.StatusCode)
  body, _ := ioutil.ReadAll(res.Body)

  //fmt.Println(string(body))

  return string(body), true

}

func gettime(url string) bool{

  fmt.Printf("Getting Camera Time from [%s]\n", url)

  camera_url := fmt.Sprintf("%s/cgi-bin/global.cgi?action=getCurrentTime", url)

  fmt.Printf("full url: %s\n", camera_url)

  rtnstr, bworked := send(camera_url, gMyconf.Userid, gMyconf.Passwd)

  if(!bworked){
    fmt.Println("Send failed")
    return false
  }else{
    fmt.Println(rtnstr)
  }

  return true
}

func settime(url string) bool{

  fmt.Printf("Setting Camera Time for [%s]\n", url)

  t := time.Now()

/*
  fmt.Println("Current Time: ", t.String())
  fmt.Println("Current Date Formated: ", t.Format("2006-01-02"))
  fmt.Println("Current Time Formated: ", t.Format("15:04:05"))
*/
  
  new_time := fmt.Sprintf("%s%c20%s", t.Format("2006-01-02"),
                                      '%',
                                     t.Format("15:04:05"))


  camera_url := fmt.Sprintf("%s/cgi-bin/global.cgi?action=setCurrentTime&time=%s" , url, new_time)
         


  //fmt.Printf("full url: %s\n", camera_url)



  rtnstr, bworked := send(camera_url, gMyconf.Userid, gMyconf.Passwd)

  if(!bworked){
    fmt.Println("Send failed")
    return false
  }else{
    fmt.Println(rtnstr)
  }


  return true
}

func getinfo(url string) bool{

  fmt.Printf("Getting Camera System Info [%s]\n", url)

  camera_url := fmt.Sprintf("%s/cgi-bin/magicBox.cgi?action=getSystemInfo", url)

  fmt.Printf("full url: %s\n", camera_url)

  rtnstr, bworked := send(camera_url, gMyconf.Userid, gMyconf.Passwd)

  if(!bworked){
    fmt.Println("Send failed")
    return false
  }else{
    fmt.Println(rtnstr)
  }

  return true
}


func main(){


  logmsg.SetLogFile("amcrest.log")

  cmdPtr := flag.String("cmd", "help", "Command to run")
  useridPtr := flag.String("userid", "notset", "Amcrest userid")
  passwdPtr := flag.String("passwd", "notset", "Amcrest password")
  confPtr := flag.String("conffile", ".amcrest.conf", "config file name")
  urlPtr := flag.String("url", "http://localhost", "URL (IP) of the camera")
  bdebugPtr := flag.Bool("debug", false, "If true, do debug magic")

  flag.Parse()

  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "useridPtr = ", *useridPtr)
  logmsg.Print(logmsg.Info, "passwdPtr = ", *passwdPtr)
  logmsg.Print(logmsg.Info, "confPtr = ", *confPtr)
  logmsg.Print(logmsg.Info, "urlPtr = ", *urlPtr)
  logmsg.Print(logmsg.Info, "bdebugPtr = ", *bdebugPtr)

  switch *cmdPtr {

    case "readconf":
      fmt.Println("Reading conf file")
      readconf(*confPtr)
      Dump()

    case "setconf":
      fmt.Println("Setting conf file")


      simple := simpleconffile.New(COMPILE_IN_KEY, *confPtr)


      gMyconf.Encrypted = true
      gMyconf.Userid = simple.EncryptString(*useridPtr)
      gMyconf.Passwd = simple.EncryptString(*passwdPtr)
      gMyconf.Filename = *confPtr

      simple.SaveConf(gMyconf)

    case "gettime":
      readconf(*confPtr)
      gettime(*urlPtr)
    
    case "settime":
      readconf(*confPtr)
      settime(*urlPtr)
    
    case "getinfo":
      readconf(*confPtr)
      getinfo(*urlPtr)

    default:
      help()
      os.Exit(2)

  }


  os.Exit(0)

}



func examplebase64(){

  s := "webcams:abc123"

  b64 := base64.StdEncoding.EncodeToString([]byte(s))

  //fmt.Println(s)
  fmt.Print(b64)

}
