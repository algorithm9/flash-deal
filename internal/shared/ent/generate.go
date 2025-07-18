package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/lock,sql/execquery,sql/upsert,sql/modifier ./schema --target ./gen
