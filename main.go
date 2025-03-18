package main

import (

  "fmt"
  "flag"
  "time"
  "os"
  "bufio"
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

type Parms struct {

  cmdPtr *string 
  useridPtr *string
  passwdPtr *string
  confPtr *string
  urlPtr *string
  urlfilePtr *string 
  ptzpresetPtr *string 
  bdebugPtr *bool

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
  fmt.Println("       getnet - Gets Camera Basic Network")
  fmt.Println("       getntp - Gets Camera System NTP settings")
  fmt.Println("       ntpenable - Turn on using NTP")
  fmt.Println("       ntpdisable - Turn off using NTP")
  fmt.Println("       getptzstatus - Get Ptz Status info")
  fmt.Println("       getptzconfig - Get Ptz Config info")
  fmt.Println("       ptzgoto - Send PTZ to preset number location")


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

func getptzstatus(url string) bool{

  fmt.Printf("Getting Camera PTZ Status from [%s]\n", url)

  camera_url := fmt.Sprintf("%s/cgi-bin/ptz.cgi?action=getStatus", url)

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

func ptzgoto(url string, ptzpreset string) bool{

  fmt.Printf("Sending PTZ to a preset location [%s]\n", url)

  camera_url := fmt.Sprintf("%s/cgi-bin/ptz.cgi?action=start&channel=0&code=GotoPreset&arg1=0&arg2=%s&arg3=0&arg4=0", url, ptzpreset)

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

func getptzconfig(url string) bool{

  fmt.Printf("Getting Camera PTZ Config from [%s]\n", url)

  camera_url := fmt.Sprintf("%s/cgi-bin/configManager.cgi?action=getConfig&name=Ptz", url)

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

func getnet(url string) bool{

  camera_url := fmt.Sprintf("%s/cgi-bin/configManager.cgi?action=getConfig&name=Network", url)

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

func getntp(url string) bool{

  fmt.Printf("Getting Camera NTP Config [%s]\n", url)

  camera_url := fmt.Sprintf("%s/cgi-bin/configManager.cgi?action=getConfig&name=NTP", url)

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

func setntponoff(url string, on bool) bool{

  if(on){
    fmt.Print("Enabling NTP ");
  }else{
    fmt.Print("Disabling NTP ");
  }

  fmt.Printf("Camera NTP Config [%s]\n", url)


  camera_url := fmt.Sprintf("%s/cgi-bin/configManager.cgi?action=setConfig&NTP.Enable=%t", url, on)

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

func printparms(parms Parms){

  fmt.Println("cmdPtr = ", *parms.cmdPtr)
  fmt.Println("useridPtr = ", *parms.useridPtr)
  fmt.Println("passwdPtr = ", *parms.passwdPtr)
  fmt.Println("confPtr = ", *parms.confPtr)
  fmt.Println("urlPtr = ", *parms.urlPtr)
  fmt.Println("urlfilePtr = ", *parms.urlfilePtr)
  fmt.Println("ptzpresetPtr = ", *parms.ptzpresetPtr)
  fmt.Println("bdebugPtr = ", *parms.bdebugPtr)

}

func submain(parms Parms) bool{

  switch *parms.cmdPtr {

    case "readconf":
      fmt.Println("Reading conf file")
      readconf(*parms.confPtr)
      Dump()

    case "setconf":
      fmt.Println("Setting conf file")


      simple := simpleconffile.New(COMPILE_IN_KEY, *parms.confPtr)


      gMyconf.Encrypted = true
      gMyconf.Userid = simple.EncryptString(*parms.useridPtr)
      gMyconf.Passwd = simple.EncryptString(*parms.passwdPtr)
      gMyconf.Filename = *parms.confPtr

      simple.SaveConf(gMyconf)

    case "gettime":
      readconf(*parms.confPtr)
      gettime(*parms.urlPtr)
    
    case "settime":
      readconf(*parms.confPtr)
      settime(*parms.urlPtr)
    
    case "getinfo":
      readconf(*parms.confPtr)
      getinfo(*parms.urlPtr)

    case "getnet":
      readconf(*parms.confPtr)
      getnet(*parms.urlPtr)

    case "getntp":
      readconf(*parms.confPtr)
      getntp(*parms.urlPtr)

    case "ntpenable":
      readconf(*parms.confPtr)
      setntponoff(*parms.urlPtr, true)

    case "ntpdisable":
      readconf(*parms.confPtr)
      setntponoff(*parms.urlPtr, false)

    case "getptzstatus":
      readconf(*parms.confPtr)
      getptzstatus(*parms.urlPtr)

    case "getptzconfig":
      readconf(*parms.confPtr)
      getptzconfig(*parms.urlPtr)

    case "ptzgoto":
      readconf(*parms.confPtr)
      ptzgoto(*parms.urlPtr, *parms.ptzpresetPtr)

    default:
      help()
      os.Exit(2)

  }

  return true
}


func main(){

  var ourParms Parms

  logmsg.SetLogFile("amcrest.log")

  ourParms.cmdPtr = flag.String("cmd", "help", "Command to run")
  ourParms.useridPtr = flag.String("userid", "notset", "Amcrest userid")
  ourParms.passwdPtr = flag.String("passwd", "notset", "Amcrest password")
  ourParms.confPtr = flag.String("conffile", ".amcrest.conf", "config file name")
  ourParms.urlPtr = flag.String("url", "http://localhost", "URL (IP) of the camera")
  ourParms.urlfilePtr = flag.String("urlfile", "notset", "File of URLs (IPs) of the cameras")
  ourParms.ptzpresetPtr = flag.String("ptzpreset", "1", "Go to a preset number default is 1, which is Home")
  ourParms.bdebugPtr = flag.Bool("debug", false, "If true, do debug magic")

  flag.Parse()


  if(*ourParms.urlfilePtr == "notset"){

      url := *ourParms.urlPtr

      if(url[0:4] != "http"){
        url = "http://" + *ourParms.urlPtr
        ourParms.urlPtr = &url
      }

      fmt.Printf("%s\n", *ourParms.urlPtr)


    submain(ourParms)
  }else{
  
    fmt.Printf("Processing file [%s] of URLs\n", *ourParms.urlfilePtr)

    file, err := os.Open(*ourParms.urlfilePtr)

    if err != nil {
      fmt.Println(err)
      return
    }

    defer file.Close()
  
    // Create a scanner
    scanner := bufio.NewScanner(file)

    var url string

    // Read and print lines
    for scanner.Scan() {
      line := scanner.Text()

      if(line[0] == '#'){
        continue
      }


      if(line[0:4] == "http"){
        url = line
      }else{
        url = "http://" + line
      }

      fmt.Printf("%s\n", url)
      ourParms.urlPtr = &url
      submain(ourParms)
    }

  } //end if using an inputfile


  os.Exit(0)

}



func examplebase64(){

  s := "webcams:abc123"

  b64 := base64.StdEncoding.EncodeToString([]byte(s))

  //fmt.Println(s)
  fmt.Print(b64)

}
