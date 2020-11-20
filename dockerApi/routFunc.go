package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/xela07ax/toolsXela/tp"
	"io"
	"net/http"
	"os"
)

type Kod struct {
	Loger chan <- [4]string
	List map[string]bool
}


func (k *Kod)StopContainer(w http.ResponseWriter, r *http.Request) {
	// Над считать id контейнера
	b, err := tp.HttpReadBody(w,r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ReadContainerId:%s\n",b)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	fmt.Print("Stopping container \n", string(b), "... ")
	if err := cli.ContainerStop(ctx, string(b), nil); err != nil {
		panic(err)
	}
	fmt.Println("Success")
	tp.HttpBytes(w, r, []byte(fmt.Sprintf("Stopped container %s %s\n", string(b), "... ")))
}
func (k *Kod)RunContainer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP: "0.0.0.0",
					HostPort: "8206",
				},
			},
		},
	}

	resp, err :=
		cli.ContainerCreate(ctx, &container.Config{
			Image:  "xaljer/kodexplorer",
			Tty:     false,
			Volumes: map[string]struct{}{
				"/home/droid/Projects/GitHub/XelaGoDoc/dockerAi/chloger:/var/www/html": {},
			},
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
		}, hostConfig, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	// Распечатаем Id контейнера
	fmt.Printf("WriteContainerId:%s\n",resp.ID)
	tp.HttpBytes(w, r, []byte(resp.ID))
	//statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	//select {
	//case err := <-errCh:
	//	if err != nil {
	//		panic(err)
	//	}
	//case <-statusCh:
	//}
	//
	//out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	//if err != nil {
	//	panic(err)
	//}
	//
	//stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}