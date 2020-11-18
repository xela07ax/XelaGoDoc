package fsKode

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/xela07ax/toolsXela/hulyttp"
	"io"
	"net/http"
	"os"
)

type Kod struct {
	Loger chan <- [4]string
	List map[string]bool
}


func (k *Kod)StopContainer(w http.ResponseWriter, r *http.Request) {

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
					HostPort: "4140",
				},
			},
		},
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:  "xaljer/kodexplorer",
		Tty:     false,
		Volumes: map[string]struct{}{
			"/home/droid/Projects/GitHub/XelaGoDoc/docker-api/static:/var/www/html": {},
		},
		ExposedPorts: nat.PortSet{
			"80/tcp": struct{}{},
		},
	}, hostConfig, nil, nil, "")
	if err != nil {
		panic(err)
	}
	k.List[resp.ID] = true
	hulyttp.Resp(w,r,"runContainer", resp.ID, 0, true)
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}
	if _, ok := k.List[resp.ID]; !ok {
		delete(k.List,resp.ID)
	}
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
