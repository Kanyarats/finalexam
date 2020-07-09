package customers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kanyarats/finalexam/database"
	"github.com/kanyarats/finalexam/middleware"
)

func createCustomersHandler(c *gin.Context) {
	cust := Customer{}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := database.Conn().QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", cust.Name, cust.Email, cust.Status)

	err := row.Scan(&cust.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, cust)
}

func getCustomersHandler(c *gin.Context) {

	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	customers := []*Customer{}
	for rows.Next() {
		cust := &Customer{}

		err = rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		customers = append(customers, cust)
	}

	c.JSON(http.StatusOK, customers)
}

func getCustomerByIDHandler(c *gin.Context) {
	id := c.Param("id")

	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cust := &Customer{}

	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func updateCustomersHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	cust := &Customer{}

	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := c.ShouldBindJSON(cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err = database.Conn().Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(id, cust.Name, cust.Email, cust.Status); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func deleteCustomersHandler(c *gin.Context) {
	id := c.Param("id")

	err := database.DeleteById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func SetupHandler() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.AuthMiddleware)

	r.GET("/customers", getCustomersHandler)
	r.GET("/customers/:id", getCustomerByIDHandler)
	r.POST("/customers", createCustomersHandler)
	r.PUT("/customers/:id", updateCustomersHandler)
	r.DELETE("/customers/:id", deleteCustomersHandler)

	return r
}
