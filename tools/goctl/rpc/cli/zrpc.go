package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	errInvalidGrpcOutput = errors.New("ZRPC: missing --go-grpc_out or --gogofaster_out")
	errInvalidGoOutput   = errors.New("ZRPC: missing --go_out or --gogofaster_out")
	errInvalidZrpcOutput = errors.New("ZRPC: missing zrpc output, please use --zrpc_out to specify the output")
)

// ZRPC generates grpc code directly by protoc and generates
// zrpc code by goctl.
func ZRPC(_ *cobra.Command, args []string) error {
	protocArgs := wrapProtocCmd("protoc", args)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	source := args[0]
	grpcOutList := VarStringSliceGoGRPCOut
	goOutList := VarStringSliceGoOut
	gogoFasterOutList := VarStringSliceGoGoFasterOut
	zrpcOut := VarStringZRPCOut
	style := VarStringStyle
	home := VarStringHome
	remote := VarStringRemote
	branch := VarStringBranch
	verbose := VarBoolVerbose
	if len(grpcOutList) == 0 && len(gogoFasterOutList) == 0 {
		return errInvalidGrpcOutput
	}
	if len(goOutList) == 0 && len(gogoFasterOutList) == 0 {
		return errInvalidGoOutput
	}
	var (
		goOut         string
		grpcOut       string
		gogoFasterOut string
	)
	if len(goOutList) != 0 {
		goOut = goOutList[len(goOutList)-1]
	}
	if len(grpcOutList) != 0 {
		grpcOut = grpcOutList[len(grpcOutList)-1]
	}
	if len(gogoFasterOutList) != 0 {
		gogoFasterOut = gogoFasterOutList[len(gogoFasterOutList)-1]
	}
	if len(goOut) == 0 && len(gogoFasterOut) == 0 {
		return errInvalidGrpcOutput
	}
	if len(zrpcOut) == 0 {
		return errInvalidZrpcOutput
	}
	goOutAbs, err := filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOutAbs, err := filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	gogoFasterAbs, err := filepath.Abs(gogoFasterOut)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(goOutAbs)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(grpcOutAbs)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(gogoFasterAbs)
	if err != nil {
		return err
	}
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if !filepath.IsAbs(zrpcOut) {
		zrpcOut = filepath.Join(pwd, zrpcOut)
	}

	isGooglePlugin := len(grpcOut) > 0
	goOut, err = filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOut, err = filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	gogoFasterOut, err = filepath.Abs(gogoFasterOut)
	if err != nil {
		return err
	}
	zrpcOut, err = filepath.Abs(zrpcOut)
	if err != nil {
		return err
	}

	var ctx generator.ZRpcContext
	ctx.Multiple = VarBoolMultiple
	ctx.Src = source
	ctx.GoOutput = goOut
	ctx.GrpcOutput = grpcOut
	ctx.GoGoFasterOutput = gogoFasterOut
	ctx.IsGooglePlugin = isGooglePlugin
	ctx.Output = zrpcOut
	ctx.ProtocCmd = strings.Join(protocArgs, " ")
	ctx.IsGenClient = VarBoolClient
	g := generator.NewGenerator(style, verbose)
	return g.Generate(&ctx)
}

func wrapProtocCmd(name string, args []string) []string {
	ret := append([]string{name}, args...)
	for _, protoPath := range VarStringSliceProtoPath {
		ret = append(ret, "--proto_path", protoPath)
	}
	for _, goOpt := range VarStringSliceGoOpt {
		ret = append(ret, "--go_opt", goOpt)
	}
	for _, goGRPCOpt := range VarStringSliceGoGRPCOpt {
		ret = append(ret, "--go-grpc_opt", goGRPCOpt)
	}
	for _, gogoFasterOpt := range VarStringSliceGoGoFasterOpt {
		ret = append(ret, "--gogofaster_opt", gogoFasterOpt)
	}
	for _, goOut := range VarStringSliceGoOut {
		ret = append(ret, "--go_out", goOut)
	}
	for _, goGRPCOut := range VarStringSliceGoGRPCOut {
		ret = append(ret, "--go-grpc_out", goGRPCOut)
	}
	for _, gogoFasterOut := range VarStringSliceGoGoFasterOut {
		ret = append(ret, "--gogofaster_out", gogoFasterOut)
	}
	for _, plugin := range VarStringSlicePlugin {
		ret = append(ret, "--plugin="+plugin)
	}
	return ret
}
