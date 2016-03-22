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
	"io"
	"io/ioutil"
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

func (game *Game) getCards() error {
	cardfile := "./nato_alph.txt"
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
	io.WriteString(os.Stdout, "Flash Cards!\n")
	game := &Game{}
	game.Cards = make(map[string]string)
	e := game.getCards()
	if e != nil {
		msg(e.Error())
		return
	}
	game.Keys = make([]string, len(game.Cards))
	idx := -1
	for k := range game.Cards {
		idx++
		game.Keys[idx] = k
	}

	flashes := len(game.Keys)
	tries := 0

	for len(game.Keys) > 0 {
		tries++
		key, keyIdx := game.getRandomKey()
		msg("\nHint: " + key)
		var answer string
		fmt.Scanln(&answer)
		//msg("" + strings.TrimSpace(answer) + "-a|")
		if game.Cards[key] == answer {
			msg("Correct!")
			game.Keys = append(game.Keys[:keyIdx], game.Keys[keyIdx+1:]...)
		} else {
			msg("Wrong!")
		}
	}
	fmt.Printf("\nIt took %d tries for %d questions\n", tries, flashes)
}