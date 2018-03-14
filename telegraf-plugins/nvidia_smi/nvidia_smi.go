package main
// Original https://dev.sigpipe.me/dashie/telegraf-plugins/src/branch/master/nvidia_smi/nvidia_smi.go
import (
  "flag"
  "fmt"
  "os"
  "os/exec"
  "strings"
)

func getResult(bin string, metric string, verbose bool, gpuId int) string {
  query := fmt.Sprintf("--query-gpu=%s", metric)
  gpu := fmt.Sprintf("--id=%d", gpuId)
  opts := []string{"--format=noheader,nounits,csv", query, gpu}

  if verbose {
    fmt.Print("Going to run ")
    fmt.Print(bin)
    fmt.Println(" with ")
    fmt.Println(opts)
  }
  ret, err := exec.Command(bin, opts...).CombinedOutput()
  if err != nil {
    fmt.Fprintf(os.Stderr, "%s: %s", err, ret)
    return ""
  }
  return string(ret)
}

func main() {
  binPath := flag.String("bin", "C:\\Program Files\\NVIDIA Corporation\\NVSMI\\nvidia-smi.exe", "nvidia-smi full path")
  verbose := flag.Bool("verbose", false, "display some things")
  gpuId := flag.Int("gpu", 0, "select GPU to query")

  flag.Parse()

  if _, err := os.Stat(*binPath); os.IsNotExist(err) {
    fmt.Fprintf(os.Stderr, "Bin path does not exists: %s", *binPath)
    return // exit
  }
  // nvidia-smi --format=csv --query-gpu=power.draw,utilization.gpu,fan.speed,temperature.gpu
  metrics := "fan.speed,memory.total,memory.used,memory.free,pstate,temperature.gpu,name,uuid,compute_mode,power.draw,utilization.gpu"
  results := getResult(*binPath, metrics, *verbose, *gpuId)

  if results == "" {
    return // exit
  }

  splitResults := strings.Split(results, ",")

  fmt.Printf("nvidiasmi,uuid=%s ", strings.TrimSpace(splitResults[7])) // it should be available ... if no, you have some problems

  fmt.Printf("gpu_name=\"%s\",", strings.TrimSpace(splitResults[6]))
  fmt.Printf("gpu_compute_mode=\"%s\",", strings.TrimSpace(splitResults[8]))

  fmt.Printf("fan_speed=%s,", strings.TrimSpace(splitResults[0])) // it's a % 0-100

  fmt.Printf("memory_total=%s,", strings.TrimSpace(splitResults[1])) // they
  fmt.Printf("memory_used=%s,", strings.TrimSpace(splitResults[2]))  // are
  fmt.Printf("memory_free=%s,", strings.TrimSpace(splitResults[3]))  // MiB

  //fmt.Printf("power_draw=%s,", strings.TrimSpace(splitResults[9]))    // W
  fmt.Printf("power_draw=\"%s\",", strings.TrimSpace(strings.Replace(splitResults[9], "W", "", -1)))     // W
  fmt.Printf("utilization=%s,", strings.TrimSpace(splitResults[10]))  // % 0-100

  fmt.Printf("pstate=%s,", strings.TrimSpace(strings.Replace(splitResults[4], "P", "", -1))) // strip the P
  fmt.Printf("temperature=%s\n", strings.TrimSpace(splitResults[5])) // in degrees Celcius

}