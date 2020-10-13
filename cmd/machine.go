package cmd

import (
	"context"
	"fmt"
	"os"

	packet_logr "github.com/packethost/pkg/log/logr"
	"github.com/spf13/cobra"
	v1 "github.com/tinkerbell/pbnj/api/v1"
	v1Client "github.com/tinkerbell/pbnj/client"
	"google.golang.org/grpc"
)

// machineCmd represents the server command
var machineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Run PBnJ client machine actions",
	Long:  `Run PBnJ client for machine actions.`,
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		logger, zlog, err := packet_logr.NewPacketLogr(
			packet_logr.WithServiceName("github.com/tinkerbell/pbnj"),
			packet_logr.WithLogLevel(logLevel),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer zlog.Sync() // nolint
		logger.V(0).Info("machine actions")

		opts = append(opts, grpc.WithInsecure())
		conn, err := grpc.Dial("localhost:50051", opts...)
		if err != nil {
			logger.V(0).Error(err, "fail to dial server")
			os.Exit(1)
		}
		defer conn.Close()
		client := v1.NewMachineClient(conn)
		taskClient := v1.NewTaskClient(conn)

		resp, err := v1Client.MachinePower(ctx, client, taskClient, &v1.PowerRequest{
			Authn: &v1.Authn{
				Authn: &v1.Authn_DirectAuthn{
					DirectAuthn: &v1.DirectAuthn{
						Host: &v1.Host{
							Host: "10.249.51.28",
						},
						Username: "Administrator",
						Password: "YFJDJC9T",
					},
				},
			},
			Vendor: &v1.Vendor{
				Name: "HP",
			},
			Action: v1.PowerRequest_STATUS,
		})
		if err != nil {
			logger.V(0).Error(err, "error calling")
			os.Exit(1)
		}

		logger.V(0).Info("resp", "resp", []interface{}{resp})

	},
}

func init() {
	clientCmd.AddCommand(machineCmd)
}

/*
func structToSlice(s v1.StatusResponse) []interface{} {
	var keysAndValues []interface{}
	keysAndValues = append(keysAndValues, "id")
}
*/
