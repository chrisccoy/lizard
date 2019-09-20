package main

import (
	"math/rand"
	"fmt"
	"time"
	"math"
)

const TTL_PLAYERS = 20
const TTL_ROUNDS = 3

type Player struct {
	name    string
	partner []*Player
}

var players []*Player

var partnerExclusion map[string]*Player
var playerPool map[string]*Player
var golfers = [20]string{"Chris Coy", "Rush Porter", "Alan Barber", "Andy Coy", "Chris Hamm",
	"Clark Seethaler", "Dan Smith", "Derek Fox", "Jeff Hull", "Jim Purdy",
	"Jon Seethaler", "Josh Peabody", "Les Tang", "Phil Quan", "Matt Lundahl",
	"Steve Seibly", "Tom Lincoln", "Tom Wilkerson", "Patrick McCormick", "Tony Bruser"}

func populatePool() {
	playerPool = make(map[string]*Player)
	for _, value := range players {
		playerPool[value.name] = value
	}

}
func main() {
	players = []*Player(make([]*Player, TTL_PLAYERS))
	for i, name := range golfers {
		players[i] = makePlayer(name)
	}
	for j := 0; j < TTL_ROUNDS; j++ {
		populatePool()
		for i, _ := range players {


			// If Player is not partnered already then find one
			if playerPool[players[i].name] != nil {
				if players[i].partner[j] == nil {
					//Remove self from the pool
					delete(playerPool, players[i].name)
					if len(playerPool) == 2 {
						//Edge Case to avoid playerA from taking playerB's only option
						next := 0
						exception := []*Player(make([]*Player, 2))
						//Move to the next
						idx := 0
						next = incrementSeed(i)
						for _, v := range playerPool {
							exception[idx] = v
							idx++
						}
						idx1 := 0
						idx2 := 1
						//player1:=i
						//player2:=next
						if !players[next].partnerExists(exception[0].name) {
							idx1 = 1
							idx2 = 0
						}

						players[i].partner[j] = exception[idx1]
						exception[idx1].partner[j] = players[i]
						players[next].partner[j] = exception[idx2]
						exception[idx2].partner[j] = players[next]
					} else {
						part := players[i].findPartner()
						//Link the players together
						players[i].partner[j] = part
						part.partner[j] = players[i]
						//Remove the partner from the pool
						delete(playerPool, part.name)
					}
				}
			}

		}
	}

	printPairings()
}
func printPairings() {
	for j := 0; j < TTL_ROUNDS; j++ {
		populatePool()
		fmt.Println()
		for len(playerPool) > 0 {
			player := getPlayerFromPool(rand.Intn(TTL_PLAYERS - 1))
			fmt.Printf("Round %d A Player: %s B Player: %s\n", j+1, player.name, player.partner[j].name)
			delete(playerPool, player.name)
			delete(playerPool, player.partner[j].name)
			/*
							partnerExclusion = make(map[string]*Player)
							for i, _ := range players {
								_, exists := partnerExclusion[players[i].name]
								if exists {
									continue
								}
								// Add the partner to the exlusion list, so we don't print twice
								partnerExclusion[players[i].partner[j].name] = players[i].partner[j]
							}*/
		}
	}

}
func makePlayer(name string) *Player {
	return &Player{name, make([]*Player, 3)}
}
func (p *Player) findPartner() *Player {

	var part *Player

	for {
		part = p.getRandomPartner()
		//Skip for self or if partner is already chosen for this player
		if p.name == part.name || p.partnerExists(part.name) {
			continue
		} else {
			break
		}
	}
	return part
}
func (p *Player) partnerExists(p1 string) bool {
	found := false
	x := playerPool[p1]

	if x != nil {

		for _, v := range p.partner {
			if v != nil && p1 == v.name {
				found = true
				break
			}
		}
	}
	return found
}
func (p *Player) getRandomPartner() *Player {
	return getPlayerFromPool(p.getRandomSeed())
}
func getPlayerFromPool(seed int) *Player {
	var thePlayer *Player
	for {
		thePlayer = playerPool[players[seed].name]
		if thePlayer != nil {
			break
		}
		// Crawl the list until we find a valid partner
		seed = incrementSeed(seed)
	}
	return thePlayer
}
func (p *Player) getRandomSeed() int {

	s1 := rand.NewSource(time.Now().UnixNano())
	num := s1.Int63()
	//Get the mod of the pseudo random number
	seed := int(math.Mod(float64(num), TTL_PLAYERS-1))
	//Skip for self
	if players[seed].name == p.name {
		seed = incrementSeed(seed)
	}
	return seed
}
func incrementSeed(seed int) int {
	if seed == TTL_PLAYERS-1 {
		seed = 0
	} else {
		seed++
	}
	return seed
}
