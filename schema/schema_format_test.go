package schema_test

import (
	"bytes"
	"github.com/chirino/graphql/internal/example/starwars"
	"github.com/chirino/graphql/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteSchemaFormatForStarwars(t *testing.T) {
	s := schema.New()
	s.Parse(starwars.Schema)
	buf := &bytes.Buffer{}
	s.WriteSchemaFormat(buf)
	assert.Equal(t, `"A character from the Star Wars universe"
interface Character {
  "The ID of the character"
  id:ID!
  "The name of the character"
  name:String!
  "The friends of the character, or an empty list if they have none"
  friends:[Character]
  "The friends of the character exposed as a connection with edges"
  friendsConnection(first:Int, after:ID):FriendsConnection!
  "The movies this character appears in"
  appearsIn:[Episode!]!
}
"An autonomous mechanical character in the Star Wars universe"
type Droid implements Character  {
  "The ID of the droid"
  id:ID!
  "What others call this droid"
  name:String!
  "This droid's friends, or an empty list if they have none"
  friends:[Character]
  "The friends of the droid exposed as a connection with edges"
  friendsConnection(first:Int, after:ID):FriendsConnection!
  "The movies this droid appears in"
  appearsIn:[Episode!]!
  "This droid's primary function"
  primaryFunction:String
}
"The episodes in the Star Wars trilogy"
enum Episode {
  "Star Wars Episode IV: A New Hope, released in 1977."
  NEWHOPE
  "Star Wars Episode V: The Empire Strikes Back, released in 1980."
  EMPIRE
  "Star Wars Episode VI: Return of the Jedi, released in 1983."
  JEDI
}
"A connection object for a character's friends"
type FriendsConnection {
  "The total number of friends"
  totalCount:Int!
  "The edges for each of the character's friends."
  edges:[FriendsEdge]
  "A list of the friends, as a convenience when edges are not needed."
  friends:[Character]
  "Information for paginating this connection"
  pageInfo:PageInfo!
}
"An edge object for a character's friends"
type FriendsEdge {
  "A cursor used for pagination"
  cursor:ID!
  "The character represented by this friendship edge"
  node:Character
}
"A humanoid creature from the Star Wars universe"
type Human implements Character  {
  "The ID of the human"
  id:ID!
  "What this human calls themselves"
  name:String!
  "Height in the preferred unit, default is meters"
  height(unit:LengthUnit=METER):Float!
  "Mass in kilograms, or null if unknown"
  mass:Float
  "This human's friends, or an empty list if they have none"
  friends:[Character]
  "The friends of the human exposed as a connection with edges"
  friendsConnection(first:Int, after:ID):FriendsConnection!
  "The movies this human appears in"
  appearsIn:[Episode!]!
  "A list of starships this person has piloted, or an empty list if none"
  starships:[Starship]
}
"Units of height"
enum LengthUnit {
  "The standard unit around the world"
  METER
  "Primarily used in the United States"
  FOOT
}
"The mutation type, represents all updates we can make to our data"
type Mutation {
  createReview(episode:Episode!, review:ReviewInput!):Review
}
"Information for paginating this connection"
type PageInfo {
  startCursor:ID
  endCursor:ID
  hasNextPage:Boolean!
}
"The query type, represents all of the entry points into our object graph"
type Query {
  hero(episode:Episode=NEWHOPE):Character
  reviews(episode:Episode!):[Review]!
  search(text:String!):[SearchResult]!
  character(id:ID!):Character
  droid(id:ID!):Droid
  human(id:ID!):Human
  starship(id:ID!):Starship
}
"Represents a review for a movie"
type Review {
  "The number of stars this review gave, 1-5"
  stars:Int!
  "Comment about the movie"
  commentary:String
}
"The input object sent when someone is creating a new review"
input ReviewInput {
  "0-5 stars"
  stars:Int!
  "Comment about the movie, optional"
  commentary:String
}
union SearchResult = Human | Droid | Starship
type Starship {
  "The ID of the starship"
  id:ID!
  "The name of the starship"
  name:String!
  "Length of the starship, along the longest axis"
  length(unit:LengthUnit=METER):Float!
}
schema {
  mutation: Mutation
  query: Query
}
`, buf.String())
}

func TestWriteSchemaFormatEdgeCases(t *testing.T) {
	s := schema.New()
	s.Parse(`
"""
Multi
Line Description
"""
directive @db_table(
    name: String
) on OBJECT

scalar Revision

schema {
  query: Query
}

type Query @db_table(name:"Hello") {
    hi: Revision
    args(a:String=null, b:Int=5, c:String="Hi", d:[String]=["a", "b"]): String
}

`)
	buf := &bytes.Buffer{}
	s.WriteSchemaFormat(buf)
	assert.Equal(t, `"""
Multi
Line Description
"""
directive @db_table(name:String) on OBJECT
type Query @db_table(name:"Hello") {
  hi:Revision
  args(a:String=null, b:Int=5, c:String="Hi", d:[String]=["a", "b"]):String
}
scalar Revision
schema {
  query: Query
}
`, buf.String())
}
