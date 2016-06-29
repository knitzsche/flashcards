/*
 * Copyright (C) 201^ Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * * Authored by:
 *	Kyle Nitzsche <kyle.nitzsche@canonical.com>
 */

package main

//#include <locale.h>
//#include <libintl.h>
import "C"

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func msg(str string) {
	fmt.Println(str)
}

type Game struct {
	Cards map[string]string
	Keys  []string
}

var game Game
var cardfile string

func (game *Game) getCards() error {
	content_b, e := ioutil.ReadFile(cardfile)
	if e != nil {
		msg_ := fmt.Sprintf("Sorry. %q does not exist", cardfile)
		return errors.New(msg_)
	}
	content := string(content_b)
	lines := strings.Split(content, "\n")
	idx := -1
	for _, line := range lines {
		kv := strings.Split(line, ":")
		if len(kv) < 2 {
			continue
		}
		if len(strings.TrimSpace(kv[0])) == 0 {
			continue
		}
		if len(strings.TrimSpace(kv[1])) == 0 {
			continue
		}
		game.Cards[kv[0]] = kv[1]
		idx++
		game.Keys = append(game.Keys, kv[0])
	}
	return nil
}

func (game *Game) showCards() error {
	return nil
}

func (game *Game) getRandomKey() (string, int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	keyNum := r.Intn(len(game.Keys))
	return game.Keys[keyNum], keyNum
}

func main() {
	if len(os.Args) > 1 {
		cardfile = os.Args[1]
	} else {
		cardfile = "./nato_alph.txt"
	}

	msg("==== cardfile: " + cardfile)

	if _, err := os.Stat(cardfile); os.IsNotExist(err) {
		msg("Cardfile does not exist")
	} else {
		msg("Cardfile does exist")
	}

	msg("\nFlash Cards!\n")
	game := &Game{}
	game.Cards = make(map[string]string)
	e := game.getCards()
	if e != nil {
		msg(e.Error())
		return
	}
	log.Printf("Cards: %v", game.Cards)
	game.Keys = make([]string, len(game.Cards)-1) //Card set description field
	idx := -1
	for k := range game.Cards {
		if k == "Card set description" {
			continue
		}
		idx++
		game.Keys[idx] = k
	}

	log.Printf("Keys: %v", game.Keys)
	flashes := len(game.Keys)
	tries := 0

	msg("Learning this: " + game.Cards["Card set description"])
	for len(game.Keys) > 0 {
		tries++
		key, keyIdx := game.getRandomKey()
		msg("\nHint: " + key)
		var answer string
		fmt.Scanln(&answer)
		//msg("" + strings.TrimSpace(answer) + "-a|")
		if strings.ToLower(game.Cards[key]) == strings.ToLower(answer) {
			msg("Correct!")
			game.Keys = append(game.Keys[:keyIdx], game.Keys[keyIdx+1:]...)
		} else {
			msg("Wrong! Correct answer is: " + game.Cards[key])
		}
	}
	fmt.Printf("\nIt took %d tries for %d questions\n", tries, flashes)
}
