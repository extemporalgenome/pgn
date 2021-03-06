package pgn

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"text/scanner"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type PGNSuite struct{}

var _ = Suite(&PGNSuite{})

var simple = `[Event "State Ch."]
[Site "New York, USA"]
[Date "1910.??.??"]
[Round "?"]
[White "Capablanca"]
[Black "Jaffe"]
[Result "1-0"]
[ECO "D46"]
[Opening "Queen's Gambit Dec."]
[Annotator "Reinfeld, Fred"]
[WhiteTitle "GM"]
[WhiteCountry "Cuba"]
[BlackCountry "United States"]

1. d4 d5 2. Nf3 Nf6 3. e3 c6 4. c4 e6 5. Nc3 Nbd7 6. Bd3 Bd6
7. O-O O-O 8. e4 dxe4 9. Nxe4 Nxe4 10. Bxe4 Nf6 11. Bc2 h6
12. b3 b6 13. Bb2 Bb7 14. Qd3 g6 15. Rae1 Nh5 16. Bc1 Kg7
17. Rxe6 Nf6 18. Ne5 c5 19. Bxh6+ Kxh6 20. Nxf7+ 1-0
`

func (s *PGNSuite) TestParse(c *C) {
	r := strings.NewReader(simple)
	sc := scanner.Scanner{}
	sc.Init(r)
	game, err := ParseGame(&sc)
	if err != nil {
		c.Fatal(err)
	}
	if game.Tags["Site"] != "New York, USA" {
		c.Fatal("Site tag wrong: ", game.Tags["Site"])
	}
	if len(game.Moves) == 0 || game.Moves[0].From != D2 || game.Moves[0].To != D4 {
		c.Fatal("first move is wrong", game.Moves[0])
	}
	if len(game.Moves) != 39 || game.Moves[38].From != E5 || game.Moves[38].To != F7 {
		c.Fatal("last move is wrong", game.Moves[38])
	}
}

func (s *PGNSuite) TestPGNScanner(c *C) {
	f, err := os.Open("polgar.pgn")
	if err != nil {
		c.Fatal(err)
	}
	ps := NewPGNScanner(f)
	for ps.Next() {
		game, err := ps.Scan()
		if err != nil {
			fmt.Println(game)
			c.Fatal(err)
		}
	}
}

func (s *PGNSuite) TestPGNParseWithCheckmate(c *C) {
	pgnstr := `[Event "Live Chess"]
[Site "Chess.com"]
[Date "2014.10.10"]
[White "MarkoMakaj"]
[Black "AndreyOstrovskiy"]
[Result "1-0"]
[WhiteElo "2196"]
[BlackElo "2226"]
[TimeControl "1|1"]
[Termination "MarkoMakaj won by checkmate"]

1.d4 g6 2.c4 Bg7 3.Nc3 c5 4.Nf3 cxd4 5.Nxd4 Nc6 6.Nc2 Nf6 7.g3 O-O 8.Bg2 b6 9.O-O Bb7 10.b3 Rc8
 11.Bb2 Qc7 12.Qd2 Qb8 13.Ne3 Rfd8 14.Rfd1 e6 15.Rac1 Qa8 16.Nb5 d5 17.cxd5 exd5 18.Bxf6 Bxf6 19.Nxd5 Bg7 20.e4 a6
 21.Nbc3 b5 22.Qf4 Qa7 23.Nf6+ Kh8 24.Ncd5 Nd4 25.Qh4 h6 26.Rxc8 Rxc8 27.e5 Ne6 28.Ng4 Rc2 29.Nde3 Rxa2 30.Nxh6 Bxg2
 31.Kxg2 Bxe5 32.Nxf7+ Kg7 33.Nxe5 Qxe3 34.Qe7+ Kh6 35.Nf7+ Kh5 36.Qh4# 1-0
`
	r := strings.NewReader(pgnstr)
	sc := scanner.Scanner{}
	sc.Init(r)
	game, err := ParseGame(&sc)
	c.Assert(err, IsNil)
	c.Assert(len(game.Moves), Equals, 71)
}

func (s *PGNSuite) TestPGNParseInfiniteLoopF4(c *C) {
	pgnstr := `[Event "BKL-Turnier"]
[Site "Leipzig"]
[Date "1984.??.??"]
[Round "5"]
[White "Polgar, Zsuzsa"]
[Black "Moehring, Guenther"]
[Result "1-0"]
[WhiteElo "2275"]
[BlackElo "2395"]
[ECO "A49"]

1.d4 Nf6 2.Nf3 d6 3.b3 g6 4.Bb2 Bg7 5.g3 c5 6.Bg2 cxd4 7.Nxd4 d5 8.O-O O-O
9.Na3 Re8 10.Nf3 Nc6 11.c4 dxc4 12.Nxc4 Be6 13.Rc1 Rc8 14.Nfe5 Nxe5 15.Bxe5 Bxc4
16.Rxc4 Rxc4 17.bxc4 Qa5 18.Bxf6 Bxf6 19.Bxb7 Rd8 20.Qb3 Rb8 21.e3 h5 22.Rb1 h4
23.Qb5 Qc7 24.a4 hxg3 25.hxg3 Be5 26.Kg2 Bd6 27.a5 Bc5 28.a6 Rd8 29.Qc6 Qxc6+
30.Bxc6 Rd2 31.Kf3 Rc2 32.Rb8+ Kg7 33.Bb5 Kf6 34.Rc8 Bb6 35.Ba4 Ra2 36.Bb5 Rc2
37.Ke4 e6 38.Kd3 Rc1 39.Kd2 Rb1 40.Kc2 Rb4 41.Rb8 Bc5 42.Rc8 Bb6 43.Rc6 Ba5
44.Rd6 g5 45.f4 gxf4 46.gxf4 Kf5 47.Rd7 Bb6 48.Rxf7+ Ke4 49.Rb7 Bc5 50.Kc3 Kxe3
51.Rc7 Bb6 52.Rc6 Ba5 53.Kc2 Kxf4 54.Rxe6 Bd8 55.Kc3 Rb1 56.Kd4 Rd1+ 57.Kc5 Kf5
58.Re8 Bb6+ 59.Kc6 Kf6 60.Kb7 Bg1 61.Ra8 Re1 62.Rf8+ Kg7 63.Rf5 Kg6 64.Rd5 Rc1
65.Ka8 Be3 66.Rd6+ Kf5 67.Rd3 Ke4 68.Rxe3+ Kxe3 69.Kxa7 Kd4 70.Kb6 Rg1 71.a7 Rg8
72.Kb7 Rg7+ 73.Kb6  1-0`

	r := strings.NewReader(pgnstr)
	sc := scanner.Scanner{}
	sc.Init(r)
	game, err := ParseGame(&sc)
	c.Assert(err, IsNil)
	//	fmt.Println(game)
	c.Assert(game.Tags["Site"], Equals, "Leipzig")
	c.Assert(len(game.Moves), Equals, 145)
}

func (s *PGNSuite) TestComments(c *C) {
	pgnstr := `[Event "Ch World (match)"]
[Site "New York (USA)"]
[Date "1886.03.24"]
[EventDate "?"]
[Round "19"]
[Result "0-1"]
[White "Johannes Zukertort"]
[Black "Wilhelm Steinitz"]
[ECO "D53"]
[WhiteElo "?"]
[BlackElo "?"]
[PlyCount "58"]

1. d4 {Notes by Robert James Fischer from a television
interview. } d5 2. c4 e6 3. Nc3 Nf6 4. Bg5 Be7 5. Nf3 O-O
6. c5 {White plays a mistake already; he should just play e3,
naturally.--Fischer} b6 7. b4 bxc5 8. dxc5 a5 9. a3 {Now he
plays this fantastic move; it's the winning move. -- Fischer}
d4 {He can't take with the knight, because of axb4.--Fischer}
10. Bxf6 gxf6 11. Na4 e5 {This kingside weakness is nothing;
the center is easily winning.--Fischer} 12. b5 Be6 13. g3 c6
14. bxc6 Nxc6 15. Bg2 Rb8 {Threatening Bb3.--Fischer} 16. Qc1
d3 17. e3 e4 18. Nd2 f5 19. O-O Re8 {A very modern move; a
quiet positional move. The rook is doing nothing now, but
later...--Fischer} 20. f3 {To break up the center, it's his
only chance.--Fischer} Nd4 21. exd4 Qxd4+ 22. Kh1 e3 23. Nc3
Bf6 24. Ndb1 d2 25. Qc2 Bb3 26. Qxf5 d1=Q 27. Nxd1 Bxd1
28. Nc3 e2 29. Raxd1 Qxc3 0-1`

	r := strings.NewReader(pgnstr)
	sc := scanner.Scanner{}
	sc.Init(r)
	game, err := ParseGame(&sc)
	c.Assert(err, Equals, nil)
	c.Assert(game, NotNil)
	c.Assert(game.Tags["Site"], Equals, "New York (USA)")
	c.Assert(len(game.Moves), Equals, 58)
}

func (s *PGNSuite) BenchmarkParse(c *C) {
	c.SetBytes(int64(len(simple)))
	r0 := strings.NewReader(simple)

	for i := 0; i < c.N; i++ {
		c.StopTimer()
		r := *r0
		sc := scanner.Scanner{}
		sc.Init(&r)

		c.StartTimer()
		ParseGame(&sc)
	}
}
