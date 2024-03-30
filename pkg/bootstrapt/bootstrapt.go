package bootstrapt

import (
	"github.com/fernandoDelPo/go_web_users/internal/domain"
	"github.com/fernandoDelPo/go_web_users/internal/user"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func NewDB() user.DB {
	return user.DB{
		Users: []domain.User{{
			ID:        1,
			FirstName: "fernando",
			LastName:  "Del Pozzi",
			Email:     "fernandodelpozzi@example.com",
		}, {
			ID:        2,
			FirstName: "Leticia",
			LastName:  "Caceres",
			Email:     "leticia@example.com",
		}, {
			ID:        3,
			FirstName: "Francesca",
			LastName:  "Del Pozzi",
			Email:     "fran@example.com",
		}, {
			ID:        4,
			FirstName: "Maxima",
			LastName:  "Perez",
			Email:     "mperez@example.com",
		}},
		MaxUserID: 4,
	}
}
