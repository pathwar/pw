package pwcompose

import (
	"encoding/json"
	fmt "fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"pathwar.land/pathwar/v2/go/pkg/errcode"
)

type PrepareOpts struct {
	ChallengeDir string
	Prefix       string
	NoPush       bool
	JSON         bool
	Version      string
	Logger       *zap.Logger
}

func NewPrepareOpts() PrepareOpts {
	return PrepareOpts{
		ChallengeDir: ".",
		Prefix:       defaultDockerPrefix,
		Version:      "dev",
		NoPush:       false,
	}
}

func (opts *PrepareOpts) applyDefaults() {
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}
	if opts.ChallengeDir == "" {
		opts.ChallengeDir = "."
	}
}

func Prepare(opts PrepareOpts) (string, error) {
	opts.applyDefaults()
	opts.Logger.Debug("prepare", zap.Any("opts", opts))

	composeBin, cleanup, err := tmpComposeBin()
	if err != nil {
		return "", err
	}
	defer cleanup()

	cleanPath, err := filepath.Abs(filepath.Clean(opts.ChallengeDir))
	if err != nil {
		return "", errcode.ErrComposeInvalidPath.Wrap(err)
	}

	if opts.Prefix[len(opts.Prefix)-1:] != "/" {
		opts.Prefix += "/"
	}

	var (
		challengeName   = filepath.Base(cleanPath)
		origComposePath = path.Join(cleanPath, "docker-compose.yml")
		tmpComposePath  = path.Join(cleanPath, "docker-compose.tmp.yml")
		dabPath         = path.Join(cleanPath, challengeName+".dab")
	)

	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return "", errcode.ErrComposeDirectoryNotFound.Wrap(err)
	}

	// parse docker-compose.yml file
	composeData, err := ioutil.ReadFile(origComposePath)
	if err != nil {
		return "", errcode.ErrComposeReadConfig.Wrap(err)
	}

	// check for error in docker-compose file
	args := append(composeCliCommonArgs(origComposePath), "config", "-q")
	opts.Logger.Debug("docker-compose", zap.Strings("args", args))
	cmd := exec.Command(composeBin.Name(), args...)
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		return "", errcode.ErrComposeInvalidConfig.Wrap(err)
	}

	composeStruct := PathwarConfig{}
	err = yaml.Unmarshal(composeData, &composeStruct)
	if err != nil {
		return "", errcode.ErrComposeInvalidConfig.Wrap(err)
	}

	// check yaml and add image name if not defined
	for name, service := range composeStruct.Services {
		if service.Labels == nil {
			service.Labels = map[string]string{}
		}
		if service.Image == "" {
			if !opts.NoPush {
				challengeName = strings.ToLower(challengeName)
				service.Image = opts.Prefix + challengeName + ":" + name
				service.Labels[serviceOrigin] = "was-built"
			} else {
				service.Build = path.Join(cleanPath, service.Build)
				service.Labels[serviceOrigin] = "was-built-dev"
			}
		} else {
			service.Labels[serviceOrigin] = "was-pulled"
		}
		service.Labels[challengeNameLabel] = challengeName
		service.Labels[serviceNameLabel] = name
		service.Labels[challengeVersionLabel] = opts.Version
		composeStruct.Services[name] = service
	}

	// create tmp docker-compose file
	tmpData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		return "", errcode.ErrComposeMarshalConfig.Wrap(err)
	}
	tmpFile, err := os.Create(tmpComposePath)
	if err != nil {
		return "", errcode.ErrComposeCreateTempFile.Wrap(err)
	}
	defer func() {
		if err = os.Remove(tmpComposePath); err != nil {
			opts.Logger.Warn("rm tmp compose file", zap.Error(err), zap.String("path", tmpComposePath))
		}
	}()
	_, err = tmpFile.Write(tmpData)
	if err != nil {
		return "", errcode.ErrComposeWriteTempFile.Wrap(err)
	}
	tmpFile.Close()
	if opts.Logger.Check(zap.DebugLevel, "") != nil {
		fmt.Fprintln(os.Stderr, string(tmpData))
	}

	if !opts.NoPush {
		// build and push images to dockerhub (don't forget to setup your credentials just type : "docker login" in bash)
		args = append(composeCliCommonArgs(tmpComposePath), "build")
		opts.Logger.Debug("docker-compose", zap.Strings("args", args))
		cmd = exec.Command(composeBin.Name(), args...)
		cmd.Dir = cleanPath
		if opts.Logger.Check(zap.DebugLevel, "") != nil {
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
		if err != nil {
			return "", errcode.ErrComposeBuild.Wrap(err)
		}

		args = append(composeCliCommonArgs(tmpComposePath), "bundle", "--push-images")
		opts.Logger.Debug("docker-compose", zap.Strings("args", args))
		cmd = exec.Command(composeBin.Name(), args...)
		if err != nil {
			return "", errcode.ErrComposeBundle.Wrap(err)
		}
		cmd.Dir = cleanPath
		if opts.Logger.Check(zap.DebugLevel, "") != nil {
			cmd.Stdout = os.Stderr
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
		if err != nil {
			return "", errcode.ErrComposeBundle.Wrap(err)
		}
		defer func() {
			if err = os.Remove(dabPath); err != nil {
				opts.Logger.Warn("rm dab file", zap.Error(err))
			}
		}()

		// parse json from .dab file
		composeDabfileJSON := dabfile{}
		composeDabfile, err := ioutil.ReadFile(dabPath)
		if err != nil {
			return "", errcode.ErrComposeReadDab.Wrap(err)
		}
		if err = json.Unmarshal(composeDabfile, &composeDabfileJSON); err != nil {
			return "", errcode.ErrComposeParseDab.Wrap(err)
		}

		// replace images from original docker-compose file with the one pushed to dockerhub
		for name, service := range composeStruct.Services {
			service.Image = composeDabfileJSON.Services[name].Image
			service.Build = "" // ensure service only has an `image:` without a `build:`
			composeStruct.Services[name] = service
		}
	}

	if opts.JSON {
		out, err := json.MarshalIndent(&composeStruct, "", "  ")
		if err != nil {
			return "", errcode.ErrComposeMarshalConfig.Wrap(err)
		}
		return string(out), nil
	}

	// print yaml
	finalData, err := yaml.Marshal(&composeStruct)
	if err != nil {
		return "", errcode.ErrComposeMarshalConfig.Wrap(err)
	}

	return string(finalData), nil
}
