package oauth

import (
	"github.com/Saku0512/CTFLab/ctflab/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func Init() {
	config.InitDB()

	goth.UseProviders(
		github.New(config.GetGitHubAuthKey(), config.GetGitHubAuthSecret(), config.GetGitHubCallbackURL()),
		google.New(config.GetGoogleAuthKey(), config.GetGoogleAuthSecret(), config.GetGoogleCallbackURL()),
	)
}
