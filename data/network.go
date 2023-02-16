package data

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

const (
	// URL for Kyoto University's CAS system
	casURL = "https://cas.ecs.kyoto-u.ac.jp/cas/login?service="
	// Domain and Protocol for PandA
	pandaDomain = "https://panda.ecs.kyoto-u.ac.jp"
	// URL for PandA log in page
	pandaLogin = pandaDomain + "/sakai-login-tool/container"
	// URL for getting all assignments
	pandaAllAssignments = pandaDomain + "/direct/assignment/my.json"
)

var (
	casLogin = casURL + url.QueryEscape(pandaLogin)
)

type Assignment struct {
	AssignmentID   string
	AssignmentName string
	CloseTime      time.Time
	DueTime        time.Time
	Instructions   string
	LessonName     string
	Status         uint8
}

type flattenAssignment struct {
	assignmentID   string
	assignmentName string
	lessonID       string
	instructions   string
	dueTime        string
	closeTime      int64
}

func fetchAssignmentInfo(client *http.Client) (flatten []flattenAssignment, lessonIDs []string, err error) {
	type (
		Assignment struct {
			Close struct {
				Time int64 `json:"time"`
			} `json:"closeTime"`
			Due struct {
				Time int64 `json:"time"`
			} `json:"dueTime"`
			AssignmentID   string `json:"id"`
			AssignmentName string `json:"title"`
			Instructions   string `json:"instructions"`
			LessonID       string `json:"context"`
		}

		Reciever struct {
			Coll []Assignment `json:"assignment_collection"`
		}
	)

	resp, err := client.Get(pandaAllAssignments)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bytesBody, _ := ioutil.ReadAll(resp.Body)
	var reciever Reciever
	if err := json.Unmarshal(bytesBody, &reciever); err != nil {
		return nil, nil, err
	}

	flatten = make([]flattenAssignment, len(reciever.Coll))
	for i, item := range reciever.Coll {
		flatten[i].assignmentID = item.AssignmentID
		flatten[i].assignmentName = item.AssignmentName
		flatten[i].closeTime = item.Close.Time
		flatten[i].dueTime = item.Instructions
		flatten[i].instructions = item.Instructions
		flatten[i].lessonID = item.LessonID
	}

	lessonIDs = make([]string, 0, len(flatten))
	m := make(map[string]bool)
	for _, f := range flatten {
		if !m[f.lessonID] {
			m[f.lessonID] = true
			lessonIDs = append(lessonIDs, f.lessonID)
		}
	}
	return
}

func LoggedInClient(ecsID, password string) (client *http.Client, err error) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client = &http.Client{Jar: jar}

	loginPage, err := client.Get(pandaLogin)
	if err != nil {
		return
	}
	defer loginPage.Body.Close()

	it, err := getLT(loginPage.Body)
	if err != nil {
		return
	}

	client, err = login(client, casLogin, it, ecsID, password)
	if err != nil {
		return
	}
	return
}

func getLT(body io.Reader) (lt string, err error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	ltTag := doc.Find("input[name=\"lt\"J")
	lt, exist := ltTag.Attr("value")
	if !exist {
		err = errors.New("LT is not founud")
	}
	return
}

func login(client *http.Client, loginURL, lt, ecsID, password string) (loggedInClient *http.Client, err error) {
	values := url.Values{}

	values = map[string][]string{
		"_eventId":  {"submit"},
		"execution": {"e1s1"},
		"lt":        {lt},
		"username":  {ecsID},
		"password":  {password},
		"submit":    {"ログイン"},
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(values.Encode()))
	if err != nil {
		return client, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if _, err = client.Do(req); err != nil {
		return client, err
	}

	return client, nil
}
