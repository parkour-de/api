package verband

import (
	"context"
	"pkv/api/src/domain/verband"
)

type Bundesland struct {
	Vereine    int `json:"vereine"`
	Mitglieder int `json:"mitglieder"`
}

func (s *Service) BundeslandInfo(ctx context.Context) (map[string]Bundesland, error) {
	vereine, mitglieder, err := s.VereineByBundesland(ctx)
	if err != nil {
		return nil, err
	}
	bundeslaender := make(map[string]Bundesland)
	for bundesland := range vereine {
		bundeslaender[bundesland] = Bundesland{
			Vereine:    vereine[bundesland],
			Mitglieder: mitglieder[bundesland],
		}
	}

	return bundeslaender, nil
}

func (s *Service) VereineByBundesland(ctx context.Context) (map[string]int, map[string]int, error) {
	vereine, vereineDetail, err := s.GetVereine(ctx)
	if err != nil {
		return nil, nil, err
	}

	aggregation := aggregateVereineByBundesland(vereine)
	mitglieder := aggregateMitgliederByBundesland(vereineDetail)

	return aggregation, mitglieder, nil
}

func aggregateVereineByBundesland(vereine []verband.Verein) map[string]int {
	aggregation := make(map[string]int)
	for _, verein := range vereine {
		aggregation[verein.Bundesland]++
	}
	return aggregation
}

func aggregateMitgliederByBundesland(vereine []VereinDetail) map[string]int {
	aggregation := make(map[string]int)
	for _, verein := range vereine {
		aggregation[verein.Bundesland] += verein.Mitglieder
	}
	return aggregation
}
