package znakdb

import (
	"agregat/domain"
	"strings"
)

func (z *DbZnak) Guide() (guide map[string]*domain.ProductGuides, err error) {
	sess := z.dbSession
	guide = make(map[string]*domain.ProductGuides)
	recs := make([]*domain.ProductGuides, 0)
	res := sess.Collection("product_guides struct").Find("product_gtin <> ?", "")
	if err := res.All(&recs); err != nil {
		return nil, err
	}
	for _, g := range recs {
		guide[strings.TrimSpace(g.ProductGtin)] = g
	}
	return guide, err
}
