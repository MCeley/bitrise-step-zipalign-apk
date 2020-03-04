package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-android/sdk"
	"github.com/bitrise-tools/go-steputils/tools"
)

type configs struct {
	UnalignedAPKPath string `env:"bitrise_unaligned_apk_path,required"`
}

func main() {

	var cfg configs
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}

	androidHome := os.Getenv("ANDROID_HOME")
	androidSDK, err := sdk.New(androidHome)
	if err != nil {
		failf("Failed to create SDK model, error: %s", err)
	}

	zipalign, err := androidSDK.LatestBuildToolPath("zipalign")
	if err != nil {
		failf("Failed to find zipalign path, error: %s", err)
	}

	unalignedPath := strings.TrimSpace(cfg.UnalignedAPKPath)
	alignedPath := getAlignedApkName(unalignedPath)

	if err := zipalignApkArtifact(zipalign, unalignedPath, alignedPath); err != nil {
		failf("Failed to zipalign Build Artifact, error: %s", err)
	}
	fmt.Println()

	if err := tools.ExportEnvironmentWithEnvman("BITRISE_ALIGNED_APK_PATH", alignedPath); err != nil {
		failf("Failed to export APK (%s) error: %s", alignedPath, err)
	} else {
		log.Donef("The aligned APK path is now available in the Environment Variable: BITRISE_ALIGNED_APK_PATH (value: %s)", alignedPath)
	}

	os.Exit(0)
}

func getAlignedApkName(unalignedApkPath string) string {
	apkDirectory := filepath.Dir(unalignedApkPath)
	unalignedApkName := filepath.Base(unalignedApkPath)
	apkExtension := filepath.Ext(unalignedApkName)
	alignedApkName := strings.TrimSuffix(unalignedApkName, apkExtension)
	alignedApkName = fmt.Sprintf("%s-aligned%s", alignedApkName, apkExtension)
	return filepath.Join(apkDirectory, alignedApkName)
}

func zipalignApkArtifact(zipalign, inPath, outPath string) error {
	cmdSlice := []string{zipalign, "-f", "-p", "4", inPath, outPath}

	prinatableCmd := command.PrintableCommandArgs(false, cmdSlice)
	log.Printf("=> %s", prinatableCmd)

	_, err := executeForOutput(cmdSlice)
	return err
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v)
	os.Exit(1)
}

func executeForOutput(cmdSlice []string) (string, error) {
	cmd, err := command.NewFromSlice(cmdSlice)
	if err != nil {
		return "", fmt.Errorf("Failed to create command, error: %s", err)
	}

	var outputBuf bytes.Buffer
	writer := io.MultiWriter(&outputBuf)
	cmd.SetStderr(writer)
	cmd.SetStdout(writer)

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%s\n%s", outputBuf.String(), err)
	}

	return outputBuf.String(), err
}
