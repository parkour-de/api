package verband

import (
	"context"
	"pkv/api/src/domain/verband"
)

type Bundesland struct {
	Vereine int `json:"vereine"`
}

func (s *Service) BundeslandInfo(ctx context.Context) (map[string]Bundesland, error) {
	vereine, err := s.VereineByBundesland(ctx)
	if err != nil {
		return nil, err
	}
	bundeslaender := make(map[string]Bundesland)
	for bundesland := range vereine {
		bundeslaender[bundesland] = Bundesland{
			Vereine: vereine[bundesland],
		}
	}

	return bundeslaender, nil
}

func (s *Service) VereineByBundesland(ctx context.Context) (map[string]int, error) {
	vereine, err := s.GetVereine(ctx)
	if err != nil {
		return nil, err
	}

	aggregation := aggregateVereineByBundesland(vereine)

	return aggregation, nil
}

func aggregateVereineByBundesland(vereine []verband.Verein) map[string]int {
	aggregation := make(map[string]int)
	for _, verein := range vereine {
		aggregation[verein.Bundesland]++
	}
	return aggregation
}
