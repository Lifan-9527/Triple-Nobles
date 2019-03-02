package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

// Cards type
type Card struct{
	Color	string
	Shape	string
	Num		int
}

// players type
type gamer struct{
	Name	string
	Score	int
	Banker	bool
	Cards	[]int
	handPattern	string
	Money	float32
}

// card pattern type
var Pattern map[string]int = map[string]int{
	"none":						-1,
	"normal":					1,
	"8 o'clock":					2,
	"9 o'clock":					3,
	"Triple Noble":					4,
	"Puma":						6,
	"Triple Noble 8 o'clock":			7,
	"Triple Noble 9 o'clock":			8,
	"Triple Face!":					9,
	"beggar":					0,
}


// Rest card pool
var RestPool []int = make([]int,0)



// utility functions
func pop(key int, pool []int)[]int {
	length := len(pool)
	Restpool := make([]int,0)
	for i:=0;i<length;i++ {
		if i!=key {
			Restpool = append(Restpool,pool[i])
		}
	}
	return Restpool
} // pop out the assigned cards

func checkType(Cards []int)string{
	var dup []int = make([]int,3)
	dup[1] = Cards[1]
	dup[2] = Cards[2]
	dup[0] = Cards[0]
	sum := 0
	for i:=0;i<3;i++{
		if dup[i]>10{dup[i]=10}
		sum = sum+dup[i]
	}
	sum = sum%10

	if Cards[0]==Cards[1] && Cards[1]==Cards[2] && Cards[0]>10					{return "Triple Face!"}
	if sum==9 && Cards[0]==Cards[1] && Cards[1]==Cards[2] 						{return "Triple Noble 9 o'clock"}
	if sum==8 && Cards[0]==Cards[1] && Cards[1]==Cards[2]						{return "Triple Noble 8 o'clock"}
	if Cards[0]==Cards[1] && Cards[1]==Cards[2] && Cards[0]<11 && Cards[0]!=3	{return	"Puma"}
	if (Cards[0]!=Cards[1] || Cards[1]!=Cards[2]) && (Cards[0]>10&&Cards[1]>10&&Cards[2]>10)	{return "Triple Noble"}
	if sum==9 && (Cards[0]<11 || Cards[1]<11 || Cards[2]<11)					{return "9 o'clock"}
	if sum==8 && (Cards[0]<11 || Cards[1]<11 || Cards[2]<11)					{return "8 o'clock"
	} else{
		return "normal"
	}
}	// The index be out of range may never happen


func cpr2Players(Cards1 []int, Cards2 []int)int{
	pat1 := Pattern[checkType(Cards1)]
	pat2 := Pattern[checkType(Cards2)]
	if pat1>pat2 {
		return pat1
	}
	if pat1<pat2{
		return -pat2
	} else{
		return fineCpr(Cards1, Cards2)
	}

}

func fineCpr(Cards1 []int, Cards2 []int)int{
	sort.Ints(Cards1)
	sort.Ints(Cards2)
	for i:=0;i<3;i++{
		if Cards1[i]>Cards2[i]{return Pattern[checkType(Cards1)]}
		if Cards1[i]<Cards2[i]{return -Pattern[checkType(Cards2)]}
	}
	return 0	// This should never happen
}

func SmartPrint(i interface{}){
	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType :=reflect.TypeOf(i)
	for i:=0;i<vValue.NumField();i++{
		kv[vType.Field(i).Name] = vValue.Field(i)
	}
	//fmt.Println(":")
	for k,v :=range kv{
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}



func main() {
	rand.Seed(time.Now().Unix())
	var unitBet float32 = 0.5
	rounds := 100 // How many rounds can we play


	// players init
	var players []gamer = []gamer{
		gamer{"Yuxiang Lu", 0, false, []int{0,0,0},"none",0},
		gamer{"Huiyu Ding", 0, false, []int{0, 0, 0},"none",0},
		gamer{"Fan Li", 0, false, []int{0, 0, 0},"none",0},
		gamer{"Xiande Liu", 0, false, []int{0, 0, 0},"none",0},

	}
	//gamer{"Yichuan Miao", 0, false, []int{0, 0, 0},"none",0},

	//gamer{"Zhiqian Huang", 0, false, []int{0, 0, 0},"none",0},

	// poker init
	var poker [52]Card
	for i := 0; i < 52; i++ {
		if i%2 == 0 {
			poker[i].Color = "black"
		} else {
			poker[i].Color = "red"
		}
		if i/2%2 == 0 {
			poker[i].Shape = "heart"
		} else {
			poker[i].Shape = "square"
		}
		poker[i].Num = i/4 + 1
	}

	// Game start!


	// Elect the first Banker
	fortune := rand.Intn(len(players))
	bankerID := fortune
	players[bankerID].Banker = true
	// Dealing cards: draw cards from poker
	for ; rounds > 0; rounds-- {
		// generate a cardRest pool
		// and reset the cards to poker
		RestPool = RestPool[:0]
		for i := 0; i < 52; i++ {RestPool = append(RestPool,i)} // full pool
		// Players draw cards from rest card pool
		for i := 0; i < 3; i++ {
			for j := 0; j < len(players); j++ {

				cardDraw := rand.Intn(len(RestPool))
				RestPool = pop(cardDraw, RestPool)
				players[j].Cards[i] = poker[cardDraw].Num
			}
		}


		// Show hands
		for i:=0;i<len(players);i++{
			players[i].handPattern = checkType(players[i].Cards)
		}

		// Compare the cards to Banker
		for i:=0; i<len(players); i++{
			players[i].Score += cpr2Players(players[i].Cards,players[bankerID].Cards)
			players[bankerID].Score -=  cpr2Players(players[i].Cards,players[bankerID].Cards)
		}

		// change the Banker
		var key []int = make([]int,0)
		for i:=0;i<len(players);i++ {
			if players[i].handPattern == "9 o'clock" || players[i].handPattern == "Triple Noble 9 o'clock" {
				key = append(key, i)
			}
		}

		if len(key)==1{
			players[bankerID].Banker = false
			players[key[0]].Banker = true
			bankerID = key[0]
			}
		if len(key)>1 && (players[bankerID].handPattern=="9 o'clock" || players[bankerID].handPattern =="Triple Noble 9 o'clock"){
			potencialBankerID := bankerID
			for i := 0; i < len(key); i++ {
				sign := cpr2Players(players[potencialBankerID].Cards, players[i].Cards)
				if sign<0 {potencialBankerID = i}
			}
		}
		if len(key)>1 && players[bankerID].handPattern!="9 o'clock" && players[bankerID].handPattern!="Triple Noble 9 o'clock"{
			potencialBankerID := 0
			for i := 0; i < len(key); i++ {
				sign := cpr2Players(players[potencialBankerID].Cards, players[i].Cards)
				if sign<0 {potencialBankerID = i}
			}
		}

	}


	for i:=0;i<len(players);i++{
		players[i].Money = float32(players[i].Score) * unitBet
	}

	// end game
	fmt.Println("Vegas！！！:")
	for i:=0;i<len(players);i++{
		SmartPrint(players[i])
		fmt.Println("")
	}



}
