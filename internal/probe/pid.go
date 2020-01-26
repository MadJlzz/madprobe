package probe

import (
	"bufio"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type Pid struct {
	Probe          `yaml:",inline"`
	Hostname       string `yaml:"hostname"`
	Port           string `yaml:"port"`
	ServiceAccount string `yaml:"service-account"`
	Pid            int    `yaml:"pid"`
}

func (p *Pid) Check() {
	if p.Hostname == "localhost" {
		p.launchLocalCheck()
	} else {
		p.launchRemoteCheck()
	}
}

func (p *Pid) launchLocalCheck() {
	for {
		pid := strconv.Itoa(p.Pid)
		cmd := exec.Command("ps", "-p", pid)
		if err := cmd.Run(); err != nil {
			log.Printf("<<PID PROBE>> Process with PID [%d] not found. got '%s'\n", p.Pid, err)
			p.app.UpdateTemplateData(p.Name, StatusDown)
		} else {
			log.Printf("<<PID PROBE>> Process with PID [%d] is currently running.", p.Pid)
			p.app.UpdateTemplateData(p.Name, StatusUp)
		}
		time.Sleep(time.Duration(p.Delay) * time.Second)
	}
}

func (p *Pid) launchRemoteCheck() {
	config := initSshClient(p.Hostname, p.ServiceAccount)
	client, err := ssh.Dial("tcp", p.Hostname+":"+p.Port, config)
	if err != nil {
		log.Fatalf("error while dialing ssh server. got %s\n", err)
	}
	pid := strconv.Itoa(p.Pid)
	for {
		session, err := client.NewSession()
		if err != nil {
			log.Fatalf("error while creating new ssh session. got %s\n", err)
		}
		if err := session.Run("ps -p " + pid); err != nil {
			log.Printf("<<PID PROBE>> Process with PID [%d] not found. got '%s'\n", p.Pid, err)
			p.app.UpdateTemplateData(p.Name, StatusDown)

		} else {
			log.Printf("<<PID PROBE>> Process with PID [%d] is currently running.", p.Pid)
			p.app.UpdateTemplateData(p.Name, StatusUp)
		}
		session.Close()
		time.Sleep(time.Duration(p.Delay) * time.Second)
	}
}

func initSshClient(hostname, username string) *ssh.ClientConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error while trying to find home directory. got %s\n", err)
	}
	hostKey := hostKeyFile(hostname, path.Join(home, ".ssh", "known_hosts"))
	privateKey := privateKeyFile(path.Join(home, ".ssh", "id_rsa"))
	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			privateKey,
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
}

func hostKeyFile(hostname, knownHostsPath string) ssh.PublicKey {
	file, err := os.Open(knownHostsPath)
	if err != nil {
		log.Fatalf("error while trying to read known hosts file. got %s\n", err)
	}
	defer file.Close()

	var hostKey ssh.PublicKey
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], hostname) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}
	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", hostname)
	}
	return hostKey
}

func privateKeyFile(privateKeyPath string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("error while trying to read private key file. got %s\n", err)
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatalf("error while parsing private key. got %s\n", err)
	}
	return ssh.PublicKeys(key)
}
