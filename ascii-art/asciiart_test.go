package asciiart

import (
	"os"
	"testing"
)

// ---------- Tests pour splitcara ----------

func TestSplitcara(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]rune
		wantErr bool
	}{
		{
			name:  "chaine vide",
			input: "",
			want:  nil,
		},
		{
			name:  "mot simple",
			input: "abc",
			want:  [][]rune{{'a', 'b', 'c'}},
		},
		{
			name:  "deux lignes separees par \\n",
			input: `a\nb`,
			want:  [][]rune{{'a'}, {'b'}},
		},
		{
			name:  "ligne vide au milieu",
			input: `a\n\nb`,
			want:  [][]rune{{'a'}, {}, {'b'}},
		},
		{
			name:    "caractere de controle interdit",
			input:   "a\x01b",
			wantErr: true,
		},
		{
			name:    "caractere hors plage ascii imprimable",
			input:   "café", // 'é' > 126
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitcara(tt.input)

			if tt.wantErr {
				if err == nil || err.Error() != "400" {
					t.Fatalf("attendu une erreur 400, obtenu: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("erreur inattendue: %v", err)
			}
			if !equalRuneSlices(got, tt.want) {
				t.Fatalf("splitcara(%q) = %v, attendu %v", tt.input, got, tt.want)
			}
		})
	}
}

func equalRuneSlices(a, b [][]rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// ---------- Tests pour changeToAsciiArt ----------

func TestChangeToAsciiArt(t *testing.T) {
	// On construit un faux fichier de banniere ne couvrant que
	// les caracteres ' ' (index 0) et '!' (index 9), 9 lignes chacun.
	fakeFile := []string{
		"L0", "L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8", // bloc pour ' '
		"L9", "L10", "L11", "L12", "L13", "L14", "L15", "L16", "L17", // bloc pour '!'
	}

	t.Run("un seul mot", func(t *testing.T) {
		cara := [][]rune{{' ', '!'}}
		got := changeToAsciiArt(fakeFile, cara)

		want := "L0L9" +
			"L1L10\n" +
			"L2L11\n" +
			"L3L12\n" +
			"L4L13\n" +
			"L5L14\n" +
			"L6L15\n" +
			"L7L16\n" +
			"L8L17\n"

		if got != want {
			t.Fatalf("resultat inattendu:\ngot:  %q\nwant: %q", got, want)
		}
	})

	t.Run("ligne vide produit un simple saut de ligne", func(t *testing.T) {
		cara := [][]rune{{}}
		got := changeToAsciiArt(fakeFile, cara)
		if got != "\n" {
			t.Fatalf("resultat inattendu: got %q, attendu %q", got, "\n")
		}
	})
}

// ---------- Tests pour whoPolis ----------

func TestWhoPolis(t *testing.T) {
	t.Run("type inconnu renvoie 400", func(t *testing.T) {
		_, err := whoPolis("inexistant")
		if err == nil || err.Error() != "400" {
			t.Fatalf("attendu erreur 400, obtenu: %v", err)
		}
	})

	t.Run("type connu mais fichier absent renvoie 500", func(t *testing.T) {
		const key = "test_missing_file"
		banners[key] = "fichier_qui_nexiste_pas_xyz.txt"
		defer delete(banners, key)

		_, err := whoPolis(key)
		if err == nil || err.Error() != "500" {
			t.Fatalf("attendu erreur 500, obtenu: %v", err)
		}
	})

	t.Run("type connu avec fichier present renvoie le contenu", func(t *testing.T) {
		const key = "test_existing_file"
		const path = "test_banner_temp.txt"
		content := []byte("contenu de test")

		if err := os.WriteFile(path, content, 0644); err != nil {
			t.Fatalf("impossible de creer le fichier temporaire: %v", err)
		}
		defer os.Remove(path)

		banners[key] = path
		defer delete(banners, key)

		got, err := whoPolis(key)
		if err != nil {
			t.Fatalf("erreur inattendue: %v", err)
		}
		if string(got) != string(content) {
			t.Fatalf("contenu inattendu: got %q, attendu %q", got, content)
		}
	})
}

// ---------- Tests d'integration pour AsciiArt ----------

func TestAsciiArt(t *testing.T) {
	t.Run("type invalide renvoie 400", func(t *testing.T) {
		_, err := AsciiArt("hello", "type_qui_nexiste_pas")
		if err == nil || err.Error() != "400" {
			t.Fatalf("attendu erreur 400, obtenu: %v", err)
		}
	})

	t.Run("caractere invalide dans la chaine renvoie 400", func(t *testing.T) {
		// On cree quand meme une banniere valide pour ne tester
		// que la validation de la chaine d'entree.
		const key = "test_ascii_art_ok"
		const path = "test_banner_ascii_art.txt"

		// Bloc minimal pour ' ' (9 lignes), suffisant car l'erreur
		// doit survenir avant tout acces au fichier de police.
		content := "L0\nL1\nL2\nL3\nL4\nL5\nL6\nL7\nL8\n"
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("impossible de creer le fichier temporaire: %v", err)
		}
		defer os.Remove(path)

		banners[key] = path
		defer delete(banners, key)

		_, err := AsciiArt("caf\x01e", key)
		if err == nil || err.Error() != "400" {
			t.Fatalf("attendu erreur 400, obtenu: %v", err)
		}
	})

	t.Run("cas nominal produit l'ascii art attendu", func(t *testing.T) {
		const key = "test_ascii_art_nominal"
		const path = "test_banner_nominal.txt"

		// Bloc pour ' ' (index 0..8) puis '!' (index 9..17).
		lines := []string{
			"L0", "L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8",
			"L9", "L10", "L11", "L12", "L13", "L14", "L15", "L16", "L17",
		}
		content := ""
		for _, l := range lines {
			content += l + "\n"
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("impossible de creer le fichier temporaire: %v", err)
		}
		defer os.Remove(path)

		banners[key] = path
		defer delete(banners, key)

		got, err := AsciiArt(" !", key)
		if err != nil {
			t.Fatalf("erreur inattendue: %v", err)
		}

		want := "L0L9" +
			"L1L10\n" +
			"L2L11\n" +
			"L3L12\n" +
			"L4L13\n" +
			"L5L14\n" +
			"L6L15\n" +
			"L7L16\n" +
			"L8L17\n"

		if got != want {
			t.Fatalf("resultat inattendu:\ngot:  %q\nwant: %q", got, want)
		}
	})
}
