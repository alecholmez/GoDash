package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alecholmez/GoDash/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Project represents a CircleCI project associated with a user
type Project struct {
	Name     string `json:"reponame"`
	User     string `json:"username"`
	Language string `json:"language"`
	VCSType  string `json:"vcs_type"`
}

// Repo is where the current build data is from circle
type Repo struct {
	Builds []Build
}

// Build contains circle build info
type Build struct {
	Commit     string `json:"subject"`
	Status     string `json:"status"`
	User       User   `json:"user"`
	Lifecyrcle string `json:"lifecycle"`
	Branch     string `json:"branch"`
	BuildNum   int    `json:"build_num"`
	StartTime  string `json:"start_time"`
	StopTime   string `json:"stop_time"`
}

// User holds info about the user who triggered the build
type User struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar_url"`
	Name   string `json:"name"`
}

// Info to send from API call
type Info struct {
	Name     string `json:"reponame"`
	Language string `json:"language"`
	Build    Build  `json:"build_info"`
}

const (
	apiURL = "https://circleci.com/api/v1.1"
)

var (
	conf     = flag.String("config", "../etc/settings.toml", "Path to config file")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println(r.Header.Get("Origin"))
			return true
		},
	}
)

// Get API Token
var token = os.Getenv("CIRCLE_CI_AUTH_TOKEN")
var client = &http.Client{}

func main() {
	c := config.Parse(*conf)

	mux := mux.NewRouter()
	mux.HandleFunc("/dash", Dash)
	s := http.Server{
		Addr:         c.Address,
		Handler:      handlers.LoggingHandler(os.Stdout, handlers.CORS()(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Listening. . .")
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Dash is the handler that exposes the polling function
func Dash(w http.ResponseWriter, r *http.Request) {
	// Upgrade web handler to a websocket connection
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}
	defer c.Close()

	// Business logic
	if token == "" {
		err = errors.New("Missing CIRCLE_CI_AUTH_TOKEN")
		w.Write([]byte(err.Error()))
		return
	}

	for {
		log.Println("Sending current builds. . .")
		projects := getProjects(fmt.Sprintf("%s/projects?circle-token=%s", apiURL, token))
		var resp struct {
			Builds []Info `json:"builds"`
		}

		for _, project := range projects {
			url := fmt.Sprintf("%s/project/%s/%s/%s?circle-token=%s", apiURL, project.VCSType, project.User, project.Name, token)
			inf := getBuildInfo(project, url)

			resp.Builds = append(resp.Builds, inf)
		}

		err = c.WriteJSON(resp)
		if err != nil {
			log.Println("Write: ", err)
		}
		time.Sleep(time.Second * 10)
	}
}

func getBuildInfo(p Project, url string) Info {
	// Hit the circleci endpoint for associated projects
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var repo Repo
	err = json.NewDecoder(resp.Body).Decode(&repo.Builds)
	if err != nil {
		panic(err)
	}
	// Get the latest build
	inf := Info{
		Language: p.Language,
		Name:     p.Name,
		Build:    repo.Builds[0],
	}

	return inf
}

func getProjects(url string) []Project {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	// Specify json header otherwise circle won't send json
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var projects []Project
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		panic(err)
	}

	return projects
}
