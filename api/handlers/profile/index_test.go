package profile

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ILOG struct{}

func (ILOG) Println(params ...interface{}) {

}

func Test_indexHandler(t *testing.T) {
	tests := []struct {
		name string
		want []profile
	}{
		{
			name: "find all profile",
			want: []profile{
				{
					Username:  "",
					Firstname: "Amaury",
					Lastname:  "Brisou",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/countries", nil)
			if err != nil {
				t.Fatal(err)
			}

			mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
			mocket.Catcher.Logging = true

			db, err := gorm.Open(mocket.DriverName, "connection_string")
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mocket.Catcher.Attach([]*mocket.FakeResponse{
				{
					Pattern: `SELECT * FROM "profiles"  `, // the same as .WithQuery()
					Response: []map[string]interface{}{
						{
							"username":  "",
							"firstname": "Amaury",
							"lastname":  "Brisou",
						},
					}, // the same as .WithReply
					Once: false, // To not use it twice if true
				},
			})

			l := ILOG{}

			ch := UserHandler{
				Store: &UserStore{Db: db, Log: l},
				Log:   l,
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ch.IndexHandler)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			// Check the response body is what we expect.
			var countries []profile
			if err := json.Unmarshal(rr.Body.Bytes(), &countries); err != nil {
				t.Errorf("cannot decode body : %s", err)
			}

			if !reflect.DeepEqual(countries, tt.want) {
				t.Errorf("indexHandler() = %v, want %v", countries, tt.want)
			}

		})
	}
}
