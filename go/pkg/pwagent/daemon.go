package pwagent

import (
	"context"
	"encoding/json"
	fmt "fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/gogo/protobuf/jsonpb"
	"go.uber.org/zap"
	"pathwar.land/go/pkg/errcode"
	"pathwar.land/go/pkg/pwapi"
	"pathwar.land/go/pkg/pwcompose"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwinit"
)

func Daemon(ctx context.Context, clean bool, runOnce bool, loopDelay time.Duration, cli *client.Client, apiClient *http.Client, httpAPIAddr string, agentName string, logger *zap.Logger) error {
	// FIXME: call API register in gRPC
	// ret, err := api.AgentRegister(ctx, &pwapi.AgentRegister_Input{Name: "dev", Hostname: "localhost", OS: "lorem ipsum", Arch: "x86_64", Version: "dev", Tags: []string{"dev"}})

	// cleanup
	if clean {
		err := pwcompose.Down(ctx, []string{}, true, true, true, cli, logger)
		if err != nil {
			return errcode.ErrCleanPathwarInstances.Wrap(err)
		}
	}

	for {
		instances, err := fetchAPIInstances(ctx, apiClient, httpAPIAddr, agentName, logger)
		if err != nil {
			logger.Error("fetch instances", zap.Error(err))

		} else {
			if err := run(ctx, instances, cli, logger); err != nil {
				logger.Error("pwdaemon", zap.Error(err))
			}
		}

		if runOnce {
			break
		}

		time.Sleep(loopDelay)
	}

	// FIXME: agent update state for each updated instances
	return nil
}

func fetchAPIInstances(ctx context.Context, apiClient *http.Client, httpAPIAddr string, agentName string, logger *zap.Logger) (*pwapi.AgentListInstances_Output, error) {
	var instances pwapi.AgentListInstances_Output

	resp, err := apiClient.Get(httpAPIAddr + "/agent/list-instances?agent_name=" + agentName)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errcode.TODO.Wrap(err)
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("received API error", zap.String("body", string(body)), zap.Int("code", resp.StatusCode))
		return nil, errcode.TODO.Wrap(fmt.Errorf("received API error"))
	}
	if err := jsonpb.UnmarshalString(string(body), &instances); err != nil {
		return nil, errcode.TODO.Wrap(err)
	}

	return &instances, nil
}

func run(ctx context.Context, apiInstances *pwapi.AgentListInstances_Output, cli *client.Client, logger *zap.Logger) error {
	// fetch local info from docker daemon
	containersInfo, err := pwcompose.GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	agentOpts := AgentOpts{
		DomainSuffix:      "local",
		HostIP:            "0.0.0.0",
		HostPort:          "8000",
		ModeratorPassword: "",
		Salt:              "1337supmyman1337",
		AllowedUsers:      map[string][]int64{},
		ForceRecreate:     false,
		NginxDockerImage:  "docker.io/library/nginx:stable-alpine",
	}

	// compute instances that needs to upped / redumped
	for _, apiInstance := range apiInstances.GetInstances() {
		found := false
		needRedump := false
		for _, flavor := range containersInfo.RunningFlavors {
			apiInstanceFlavor := apiInstance.GetFlavor()
			apiInstanceFlavorChallenge := apiInstanceFlavor.GetChallenge()
			if apiInstanceFlavor != nil && apiInstanceFlavorChallenge != nil {
				if flavor.InstanceKey == strconv.FormatInt(apiInstance.GetID(), 10) {
					found = true
					if apiInstance.GetStatus() == pwdb.ChallengeInstance_NeedRedump {
						needRedump = true
					}
				}
			}
		}
		if !found || needRedump {
			// parse pwinit config
			var configData pwinit.InitConfig
			err = json.Unmarshal(apiInstance.GetInstanceConfig(), &configData)
			if err != nil {
				return errcode.ErrParseInitConfig.Wrap(err)
			}

			err = pwcompose.Up(ctx, apiInstance.GetFlavor().GetComposeBundle(), strconv.FormatInt(apiInstance.GetID(), 10), true, &configData, cli, logger)
			if err != nil {
				return errcode.ErrUpPathwarInstance.Wrap(err)
			}
		}
	}

	// update pathwar infos
	containersInfo, err = pwcompose.GetContainersInfo(ctx, cli)
	if err != nil {
		return errcode.ErrComposeGetContainersInfo.Wrap(err)
	}

	// update nginx configuration
	for _, apiInstance := range apiInstances.GetInstances() {
		if apiInstanceFlavor := apiInstance.GetFlavor(); apiInstanceFlavor != nil {
			if seasonChallenges := apiInstanceFlavor.GetSeasonChallenges(); seasonChallenges != nil {
				for _, seasonChallenge := range seasonChallenges {
					if subscriptions := seasonChallenge.GetActiveSubscriptions(); subscriptions != nil {
						for _, subscription := range subscriptions {
							if team := subscription.GetTeam(); team != nil {
								if members := team.GetMembers(); members != nil {
									for _, member := range members {
										for _, flavor := range containersInfo.RunningFlavors {
											if flavor.InstanceKey == strconv.FormatInt(apiInstance.GetID(), 10) {
												for _, instance := range flavor.Instances {
													for _, port := range instance.Ports {
														if port.PublicPort != 0 {
															// configure nginx
															// generate a hash per user for challenge dns prefix, based on their userIDs
															instanceName := instance.Names[0][1:]
															_, entryFound := agentOpts.AllowedUsers[instanceName]
															if !entryFound {
																agentOpts.AllowedUsers[instanceName] = []int64{member.GetID()}
															} else {
																allowedUsersSlice := agentOpts.AllowedUsers[instanceName]
																allowedUsersSlice = append(allowedUsersSlice, member.GetID())
																agentOpts.AllowedUsers[instanceName] = allowedUsersSlice
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	err = Nginx(ctx, agentOpts, cli, logger)
	if err != nil {
		return errcode.ErrUpdateNginx.Wrap(err)
	}

	return nil
}