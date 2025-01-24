package database

import (
	"fmt"
	"os"
)

// type NewShotFromClient struct {
// 	Shot_id    int       `json:"Id"`
// 	Date_given time.Time `json:"Date_given"`
// 	Date_due   time.Time `json:"Next_due"`
// }

func inserQueryFail(err error, name string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed of %s: %v\n", name, err)
	}
	fmt.Printf("command  '%s' created successfully\n", name)
}
