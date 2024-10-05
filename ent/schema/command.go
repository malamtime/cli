package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Command holds the schema definition for the Command entity.
type Command struct {
	ent.Schema
}

func (Command) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Command.
func (Command) Fields() []ent.Field {
	return []ent.Field{
		field.String("shell"),
		field.Int64("sessionId"),
		field.String("command"),
		field.String("main").Default(""),
		field.String("hostname"),
		field.String("username"),
		field.Time("time"),
		field.Time("endTime"),
		field.Int("result").Default(0),
		field.Enum("phase").Values("pre", "post").Default("post"),
		field.Bool("sentToServer").Default(false),
	}
}

// Edges of the Command.
func (Command) Edges() []ent.Edge {
	return nil
}
