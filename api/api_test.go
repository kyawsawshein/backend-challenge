package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"bytes"
	"testing"

	"backend-challenge/config"
	"backend-challenge/db"
	"backend-challenge/user"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.T) {
	config.LoadConfig()
	connStr := "mongodb://root:admin@localhost:27017"
	clientOpts := options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	db.Connection = client.Database("mdb")
	db.GetCollection("user")
}

func login_test(email string, pwd string) (*http.Request, *fiber.App) {
	// config.LoadConfig()
	webApi := GetWebApi(context.Background())
	payloadStr := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, pwd)
	payload := []byte(payloadStr)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	app := webApi.GetApp()
	user.SetUpRoutes(app)
	return req, app
}

func Test_listUsers(t *testing.T) {
	req, app := login_test("admin@mail.com", "admin")
	resp, err := app.Test(req, 2000)
	if err != nil {
		t.Error(err)
	} else {
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		fmt.Println(resp)
		tokenData := struct {
			Token string `json:"token"`
		}{}
		if err := dec.Decode(&tokenData); err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodGet, "/user/", nil)
		req.Header.Add("Content", "application/json")
		req.Header.Add("Authorization", "Bearer "+ tokenData.Token)
		resp, err := app.Test(req, 2000)
		if err != nil {
			t.Error(err)
		}
		decoder := json.NewDecoder(resp.Body)
		var users []user.User
		if err := decoder.Decode(&users); err != nil {
			t.Error(err)
		}

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, 2, len(users))
		assert.Equal(t, "admin", users[0].Name)

		idStr := ""
		if oid, ok := users[1].ID.(primitive.ObjectID); ok {
			idStr = oid.Hex()
		}

		payloadStr := fmt.Sprintf(`{"name":"%s"}`, "kss")
		payload := []byte(payloadStr)
		req = httptest.NewRequest(http.MethodPut, "/user/"+idStr, bytes.NewReader(payload))
		req.Header.Add("Content", "application/json")
		req.Header.Add("Authorization", "Bearer "+tokenData.Token)
		resp, err = app.Test(req, 2000)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(resp)
		assert.Equal(t, 422, resp.StatusCode)
	}
}
