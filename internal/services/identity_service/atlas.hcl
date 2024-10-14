data "external_schema" "gorm" {
  program = ["go", "run", "./cmd/schema/main.go"]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/16/dev"
  migration {
    dir = "file://data/migrations?format=goose"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}