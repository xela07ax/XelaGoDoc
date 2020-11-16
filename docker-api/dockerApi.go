package main

import (
	"fmt"
	"bytes"
	"strings"
	"regexp"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	var endpoint string = "unix:///var/run/docker.sock"
	var client *docker.Client
	var err error

	client, err = docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	listImages(client)
	listContainers(client)
	listServices(client)
	listTasks(client)
	listNodes(client)
	getServiceLogs(client)
	getContainerLogs(client)
}


func listImages(client * docker.Client) {
	opts := docker.ListImagesOptions{All: false}

	images, err := client.ListImages(opts)
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		fmt.Println("Image  ID      : ", image.ID)
		fmt.Println("Imager RepoTags: ", image.RepoTags)
	}
}

func listContainers(client * docker.Client) {
	opts := docker.ListContainersOptions{}

	containers, err := client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		fmt.Println("Container ID   : ", container.ID)
		fmt.Println("Container Names: ", container.Names)
	}
}

func listServices(client * docker.Client) {
	//opts := docker.ListServicesOptions{Filters: map[string][]string {"name":{"<servicename>"}}}
	opts := docker.ListServicesOptions{}

	services, err := client.ListServices(opts)
	if err != nil {
		panic(err)
	}
	for _, service := range services {
		fmt.Println("Service ID      : ", service.ID)
		fmt.Println("Service Name    : ", service.Spec.Name)
		fmt.Println("Service Replicas: ", *(service.Spec.Mode.Replicated.Replicas))
	}
}

func listTasks(client * docker.Client) {
	opts := docker.ListTasksOptions{Filters: map[string][]string {"service":{"<servicename>"}}}

	tasks, err := client.ListTasks(opts)
	if err != nil {
		panic(err)
	}
	for _, task := range tasks {
		fmt.Println("Task ID         : ", task.ID)
		fmt.Println("Task Name       : ", task.Name)
		fmt.Println("Task ServiceID  : ", task.ServiceID)
		fmt.Println("Task Slot       : ", task.Slot)
		fmt.Println("Task NodeID     : ", task.NodeID)
		fmt.Println("Task Message    : ", task.Status.Message)
		fmt.Println("Task Error      : ", task.Status.Err)
		fmt.Println("Task ContainerID: ", task.Status.ContainerStatus.ContainerID)
		fmt.Println("Task PID        : ", task.Status.ContainerStatus.PID)
		fmt.Println("Task ExitCode   : ", task.Status.ContainerStatus.ExitCode)
	}
}

func listNodes(client * docker.Client) {
	opts := docker.ListNodesOptions{}

	nodes, err := client.ListNodes(opts)
	if err != nil {
		panic(err)
	}
	for _, node := range nodes {
		fmt.Println("Node ID           : ", node.ID)
		fmt.Println("Node Hostname     : ", node.Description.Hostname)
		fmt.Println("Node State        : ", node.Status.State)
		fmt.Println("Node Message      : ", node.Status.Message)
		fmt.Println("Node Addr         : ", node.Status.Addr)
		fmt.Println("Node ManagerStatus: ", node.ManagerStatus)
	}
}

func getServiceLogs(client * docker.Client) {
	var buf bytes.Buffer
	opts := docker.LogsServiceOptions{
		Service:      "<servicename>",
		OutputStream: &buf,
		ErrorStream:  &buf,
		Follow:       false,
		Stdout:       true,
		Stderr:       true,
		Timestamps:   true,
		Tail:         "10",      // get the last 10 lines log
	}

	err := client.GetServiceLogs(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}


func getContainerLogs(client * docker.Client) {
	var buf bytes.Buffer
	opts := docker.LogsOptions {
		Container:      "<containername>",
		OutputStream: &buf,
		ErrorStream:  &buf,
		Follow:       false,
		Stdout:       true,
		Stderr:       true,
		//      Timestamps:   true,
		Tail:         "10",
	}

	err := client.Logs(opts)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if matched, _ := regexp.MatchString(`: ERROR  : `, line); matched {
			fmt.Println(line)
		}
	}
}