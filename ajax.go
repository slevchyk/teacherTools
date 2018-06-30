package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	"github.com/slevchyk/teacherTools/dbase"
)

func checkUsername(w http.ResponseWriter, r *http.Request) {

	bs, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	sBody := string(bs)
	xs := strings.Split(sBody, "|")

	if len(xs) == 2 && xs[0] == xs[1] {
		fmt.Fprint(w, "true")
		return
	} else {
		rows, err := db.Query(dbase.GetQuery(dbase.SUserByEmail), xs[0])
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		fmt.Println(xs[0])
		if rows.Next() {
			fmt.Fprint(w, "false")
		} else {
			fmt.Fprint(w, "true")
		}
	}
}
