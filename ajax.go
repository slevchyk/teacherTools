package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	"regexp"

	"github.com/slevchyk/teacherTools/dbase"
)

func checkEmail(w http.ResponseWriter, r *http.Request) {

	bs, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "Internal server error. Can't check email")
		return
	}

	sBody := string(bs)
	xs := strings.Split(sBody, "|")

	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(xs[0]) {
		fmt.Fprint(w, "Invalid format")
		return
	}

	if len(xs) == 2 && xs[0] == xs[1] {
		fmt.Fprint(w, "current")
		return
	} else {
		rows, err := db.Query(dbase.SelectUserByEmail(), xs[0])
		if err != nil {
			fmt.Fprint(w, "Internal server error. Can't check email")
			return
		}

		if rows.Next() {
			fmt.Fprint(w, "false")
		} else {
			fmt.Fprint(w, "true")
		}
	}
}
