package libs

import (
	"bytes"
	"encoding/json"
	"flag"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	logger "github.com/accuknox/knoxAutoPolicy/src/logging"
	"github.com/accuknox/knoxAutoPolicy/src/types"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v2"
)

// =================== //
// == Configuration == //
// =================== //

func LoadConfigurationFile() {
	configFilePath := flag.String("config-path", "conf/", "conf/")
	flag.Parse()

	viper.SetConfigName(GetEnv("CONF_FILE_NAME", "conf"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(*configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		if readErr, ok := err.(viper.ConfigFileNotFoundError); ok {
			var log *zerolog.Logger = logger.GetInstance()
			log.Panic().Msgf("No config file found at %s\n", *configFilePath)
		} else {
			var log *zerolog.Logger = logger.GetInstance()
			log.Panic().Msgf("Error reading config file: %s\n", readErr)
		}
	}
}

// ================== //
// == Print Pretty == //
// ================== //

func PrintPolicyJSON(data interface{}) (string, error) {
	empty := ""
	tab := "  "

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return empty, err
	}

	return buffer.String(), nil

}

func PrintPolicyYaml(data interface{}) (string, error) {
	b, _ := yaml.Marshal(&data)
	return string(b), nil
}

// ============= //
// == Network == //
// ============= //

func getIPAddr(ifname string) string {
	if interfaces, err := net.Interfaces(); err == nil {
		for _, iface := range interfaces {
			if iface.Name == ifname {
				addrs, err := iface.Addrs()
				if err != nil {
					panic(err)
				}
				ipaddr := strings.Split(addrs[0].String(), "/")[0]
				return ipaddr
			}
		}
	}

	return "None"
}

func getExternalInterface() string {
	route := GetCommandOutput("ip", []string{"route", "get", "8.8.8.8"})
	routeData := strings.Split(strings.Split(route, "\n")[0], " ")

	for idx, word := range routeData {
		if word == "dev" {
			return routeData[idx+1]
		}
	}

	return "None"
}

// GetExternalIPAddr Function
func GetExternalIPAddr() string {
	iface := getExternalInterface()
	if iface != "None" {
		return getIPAddr(iface)
	}

	return "None"
}

// GetProtocol Function
func GetProtocol(protocol int) string {
	protocolMap := map[int]string{
		1:   "icmp",
		6:   "tcp",
		17:  "udp",
		132: "stcp",
	}

	return protocolMap[protocol]
}

// ============ //
// == Common == //
// ============ //

func DeepCopy(dst, src interface{}) {
	byt, _ := json.Marshal(src)
	json.Unmarshal(byt, dst)
}

// exists Function
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func IsK8sEnv() bool {
	if _, ok := os.LookupEnv("KUBERNETES_PORT"); ok {
		return true
	}

	k8sConfig := os.Getenv("HOME") + "./kube"
	if exist, _ := exists(k8sConfig); exist {
		return true
	}

	return false
}

func GetOSSigChannel() chan os.Signal {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	return c
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func ContainsElement(slice interface{}, element interface{}) bool {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)

		for i := 0; i < s.Len(); i++ {
			val := s.Index(i).Interface()
			if reflect.DeepEqual(val, element) {
				return true
			}
		}
	}

	return false
}

func RandSeq(n int) string {
	var lowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)

	for i := range b {
		b[i] = lowerLetters[rand.Intn(len(lowerLetters))]
	}

	return string(b)
}

// GetCommandOutput Function
func GetCommandOutput(cmd string, args []string) string {
	res := exec.Command(cmd, args...)
	out, err := res.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// ============== //
// == File I/O == //
// ============== //

func WriteKnoxPolicyToYamlFile(namespace string, policies []types.KnoxNetworkPolicy) {
	fileName := GetEnv("POLICY_DIR", "./") + "knox_policies_" + namespace + ".yaml"

	os.Remove(fileName)

	// create policy file
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	for _, policy := range policies {
		// set flow ids null
		policy.FlowIDs = nil

		b, _ := yaml.Marshal(&policy)
		f.Write(b)
		f.WriteString("---\n")
		f.Sync()
	}

	f.Close()
}

func WriteCiliumPolicyToYamlFile(namespace string, policies []types.CiliumNetworkPolicy) {
	// create policy file
	fileName := GetEnv("POLICY_DIR", "./") + "cilium_policies_" + namespace + ".yaml"

	os.Remove(fileName)

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	for _, policy := range policies {
		b, _ := yaml.Marshal(&policy)
		f.Write(b)
		f.WriteString("---\n")
		f.Sync()
	}

	f.Close()
}

func WriteKubeArmorPolicyToYamlFile(namespace string, policies []types.KubeArmorSystemPolicy) {
	// create policy file
	fileName := GetEnv("POLICY_DIR", "./") + "kubearmor_policies_" + namespace + ".yaml"

	os.Remove(fileName)

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	for _, policy := range policies {
		b, _ := yaml.Marshal(&policy)
		f.Write(b)
		f.WriteString("---\n")
		f.Sync()
	}

	f.Close()
}

// ========== //
// == Time == //
// ========== //

// Time Format
const (
	TimeForm       string = "2006-01-02T15:04:05.000000"
	TimeFormSimple string = "2006-01-02 15:04:05"
	TimeFormUTC    string = "2006-01-02T15:04:05.000000Z"
	TimeFormHuman  string = "2006-01-02 15:04:05.000000"
	TimeCilium     string = "2006-01-02T15:04:05.000000000Z"
)

func ConvertUnixTSToDateTime(ts int64) primitive.DateTime {
	t := time.Unix(ts, 0)
	dateTime := primitive.NewDateTimeFromTime(t)
	return dateTime
}

// ConvertStrToUnixTime function: str -> unix seconds for mysql
func ConvertStrToUnixTime(strTime string) int64 {
	if strTime == "now" {
		return time.Now().UTC().Unix()
	}

	t, _ := time.Parse(TimeFormSimple, strTime)
	return t.UTC().Unix()
}
