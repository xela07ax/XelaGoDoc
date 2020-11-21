package main

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/xela07ax/toolsXela/tp"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (k *Kod)RunContainerTmpfs(w http.ResponseWriter, r *http.Request) {
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
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	dir := filepath.Join(path,"rw")
	fmt.Println(dir)

	destPath := "/home/droid/Projects/GitHub/XelaGoDoc/dockerApi/fs"

	hostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   "volume",
				Source: "hello",
				Target: destPath,
				VolumeOptions: &mount.VolumeOptions{
					DriverConfig: &mount.Driver{
						Name:    "local",
						Options: map[string]string{"o": "size=100", "type": "tmpfs", "device": "tmpfs"},
					},
				},
			},
		},
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
				fmt.Sprintf("%s:/var/www/html",dir): {},
			},
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
		}, &hostConfig, nil, "")
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
