package snek

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMove1(t *testing.T) {
	jsonStr := "{\"you\":\"810286e5-a75b-4e7d-940d-035bbcab4656\",\"width\":20,\"turn\":54,\"snakes\":[{\"taunt\":\"\",\"name\":\"dsnek\",\"id\":\"810286e5-a75b-4e7d-940d-035bbcab4656\",\"health_points\":70,\"coords\":[[13,1],[13,2],[12,2],[12,3],[11,3],[10,3]]},{\"taunt\":\"Latitude is starting.\",\"name\":\"Latitude Snake\",\"id\":\"be0c9f05-2d28-4045-8d4c-561237b0d2e7\",\"health_points\":92,\"coords\":[[1,10],[1,9],[1,8],[1,7],[1,6]]},{\"taunt\":\"\",\"name\":\"SWU Bounty Snake\",\"id\":\"d1ea575f-bf5f-4d59-9117-cdaa644afd94\",\"health_points\":76,\"coords\":[[11,0],[11,1],[11,2],[10,2]]},{\"taunt\":\"I'm a snake, not a doctor\",\"name\":\"Genetisnake NBK\",\"id\":\"d9c06d7b-a3ad-4f67-83f8-2956cb1bd21a\",\"health_points\":57,\"coords\":[[3,2],[2,2],[1,2],[0,2]]}],\"height\":20,\"game_id\":\"38637171-78be-47fc-b177-b0611331d791\",\"food\":[[1,16],[17,13]],\"dead_snakes\":[]}"
	m := MoveRequest{}
	json.Unmarshal([]byte(jsonStr), &m)
	for i := 0; i < 100; i++ {
		assert.NotEqual(t, LEFT, m.GenerateMove())
	}
}

func TestMove2(t *testing.T) {
	jsonStr := "{\"you\":\"5505109c-b0cc-438f-91b0-1bdfcf5804c6\",\"width\":20,\"turn\":363,\"snakes\":[{\"taunt\":\"\",\"name\":\"dsnek\",\"id\":\"5505109c-b0cc-438f-91b0-1bdfcf5804c6\",\"health_points\":88,\"coords\":[[9,3],[9,4],[9,5],[10,5],[10,6],[11,6],[12,6],[13,6],[14,6],[15,6],[16,6],[17,6],[18,6],[18,7],[18,8],[18,9],[18,10],[18,11],[18,12],[18,13]]},{\"taunt\":\"\",\"name\":\"SWU Bounty Snake\",\"id\":\"50b96347-bcde-4a80-8899-1d47a8594b10\",\"health_points\":39,\"coords\":[[8,1],[8,0],[9,0],[10,0],[10,1],[10,2],[10,3],[10,4],[11,4],[11,5],[12,5],[13,5],[14,5],[15,5],[16,5],[16,4],[16,3],[16,2],[16,1],[15,1],[15,2],[14,2]]}],\"height\":20,\"game_id\":\"38637171-78be-47fc-b177-b0611331d791\",\"food\":[[16,0],[11,0]],\"dead_snakes\":[{\"taunt\":\"I'm a snake, not a doctor\",\"name\":\"Genetisnake NBK\",\"id\":\"f1481413-b868-45c0-b8b8-988ba3c72807\",\"health_points\":92,\"coords\":[[19,14],[19,15],[19,16],[18,16],[17,16],[16,16]]},{\"taunt\":\"Latitude is starting.\",\"name\":\"Latitude Snake\",\"id\":\"6c0cb2e5-76fb-4266-b481-a062ed1ea0d2\",\"health_points\":72,\"coords\":[[9,9],[9,10],[8,10]]}]}"
	m := MoveRequest{}
	json.Unmarshal([]byte(jsonStr), &m)
	for i := 0; i < 100; i++ {
		assert.NotEqual(t, UP, m.GenerateMove())
	}
}
