package reductor

func (rdc *Reductor) ChanIn() chan Message {
	return rdc.in
}
