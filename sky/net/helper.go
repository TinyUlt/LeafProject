package net
import (
	"time"
)
func push(ch chan uint32, v uint32) {
    for {
        select {
        case ch <- v:
            return
        default:
        }
        select {
        case <- ch:
        default:
            //return
        }
    }
}
func pushex(ch chan uint32, v uint32) {
	for {
		c := time.After(time.Duration(200) * time.Millisecond)

		select {
		case ch <- v:
			return
		case <-c:
			//return
		}
		select {
		case <- ch:
		default:
			//return
		}
	}
}

func pop(ch chan uint32) {
    for {
        select {
        case <- ch:
        default:
            return
        }
    }
}
