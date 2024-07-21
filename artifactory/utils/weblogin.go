package utils

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/pkg/browser"
)

func DoWebLogin(serverDetails *config.ServerDetails) (token auth.CommonTokenParams, err error) {
	if err = sendUnauthenticatedPing(serverDetails); err != nil {
		return
	}

	uuidToken, err := uuid.NewRandom()
	if errorutils.CheckError(err) != nil {
		return
	}
	uuidStr := uuidToken.String()
	accessManager, err := CreateAccessServiceManager(serverDetails, false)
	if err != nil {
		return
	}
	if err = accessManager.SendLoginAuthenticationRequest(uuidStr); err != nil {
		err = errors.Join(err,
			errorutils.CheckErrorf("Oops! it looks like the web login option is unavailable right now. This could be because your JFrog Artifactory version is less than 7.64.0. \n"+
				"Don't worry! You can use the \"jf c add\" command to authenticate with the JFrog Platform using other methods"))
		return
	}
	log.Info("After logging in via your web browser, please enter the code if prompted: " + coreutils.PrintBoldTitle(uuidStr[len(uuidStr)-4:]))
	if err = browser.OpenURL(clientUtils.AddTrailingSlashIfNeeded(serverDetails.Url) + "ui/login?jfClientSession=" + uuidStr + "&jfClientName=JFrog-CLI&jfClientCode=1"); err != nil {
		return
	}
	time.Sleep(1 * time.Second)
	log.Debug("Attempting to get the authentication token...")
	token, err = accessManager.GetLoginAuthenticationToken(uuidStr)
	if err != nil {
		return
	}
	if token.AccessToken == "" {
		return token, errorutils.CheckErrorf("failed getting authentication token after web log")
	}
	log.Info("You're now logged in!")
	return
}

func sendUnauthenticatedPing(serverDetails *config.ServerDetails) error {
	artifactoryManager, err := CreateServiceManager(serverDetails, 3, 0, false)
	if err != nil {
		return err
	}
	_, err = artifactoryManager.Ping()
	return err
}
