package token

import (
	"testing"
	"time"
)

func TestTrack_And_Count(t *testing.T) {
	s := NewStore()
	tok := Token{Value: "tok1", ExpiresAt: time.Now().Add(time.Minute)}
	s.Track(tok)
	if s.Count() != 1 {
		t.Errorf("Count = %d; want 1", s.Count())
	}
}

func TestRevoke_MarksToken(t *testing.T) {
	s := NewStore()
	s.Track(Token{Value: "tok2", ExpiresAt: time.Now().Add(time.Minute)})
	s.Revoke("tok2")
	if !s.IsRevoked("tok2") {
		t.Error("expected tok2 to be revoked")
	}
}

func TestIsRevoked_UnknownToken_ReturnsFalse(t *testing.T) {
	s := NewStore()
	if s.IsRevoked("unknown") {
		t.Error("expected unknown token to not be revoked")
	}
}

func TestPurge_RemovesExpiredTokens(t *testing.T) {
	s := NewStore()
	past := time.Now().Add(-time.Second)
	future := time.Now().Add(time.Minute)
	s.Track(Token{Value: "expired", ExpiresAt: past})
	s.Track(Token{Value: "valid", ExpiresAt: future})
	removed := s.Purge(time.Now())
	if removed != 1 {
		t.Errorf("Purge removed %d; want 1", removed)
	}
	if s.Count() != 1 {
		t.Errorf("Count = %d; want 1 after purge", s.Count())
	}
}

func TestPurge_AlsoRemovesRevokedExpired(t *testing.T) {
	s := NewStore()
	past := time.Now().Add(-time.Second)
	s.Track(Token{Value: "old-revoked", ExpiresAt: past})
	s.Revoke("old-revoked")
	s.Purge(time.Now())
	// After purge the revocation entry should also be cleaned up
	if s.IsRevoked("old-revoked") {
		t.Error("expected revoked+expired token to be purged from revocation list")
	}
}

func TestPurge_EmptyStore_ReturnsZero(t *testing.T) {
	s := NewStore()
	if n := s.Purge(time.Now()); n != 0 {
		t.Errorf("Purge on empty store = %d; want 0", n)
	}
}
