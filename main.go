package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const screenS = 1000

var running bool = true

var light = sdl.Color{R: 255, G: 255, B: 255, A: 255}
var dark = sdl.Color{R: 120, G: 142, B: 173, A: 255}
var oldSquareColor = sdl.Color{R: 146, G: 177, B: 102, A: 255}
var newSquareColor = sdl.Color{R: 195, G: 216, B: 135, A: 255}
var selectedSquareColor = sdl.Color{R: 121, G: 156, B: 130, A: 255}
var board [64]*Square

var whiteCastleShort bool = true
var whiteCastleLong bool = true
var blackCastleShort bool = true
var blackCastleLong bool = true

var pieces [32]*Piece = [32]*Piece{&kw, &qw, &rw1, &rw2, &nw1, &nw2, &bw1, &bw2, &pw1, &pw2, &pw3, &pw4, &pw5, &pw6, &pw7, &pw8, &kb, &qb, &rb1, &rb2, &nb1, &nb2, &bb1, &bb2, &pb1, &pb2, &pb3, &pb4, &pb5, &pb6, &pb7, &pb8}

var selectedSquare *Square
var selectedPiece *Piece

var turn int = 1 //1 - white to move, 0 - black to move

var cemetary *Square = &Square{X: -1, Y: -1, Rectangle: nil, isOccupiedWhite: false, isOccupiedBlack: false}
var lastMove Move = Move{piece: &Piece{Type: 0, Color: -1, Pos: cemetary}, oldSquare: &Square{Rectangle: &sdl.Rect{X: screenS, Y: screenS, H: 0, W: 0}}, newSquare: &Square{Rectangle: &sdl.Rect{X: screenS, Y: screenS, H: 0, W: 0}}}

var klTex *sdl.Texture
var qlTex *sdl.Texture
var rlTex *sdl.Texture
var nlTex *sdl.Texture
var blTex *sdl.Texture
var plTex *sdl.Texture

var kdTex *sdl.Texture
var qdTex *sdl.Texture
var rdTex *sdl.Texture
var ndTex *sdl.Texture
var bdTex *sdl.Texture
var pdTex *sdl.Texture

var queenPromotionButton Button
var rookPromotionButton Button
var knightPromotionButton Button
var bishopPromotionButton Button

var queenLightPromotionTex *sdl.Texture
var rookLightPromotionTex *sdl.Texture
var knightLightPromotionTex *sdl.Texture
var bishopLightPromotionTex *sdl.Texture
var queenDarkPromotionTex *sdl.Texture
var rookDarkPromotionTex *sdl.Texture
var knightDarkPromotionTex *sdl.Texture
var bishopDarkPromotionTex *sdl.Texture

var renderer *sdl.Renderer
var window *sdl.Window
var err error

var kw Piece
var qw Piece
var rw1 Piece
var rw2 Piece
var nw1 Piece
var nw2 Piece
var bw1 Piece
var bw2 Piece
var pw1 Piece
var pw2 Piece
var pw3 Piece
var pw4 Piece
var pw5 Piece
var pw6 Piece
var pw7 Piece
var pw8 Piece
var kb Piece
var qb Piece
var rb1 Piece
var rb2 Piece
var nb1 Piece
var nb2 Piece
var bb1 Piece
var bb2 Piece
var pb1 Piece
var pb2 Piece
var pb3 Piece
var pb4 Piece
var pb5 Piece
var pb6 Piece
var pb7 Piece
var pb8 Piece

var promotion int

var toBeDestroyed *Piece = nil

var choosing bool = false

var promoted bool = false

var posSave *Square

type Button struct {
	Rect *sdl.Rect
	Tex  *sdl.Texture
}

type Piece struct {
	Type  int //1 - king, 2 - queen, 3  - rook, 4 - knight, 5 - bishop, 6 - pawn
	Color int //1 - white, 0 - black
	Pos   *Square
}

type Square struct {
	X               int32
	Y               int32
	Rectangle       *sdl.Rect
	isOccupiedWhite bool
	isOccupiedBlack bool
}

type Move struct {
	piece     *Piece
	oldSquare *Square
	newSquare *Square
}

func initBoard() [64]*Square {
	var s [64]*Square
	var k int = 0

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			s[k] = &Square{X: int32(i), Y: int32(j), Rectangle: &sdl.Rect{X: int32(i * screenS / 8), Y: int32(j * screenS / 8), H: int32(screenS / 8), W: int32(screenS / 8)}}
			k += 1
		}
	}
	return s
}

func drawPiece(self Piece) {
	if self.Pos != cemetary {
		if self.Color == 1 {
			if self.Type == 1 {
				renderer.Copy(klTex, nil, self.Pos.Rectangle)
				return
			} else if self.Type == 2 {
				renderer.Copy(qlTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 3 {
				renderer.Copy(rlTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 4 {
				renderer.Copy(nlTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 5 {
				renderer.Copy(blTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 6 {
				renderer.Copy(plTex, nil, self.Pos.Rectangle)
				return

			}
		} else if self.Color == 0 {
			if self.Type == 1 {
				renderer.Copy(kdTex, nil, self.Pos.Rectangle)
				return
			} else if self.Type == 2 {
				renderer.Copy(qdTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 3 {
				renderer.Copy(rdTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 4 {
				renderer.Copy(ndTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 5 {
				renderer.Copy(bdTex, nil, self.Pos.Rectangle)
				return

			} else if self.Type == 6 {
				renderer.Copy(pdTex, nil, self.Pos.Rectangle)
				return

			}
		}

	}

}
func move(self *Piece, t *Square, test bool) error {
	lastMoveSave := lastMove
	var castle string = ""
	var rookPrevPos *Square
	if !test {
		toBeDestroyed = nil
		promotion = 0
		promoted = false
	}

	switch self.Type {
	case -1:
		return errors.New("this piece is taken")
	case 1: //king
		{
			//castling
			if !test && !isChecked() {
				if self.Color == 1 {
					if t.X == 6 && t.Y == 7 && whiteCastleShort && !board[47].isOccupiedWhite && !board[55].isOccupiedWhite && !board[47].isOccupiedBlack && !board[55].isOccupiedBlack {
						if !isAttacked(board[47]) && !isAttacked(board[55]) {
							castle = "whiteshort"
						}
					} else if t.X == 2 && t.Y == 7 && whiteCastleLong && !board[15].isOccupiedWhite && !board[23].isOccupiedWhite && !board[31].isOccupiedWhite && !board[15].isOccupiedBlack && !board[23].isOccupiedBlack && !board[31].isOccupiedBlack {
						if !isAttacked(board[23]) && !isAttacked(board[31]) {
							castle = "whitelong"
						}
					}
				} else if self.Color == 0 {
					if t.X == 6 && t.Y == 0 && blackCastleShort && !board[40].isOccupiedWhite && !board[48].isOccupiedWhite && !board[40].isOccupiedBlack && !board[48].isOccupiedBlack {
						if !isAttacked(board[40]) && !isAttacked(board[48]) {
							castle = "blackshort"
						}
					} else if t.X == 2 && t.Y == 0 && blackCastleLong && !board[8].isOccupiedWhite && !board[16].isOccupiedWhite && !board[24].isOccupiedWhite && !board[8].isOccupiedBlack && !board[16].isOccupiedBlack && !board[24].isOccupiedBlack {
						if !isAttacked(board[16]) && !isAttacked(board[24]) {
							castle = "blacklong"
						}
					}
				}
			}
			//checks for the move rules of the king
			if ((self.Pos.X == t.X) && (self.Pos.Y+1 == t.Y || self.Pos.Y-1 == t.Y)) || ((self.Pos.Y == t.Y) && (self.Pos.X+1 == t.X || self.Pos.X-1 == t.X)) || (int32(math.Abs(float64(self.Pos.X-t.X))) == 1 && int32(math.Abs(float64(self.Pos.Y-t.Y))) == 1) {
				//checks if the square is occupied by a same colored piece
				if (t.isOccupiedWhite && self.Color == 1) || (t.isOccupiedBlack && self.Color == 0) {
					return errors.New("square occupied by an ally piece")
				}
				//checks if the square is occupied by an enemy piece
				if (t.isOccupiedBlack && self.Color == 1) || (t.isOccupiedWhite && self.Color == 0) {
					//looks through all pieces to take the one thats on the target square
					for _, piece := range pieces {
						if piece.Pos == t {
							if piece.Type != 1 {
								toBeDestroyed = piece
							} else {
								return errors.New("king")
							}
						}
					}
				}
				if self.Color == 1 {
					whiteCastleLong = false
					whiteCastleShort = false
				} else {
					blackCastleShort = false
					blackCastleLong = false
				}
				//}
			} else if castle == "" {
				return errors.New("illegal king move")
			}
		}
	case 2: //queen
		{
			if (self.Pos.X == t.X) || (self.Pos.Y == t.Y) || (int32(math.Abs(float64(self.Pos.X-t.X))) == int32(math.Abs(float64(self.Pos.Y-t.Y)))) {
				if self.Pos.X == t.X {
					if self.Pos.Y > t.Y {

						for i := self.Pos.Y - 1; i > t.Y; i-- {
							for j := 0; j < len(board); j++ {
								if board[j].X == self.Pos.X && board[j].Y == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("queen blocked")
									}
								}
							}
						}
					} else if self.Pos.Y < t.Y {
						for i := self.Pos.Y + 1; i < t.Y; i++ {
							for j := 0; j < len(board); j++ {
								if board[j].X == self.Pos.X && board[j].Y == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("queen blocked")
									}
								}
							}
						}
					}
				} else if self.Pos.Y == t.Y {
					if self.Pos.X > t.X {

						for i := self.Pos.X - 1; i > t.X; i-- {
							for j := 0; j < len(board); j++ {
								if board[j].Y == self.Pos.Y && board[j].X == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("queen blocked")
									}
								}
							}
						}
					} else if self.Pos.X < t.X {
						for i := self.Pos.X + 1; i < t.X; i++ {
							for j := 0; j < len(board); j++ {
								if board[j].Y == self.Pos.Y && board[j].X == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("queen blocked")
									}
								}
							}
						}
					}
				} else if int32(math.Abs(float64(self.Pos.X-t.X))) == int32(math.Abs(float64(self.Pos.Y-t.Y))) {
					n := int32(math.Abs(float64(self.Pos.X - t.X)))
					xd := (t.X - self.Pos.X) / n
					yd := (t.Y - self.Pos.Y) / n
					for k := 1; int32(k) < n; k++ {
						for i := 0; i < len(board); i++ {
							if (self.Pos.X+(int32(k)*xd)) == board[i].X && (self.Pos.Y+(int32(k)*yd)) == board[i].Y {
								if board[i].isOccupiedBlack || board[i].isOccupiedWhite {
									return errors.New("queen blocked")
								}
							}
						}
					}

				}
				if (t.isOccupiedBlack && self.Color == 0) || (t.isOccupiedWhite && self.Color == 1) {
					return errors.New("square occupied by an ally piece")
				}
				if (t.isOccupiedBlack && self.Color == 1) || (t.isOccupiedWhite && self.Color == 0) {
					for _, piece := range pieces {
						if piece.Pos == t {
							if piece.Type != 1 {
								toBeDestroyed = piece
							} else {
								return errors.New("king")
							}
						}
					}
				}

			} else {
				return errors.New("illegal queen move")
			}

		}
	case 3: //rook
		{
			if self.Pos.X == t.X || self.Pos.Y == t.Y {
				if self.Pos.X == t.X {
					if self.Pos.Y > t.Y {

						for i := self.Pos.Y - 1; i > t.Y; i-- {
							for j := 0; j < len(board); j++ {
								if board[j].X == self.Pos.X && board[j].Y == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("rook blocked")
									}
								}
							}
						}
					} else if self.Pos.Y < t.Y {
						for i := self.Pos.Y + 1; i < t.Y; i++ {
							for j := 0; j < len(board); j++ {
								if board[j].X == self.Pos.X && board[j].Y == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("rook blocked")
									}
								}
							}
						}
					}
				} else if self.Pos.Y == t.Y {
					if self.Pos.X > t.X {

						for i := self.Pos.X - 1; i > t.X; i-- {
							for j := 0; j < len(board); j++ {
								if board[j].Y == self.Pos.Y && board[j].X == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("rook blocked")
									}
								}
							}
						}
					} else if self.Pos.X < t.X {
						for i := self.Pos.X + 1; i < t.X; i++ {
							for j := 0; j < len(board); j++ {
								if board[j].Y == self.Pos.Y && board[j].X == i {
									if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
										return errors.New("rook blocked")
									}
								}
							}
						}
					}
				}
				if (t.isOccupiedWhite && self.Color == 1) || (t.isOccupiedBlack && self.Color == 0) {
					return errors.New("square occupied by an ally piece")
				}
				if (t.isOccupiedBlack && self.Color == 1) || (t.isOccupiedWhite && self.Color == 0) {
					for _, piece := range pieces {
						if piece.Pos == t {
							if piece.Type != 1 {
								toBeDestroyed = piece
							} else {
								return errors.New("king")
							}

						}
					}
				}
				if self == &rw1 {
					whiteCastleLong = false
				} else if self == &rw2 {
					whiteCastleShort = false
				} else if self == &rb1 {
					blackCastleLong = false
				} else if self == &rb2 {
					blackCastleShort = false
				}

			} else {
				return errors.New("illegal rook move")
			}
		}
	case 4: //knight
		{
			if ((int32(math.Abs(float64(self.Pos.X-t.X)))) == 2 && (int32(math.Abs(float64(self.Pos.Y-t.Y))) == 1)) || ((int32(math.Abs(float64(self.Pos.Y-t.Y)))) == 2 && (int32(math.Abs(float64(self.Pos.X-t.X))) == 1)) {
				if (t.isOccupiedBlack && self.Color == 0) || (t.isOccupiedWhite && self.Color == 1) {
					return errors.New("square occupied by an ally piece")
				}
				if (t.isOccupiedBlack && self.Color == 1) || (t.isOccupiedWhite && self.Color == 0) {
					for _, piece := range pieces {
						if piece.Pos == t {
							if piece.Type != 1 {
								toBeDestroyed = piece
							} else {
								return errors.New("king")
							}
						}
					}
				}
			} else {
				return errors.New("illegal knight move")
			}
		}
	case 5: //bishop
		{
			if int32(math.Abs(float64(self.Pos.X-t.X))) == int32(math.Abs(float64(self.Pos.Y-t.Y))) {
				//checking if the piece is blocked
				n := int32(math.Abs(float64(self.Pos.X - t.X)))
				var xd int32
				var yd int32
				if n != 0 {
					xd = (t.X - self.Pos.X) / n
					yd = (t.Y - self.Pos.Y) / n
				} else {
					return errors.New("invalid move")
				}

				for k := 1; int32(k) < n; k++ {
					for i := 0; i < len(board); i++ {
						if (self.Pos.X+(int32(k)*xd)) == board[i].X && (self.Pos.Y+(int32(k)*yd)) == board[i].Y {
							if board[i].isOccupiedBlack || board[i].isOccupiedWhite {
								return errors.New("bishop blocked")
							}
						}
					}
				}
				if (t.isOccupiedBlack && self.Color == 0) || (t.isOccupiedWhite && self.Color == 1) {
					return errors.New("square occupied by an ally piece")
				}
				if (t.isOccupiedBlack && self.Color == 1) || (t.isOccupiedWhite && self.Color == 0) {
					for _, piece := range pieces {
						if piece.Pos == t {
							if piece.Type != 1 {
								toBeDestroyed = piece
							} else {
								return errors.New("king")
							}
						}
					}
				}
			} else {
				return errors.New("illegal bishop move")
			}

		}
	case 6: //pawn
		{
			var blocked bool = false

			//checking for normal move / two squares
			if ((self.Pos.X == t.X) && ((self.Color == 1 && self.Pos.Y-1 == t.Y) || (self.Color == 0 && self.Pos.Y+1 == t.Y)) || ((self.Color == 1 && self.Pos.Y == 6 && self.Pos.Y-2 == t.Y) || (self.Color == 0 && self.Pos.Y == 1 && self.Pos.Y+2 == t.Y)) && self.Pos.X == t.X) && (!t.isOccupiedBlack && !t.isOccupiedWhite) {
				if t.isOccupiedBlack && t.isOccupiedWhite {
					return errors.New("square occupied by an ally piece")
				}
				if self.Color == 1 {
					for i := self.Pos.Y - 1; i > t.Y; i-- {
						for j := 0; j < len(board); j++ {
							if board[j].X == self.Pos.X && board[j].Y == i {
								if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
									blocked = true
									break
								}
							}
						}
					}
				} else if self.Color == 0 {
					for i := self.Pos.Y + 1; i < t.Y; i++ {
						for j := 0; j < len(board); j++ {
							if board[j].X == self.Pos.X && board[j].Y == i {
								if board[j].isOccupiedBlack || board[j].isOccupiedWhite {
									blocked = true
									break
								}
							}
						}
					}
				}
				if blocked {
					return errors.New("pawn move blocked")
				}

				//cheking for taking move
			} else if (int32(math.Abs(float64(self.Pos.X-t.X))) == 1 && int32(math.Abs(float64(self.Pos.Y-t.Y))) == 1) && ((self.Color == 0 && !t.isOccupiedBlack && t.isOccupiedWhite) || (self.Color == 1 && !t.isOccupiedWhite && t.isOccupiedBlack)) {
				for _, piece := range pieces {
					if piece.Pos == t {
						if piece.Type != 1 {
							toBeDestroyed = piece
						} else {
							return errors.New("king")
						}

					}
				}
				//holy hell
			} else if (int32(math.Abs(float64(self.Pos.X-t.X))) == 1 && int32(math.Abs(float64(self.Pos.Y-t.Y))) == 1) &&
				(((self.Color == 0 && self.Pos.Y == 4) || (self.Color == 1 && self.Pos.Y == 3)) && (lastMove.piece.Type == 6) && (lastMove.piece.Pos.X+1 == self.Pos.X || lastMove.piece.Pos.X-1 == self.Pos.X) && (lastMove.piece.Pos.Y == self.Pos.Y) && (math.Abs(float64(lastMove.oldSquare.Y-lastMove.newSquare.Y)) == 2)) {
				toBeDestroyed = lastMove.piece
			} else {
				return errors.New("not legal move")
			}
			//checking for promotion
			if (t.Y == 0 && self.Color == 1) || (t.Y == 7 && self.Color == 0) {
				promotion = choosePromotion(self)
				if promotion == 0 {
					return errors.New("promotion cancelled")
				}
				promoted = true
				self.Type = promotion
			}
		}
	}
	lastMove = Move{piece: self, oldSquare: self.Pos, newSquare: t}
	self.Pos = t
	if castle != "" {
		if castle == "whiteshort" {
			rookPrevPos = rw2.Pos
			rw2.Pos = board[47]
		} else if castle == "whitelong" {
			rookPrevPos = rw1.Pos
			rw1.Pos = board[31]
		} else if castle == "blackshort" {
			rookPrevPos = rb2.Pos
			rb2.Pos = board[40]
		} else if castle == "blacklong" {
			rookPrevPos = rb1.Pos
			rb1.Pos = board[24]
		}
	}
	updateBoard()
	if test {
		if toBeDestroyed != nil {
			if toBeDestroyed.Color == 1 {
				toBeDestroyed.Pos.isOccupiedWhite = false
			} else if toBeDestroyed.Color == 0 {
				toBeDestroyed.Pos.isOccupiedBlack = false
			}
		}
	}
	inCheck := isChecked()
	if !test && !inCheck {
		if castle != "" {
			if self.Color == 1 {
				whiteCastleLong = false
				whiteCastleShort = false

			} else if self.Color == 0 {
				blackCastleLong = false
				blackCastleShort = false
			}
		}
		if toBeDestroyed != nil {
			destroyPiece(toBeDestroyed)
		}
		updateBoard()
		if turn == 0 {
			turn = 1
		} else if turn == 1 {
			turn = 0
		}
		return nil
	} else {
		if castle != "" {
			if castle == "whiteshort" {
				rw2.Pos = rookPrevPos
			} else if castle == "whitelong" {
				rw1.Pos = rookPrevPos
			} else if castle == "blackshort" {
				rb2.Pos = rookPrevPos
			} else if castle == "blacklong" {
				rb1.Pos = rookPrevPos
			}
		}
		self = lastMove.piece
		self.Pos = lastMove.oldSquare
		lastMove = lastMoveSave
		//updateBoard()
		if inCheck {
			return errors.New("king left in check")
		} else if test {
			return errors.New("test")
		}
		// return errors.New("test or king left in check")
	}
	return errors.New("unreachable")
}

func isChecked() bool {
	var kingToBeChecked *Piece
	if turn == 1 {
		kingToBeChecked = &kw
	} else {
		kingToBeChecked = &kb
	}
	for _, piece := range pieces {
		if piece.Color != turn && piece != toBeDestroyed && piece.Type != 6 {
			err := move(piece, kingToBeChecked.Pos, true)
			if err != nil {
				if err.Error() == "king" {
					return true
				}
			}
		} else if piece.Color != turn && piece != toBeDestroyed && piece.Type == 6 && piece.Pos.X != kingToBeChecked.Pos.X {
			err := move(piece, kingToBeChecked.Pos, true)
			if err != nil {
				if err.Error() == "king" {
					return true
				}
			}
		}
	}
	return false
}

// func avaibleMoves() int {
// 	boardSave := board
// 	piecesSave := pieces
// 	numMoves := 0
// 	for _, piece := range pieces {
// 		for _, square := range board {
// 			if square != piece.Pos && piece.Color != turn {
// 				err := move(piece, square, true)
// 				//fmt.Println(piece, ",", square.X, ",", square.Y, ",", err)
// 				if err.Error() == "test" {
// 					numMoves += 1
// 				}
// 			}
// 		}
// 	}
// 	board = boardSave
// 	pieces = piecesSave
// 	return numMoves
// }

func isAttacked(s *Square) bool {
	for _, piece := range pieces {
		if piece.Color != turn && piece != toBeDestroyed {
			err := move(piece, s, true)
			if err.Error() == "test" || err.Error() == "king left in check" {
				fmt.Println(&piece)
				return true
			}
		}
	}

	return false
}

func choosePromotion(p *Piece) int {
	choosing = true
	queenRect := &sdl.Rect{X: p.Pos.X * screenS / 8, Y: p.Pos.Y * screenS / 8, H: screenS / 10, W: screenS / 10}
	rookRect := &sdl.Rect{X: p.Pos.X * screenS / 8, Y: (p.Pos.Y * screenS / 8) + queenRect.H, H: screenS / 10, W: screenS / 10}
	knightRect := &sdl.Rect{X: p.Pos.X * screenS / 8, Y: (p.Pos.Y * screenS / 8) + (2 * queenRect.H), H: screenS / 10, W: screenS / 10}
	bishopRect := &sdl.Rect{X: p.Pos.X * screenS / 8, Y: (p.Pos.Y * screenS / 8) + (3 * queenRect.H), H: screenS / 10, W: screenS / 10}

	if p.Color == 1 {
		queenPromotionButton = Button{Tex: queenLightPromotionTex, Rect: queenRect}
		rookPromotionButton = Button{Tex: rookLightPromotionTex, Rect: rookRect}
		knightPromotionButton = Button{Tex: knightLightPromotionTex, Rect: knightRect}
		bishopPromotionButton = Button{Tex: bishopLightPromotionTex, Rect: bishopRect}
	} else if p.Color == 0 {
		queenRect.Y -= 3 * queenRect.H
		rookRect.Y -= 3 * queenRect.H
		knightRect.Y -= 3 * queenRect.H
		bishopRect.Y -= 3 * queenRect.H

		queenPromotionButton = Button{Tex: queenDarkPromotionTex, Rect: queenRect}
		rookPromotionButton = Button{Tex: rookDarkPromotionTex, Rect: rookRect}
		knightPromotionButton = Button{Tex: knightDarkPromotionTex, Rect: knightRect}
		bishopPromotionButton = Button{Tex: bishopDarkPromotionTex, Rect: bishopRect}
	}
	drawBoard()
	renderer.Copy(queenPromotionButton.Tex, nil, queenRect)
	renderer.Copy(rookPromotionButton.Tex, nil, rookRect)
	renderer.Copy(knightPromotionButton.Tex, nil, knightRect)
	renderer.Copy(bishopPromotionButton.Tex, nil, bishopRect)
	renderer.Present()
	for choosing {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				{
					running = false
					choosing = false
					return 0
				}
			case *sdl.MouseButtonEvent:
				{
					if t.State == sdl.PRESSED {
						if inRectangle(t.X, t.Y, queenRect) {
							choosing = false
							return 2
						} else if inRectangle(t.X, t.Y, rookRect) {
							choosing = false
							return 3
						} else if inRectangle(t.X, t.Y, knightRect) {
							choosing = false
							return 4
						} else if inRectangle(t.X, t.Y, bishopRect) {
							choosing = false
							return 5
						} else {
							choosing = false
							return 0
						}
					}
				}

			}
		}

	}
	return 0
}

func updateBoard() {
	for i := 0; i < len(board); i++ {
		foundW := false
		foundB := false
		for _, piece := range pieces {
			if piece != toBeDestroyed {
				if piece.Pos == board[i] {
					if piece.Color == 1 {
						foundW = true
						board[i].isOccupiedWhite = true

					} else if piece.Color == 0 {
						foundB = true
						board[i].isOccupiedBlack = true
					}
				}
			}
		}
		if !foundW {
			board[i].isOccupiedWhite = false

		}
		if !foundB {
			board[i].isOccupiedBlack = false
		}
	}
}

func inSquare(x int32, y int32, s *Square) bool {
	if (x >= s.X*screenS/8) && (x <= (s.X*screenS/8)+screenS/8) && (y >= s.Y*screenS/8) && (y <= (s.Y*screenS/8)+screenS/8) {
		return true
	} else {
		return false
	}
}

func destroyPiece(self *Piece) {
	if self.Color == 0 {
		self.Pos.isOccupiedBlack = false
	} else if self.Color == 1 {
		for _, square := range board {
			if square == self.Pos {
				square.isOccupiedWhite = false
			}
		}
	}
	self.Pos = cemetary
	self.Type = -1
}

func loadTexture(file string) *sdl.Texture {
	path := `.\Pieces\` + file
	Surf, err := sdl.LoadBMP(path)
	if err != nil {
		fmt.Println("error loading surface", file, err)
		return nil
	}

	Tex, err := renderer.CreateTextureFromSurface(Surf)
	if err != nil {
		fmt.Println("error loading texture", file, err)
		return nil
	}
	return Tex
}

func inRectangle(pX int32, pY int32, rectangle *sdl.Rect) bool {
	in := false
	if (pX <= (rectangle.X + rectangle.W)) && (pY <= (rectangle.Y + rectangle.H)) && (pX >= rectangle.X) && (pY >= rectangle.Y) {
		in = true
	}
	return in
}

func drawBoard() {
	renderer.Clear()

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			isLight := (rank+file)%2 != 0

			if selectedSquare != nil && (file == int(selectedSquare.X) && rank == int(selectedSquare.Y)) {
				renderer.SetDrawColor(selectedSquareColor.R, selectedSquareColor.G, selectedSquareColor.B, selectedSquareColor.A)
			} else if isLight {
				renderer.SetDrawColor(dark.R, dark.G, dark.B, dark.A)

			} else if !isLight {
				renderer.SetDrawColor(light.R, light.G, light.B, light.A)
			}

			position := sdl.Point{X: int32(file * screenS / 8), Y: int32(rank * screenS / 8)}
			renderer.FillRect(&sdl.Rect{X: position.X, Y: position.Y, W: int32(screenS / 8), H: int32(screenS / 8)})

		}

	}
	renderer.SetDrawColor(oldSquareColor.R, oldSquareColor.G, oldSquareColor.B, oldSquareColor.A)
	renderer.FillRect(lastMove.oldSquare.Rectangle)
	renderer.SetDrawColor(newSquareColor.R, newSquareColor.G, newSquareColor.B, newSquareColor.A)
	renderer.FillRect(lastMove.newSquare.Rectangle)

	for _, piece := range pieces {
		drawPiece(*piece)
	}
	renderer.Present()

}
func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}

	window, err = sdl.CreateWindow(
		"Chess",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenS, screenS,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer:", err)
		return
	}
	defer renderer.Destroy()

	board = initBoard()

	//Loading textures

	plTex = loadTexture("pl.bmp")
	klTex = loadTexture("kl.bmp")
	qlTex = loadTexture("ql.bmp")
	rlTex = loadTexture("rl.bmp")
	nlTex = loadTexture("nl.bmp")
	blTex = loadTexture("bl.bmp")

	pdTex = loadTexture("pd.bmp")
	kdTex = loadTexture("kd.bmp")
	qdTex = loadTexture("qd.bmp")
	rdTex = loadTexture("rd.bmp")
	ndTex = loadTexture("nd.bmp")
	bdTex = loadTexture("bd.bmp")

	queenLightPromotionTex = loadTexture("qlPromotion.bmp")
	rookLightPromotionTex = loadTexture("rlPromotion.bmp")
	knightLightPromotionTex = loadTexture("nlPromotion.bmp")
	bishopLightPromotionTex = loadTexture("blPromotion.bmp")
	queenDarkPromotionTex = loadTexture("qdPromotion.bmp")
	rookDarkPromotionTex = loadTexture("rdPromotion.bmp")
	knightDarkPromotionTex = loadTexture("ndPromotion.bmp")
	bishopDarkPromotionTex = loadTexture("bdPromotion.bmp")

	kw = Piece{Type: 1, Color: 1, Pos: board[39]}

	qw = Piece{Type: 2, Color: 1, Pos: board[31]}

	rw1 = Piece{Type: 3, Color: 1, Pos: board[7]}
	rw2 = Piece{Type: 3, Color: 1, Pos: board[63]}

	nw1 = Piece{Type: 4, Color: 1, Pos: board[15]}
	nw2 = Piece{Type: 4, Color: 1, Pos: board[55]}

	bw1 = Piece{Type: 5, Color: 1, Pos: board[23]}
	bw2 = Piece{Type: 5, Color: 1, Pos: board[47]}

	pw1 = Piece{Type: 6, Color: 1, Pos: board[6]}
	pw2 = Piece{Type: 6, Color: 1, Pos: board[14]}
	pw3 = Piece{Type: 6, Color: 1, Pos: board[22]}
	pw4 = Piece{Type: 6, Color: 1, Pos: board[30]}
	pw5 = Piece{Type: 6, Color: 1, Pos: board[38]}
	pw6 = Piece{Type: 6, Color: 1, Pos: board[46]}
	pw7 = Piece{Type: 6, Color: 1, Pos: board[54]}
	pw8 = Piece{Type: 6, Color: 1, Pos: board[62]}

	kb = Piece{Type: 1, Color: 0, Pos: board[32]}

	qb = Piece{Type: 2, Color: 0, Pos: board[24]}

	rb1 = Piece{Type: 3, Color: 0, Pos: board[0]}
	rb2 = Piece{Type: 3, Color: 0, Pos: board[56]}

	nb1 = Piece{Type: 4, Color: 0, Pos: board[8]}
	nb2 = Piece{Type: 4, Color: 0, Pos: board[48]}

	bb1 = Piece{Type: 5, Color: 0, Pos: board[16]}
	bb2 = Piece{Type: 5, Color: 0, Pos: board[40]}

	pb1 = Piece{Type: 6, Color: 0, Pos: board[1]}
	pb2 = Piece{Type: 6, Color: 0, Pos: board[9]}
	pb3 = Piece{Type: 6, Color: 0, Pos: board[17]}
	pb4 = Piece{Type: 6, Color: 0, Pos: board[25]}
	pb5 = Piece{Type: 6, Color: 0, Pos: board[33]}
	pb6 = Piece{Type: 6, Color: 0, Pos: board[41]}
	pb7 = Piece{Type: 6, Color: 0, Pos: board[49]}
	pb8 = Piece{Type: 6, Color: 0, Pos: board[57]}

	updateBoard()
	renderer.SetDrawColor(255, 0, 0, 255)
	drawBoard()
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return

			case *sdl.MouseButtonEvent:
				if !choosing {

					if t.State == sdl.PRESSED {
						if selectedPiece == nil {
							for i := 0; i < len(board); i++ {
								if inSquare(t.X, t.Y, board[i]) {
									selectedSquare = board[i]
									break
								}
							}
							for i := 0; i < len(pieces); i++ {
								if pieces[i].Pos.X == selectedSquare.X && pieces[i].Pos.Y == selectedSquare.Y {
									if pieces[i].Color == turn {
										selectedPiece = pieces[i]
									}
								}
							}
							drawBoard()
						} else if selectedPiece != nil {
							for i := 0; i < len(board); i++ {
								if inSquare(t.X, t.Y, board[i]) {
									selectedSquare = board[i]
								}
							}
							newPieceClicked := false
							for i := 0; i < len(pieces); i++ {
								if pieces[i].Pos == selectedSquare && pieces[i].Color == selectedPiece.Color {
									selectedPiece = pieces[i]
									newPieceClicked = true
									break
								}
							}
							if newPieceClicked {
								drawBoard()
								break
							}
							if (selectedPiece == &Piece{}) {
								return
							}
							posSave = selectedPiece.Pos
							err = move(selectedPiece, selectedSquare, false)
							if err != nil {
								move(selectedPiece, posSave, false)
							}
							selectedPiece = nil
							selectedSquare = nil
							updateBoard()
							drawBoard()
							break
						}
					}
				}
			}

		}
	}
}
