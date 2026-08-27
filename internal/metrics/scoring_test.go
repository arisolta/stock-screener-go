package metrics

import (
	"math"
	"testing"
)

func TestScoreHigher(t *testing.T) {
	// Full credit threshold test
	s := ScoreHigher(0.25, 17, 0.25, 0.0, ptr(0.50))
	expected := 17.0 * 0.85
	if math.Abs(s-expected) > 1e-4 {
		t.Errorf("expected %.4f, got %.4f", expected, s)
	}

	// Exceptional credit test
	sExp := ScoreHigher(0.50, 17, 0.25, 0.0, ptr(0.50))
	if math.Abs(sExp-17.0) > 1e-4 {
		t.Errorf("expected 17.0, got %.4f", sExp)
	}

	// Zero credit test
	sZero := ScoreHigher(0.0, 17, 0.25, 0.0, ptr(0.50))
	if sZero != 0.0 {
		t.Errorf("expected 0.0, got %.4f", sZero)
	}
}

func TestScoreLower(t *testing.T) {
	// Best at test (1.05 earned 85% of 18)
	s := ScoreLower(1.05, 18, 1.05, 2.0, ptr(1.0))
	expected := 18.0 * 0.85
	if math.Abs(s-expected) > 1e-4 {
		t.Errorf("expected %.4f, got %.4f", expected, s)
	}

	// Exceptional at test (1.0 earned full 18)
	sExp := ScoreLower(1.0, 18, 1.05, 2.0, ptr(1.0))
	if math.Abs(sExp-18.0) > 1e-4 {
		t.Errorf("expected 18.0, got %.4f", sExp)
	}

	// Zero credit test
	sZero := ScoreLower(2.0, 18, 1.05, 2.0, ptr(1.0))
	if sZero != 0.0 {
		t.Errorf("expected 0.0, got %.4f", sZero)
	}
}

func TestNormalizeBucketScore(t *testing.T) {
	norm := NormalizeBucketScore(21.25, 25)
	if norm != 85.0 {
		t.Errorf("expected 85.0, got %.1f", norm)
	}
}
