package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

func searchPUT(c *gin.Context) {
	searchKey := os.Getenv("SEARCH_KEY")
	apiKey := os.Getenv("API_KEY")

	var requestData map[string]interface{}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searchTerms, ok := requestData["searchTerms"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing searchTerms"})
		return
	}

	startStr, ok := requestData["start"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing start"})
		return
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start"})
		return
	}

	searchParameters := formatSearchParameters(searchTerms)

	resp, err := http.Get("https://www.googleapis.com/customsearch/v1?key=" + apiKey + "&cx=" + searchKey + "&q=" + searchParameters + "&start=" + strconv.Itoa(start) + "&num=10")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to reach search api"})
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items, ok := data["items"].([]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected response format"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func formatSearchParameters(searchTerms string) string {
	return strings.ReplaceAll(strings.ToLower(searchTerms), " ", "+")
}
