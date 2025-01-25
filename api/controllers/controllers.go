package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/daiki-trnsk/test-debatemap/api/database"
	"github.com/daiki-trnsk/test-debatemap/api/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var DebateMapCollection *mongo.Collection = database.DebateMapData(database.Client, "testdebatemaps")

func AddDebateMap(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var debatemap models.DebateMap
	defer cancel()
	if err := c.Bind(&debatemap); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Printf("Received DebateMap: %+v\n", debatemap)
	debatemap.ID = primitive.NewObjectID()
	debatemap.RegistrationDate = time.Now()
	debatemap.NodesJSON = map[string]interface{}{
		"id":       "root",
		"topic":    debatemap.RootNodeTopic,
		"expanded": true,
		"children": []interface{}{},
	}
	if debatemap.NodesJSON["topic"] == "" {
		return c.JSON(http.StatusBadRequest, "Enter Topic")
	}
	_, anyerr := DebateMapCollection.InsertOne(ctx, debatemap)
	if anyerr != nil {
		return c.JSON(http.StatusInternalServerError, "Not Created")
	}
	defer cancel()
	return c.JSON(http.StatusOK, "Map Created!")
}

func GetDebateMaps(c echo.Context) error {
	var debatemaplist []models.DebateMap
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	cursor, err := DebateMapCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
	}
	err = cursor.All(ctx, &debatemaplist)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, "Error occurred while fetching debatemaps")
	}
	defer cursor.Close(ctx)
	if err := cursor.Err(); err != nil {
		// Don't forget to log errors. I log them really simple here just
		// to get the point across.
		log.Println(err)
		return c.JSON(400, "invalid")
	}
	defer cancel()
	return c.JSON(200, debatemaplist)
}

func GetDebateMap(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var debatemap models.DebateMap
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	err := DebateMapCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&debatemap)
	if err != nil {
		return c.JSON(404, "Invalid DebateMap ID")
	}
	defer cancel()
	return c.JSON(200, debatemap)
}

func UpdateDebateMap(c echo.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var debatemap models.DebateMap
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := c.Bind(&debatemap); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	update := bson.M{
		"$set": bson.M{
			"root_node_topic": debatemap.NodesJSON["topic"],
			"nodes_json": debatemap.NodesJSON,
		},
	}
	_, err := DebateMapCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error occurred while updating debatemap"})
	}
	defer cancel()
	return c.JSON(http.StatusOK, "DebateMap updated successfully")
}

func SearchDebateMapByQuery(c echo.Context) error {
	var searchdebatemaps []models.DebateMap
	queryParam := c.QueryParam("name")
	if queryParam == "" {
		log.Println("query is empty")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid Search Index"})
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	searchquerydb, err := DebateMapCollection.Find(ctx, bson.M{"root_node_topic": bson.M{"$regex": queryParam, "$options": "i"}})
	if err != nil {
		return c.JSON(404, "something went wrong in fetching the dbquery")
	}
	err = searchquerydb.All(ctx, &searchdebatemaps)
	if err != nil {
		log.Println(err)
		return c.JSON(400, "invalid")
	}
	defer searchquerydb.Close(ctx)
	if err := searchquerydb.Err(); err != nil {
		log.Println(err)
		return c.JSON(400, "invalid request")
	}
	defer cancel()
	return c.JSON(200, searchdebatemaps)
}
