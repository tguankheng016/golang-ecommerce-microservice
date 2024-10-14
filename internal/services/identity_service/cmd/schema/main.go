package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ariga.io/atlas-provider-gorm/gormschema"
	rolemodel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	usermodel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
)

func main() {
	sb := &strings.Builder{}
	loadModels(sb)

	io.WriteString(os.Stdout, sb.String())
}

func loadModels(sb *strings.Builder) {
	models := []interface{}{
		&usermodel.User{},
		&rolemodel.Role{},
		&usermodel.UserRolePermission{},
		&usermodel.UserToken{},
	}
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	sb.WriteString(stmts)
	sb.WriteString(";\n")
}
