# Render SD

![License](https://img.shields.io/badge/license-MIT-blue.svg)

**Render SD** to narzędzie wiersza poleceń do renderowania opisów produktów w formacie HTML na podstawie danych JSON. Projekt wykorzystuje bibliotekę [templ](https://github.com/a-h/templ) do generowania szablonów HTML w języku Go.

## 🚀 Szybki start

### Wymagania

- Go 1.24.4 lub nowszy

### Instalacja

```bash
# Klonowanie repozytorium
git clone https://github.com/marad/render-allegro-description.git
cd render-sd

# Instalacja zależności
go mod tidy

# Kompilacja szablonów templ
templ generate

# Budowanie aplikacji
go build -o render-sd
```

### Podstawowe użycie

```bash
# Renderowanie z pliku JSON
./render-sd --file test.json

# Renderowanie ze standardowego wejścia
cat test.json | ./render-sd

# Renderowanie tylko zawartości (bez pełnej strony HTML)
./render-sd --file test.json --no-page
```

## 📋 Składnia

```
render-sd [--file=<path> | -f <path>] [--no-page]
render-sd -h | --help
```

### Opcje

| Opcja | Opis |
|-------|------|
| `-f <path>`, `--file=<path>` | Ścieżka do pliku JSON (jeśli nie podano, czyta ze stdin) |
| `--no-page` | Renderuje tylko zawartość opisu bez pełnej strony HTML |
| `-h`, `--help` | Pokazuje pomoc |

## 📊 Format danych JSON

Narzędzie oczekuje danych JSON w formacie Allegro Standardized Description:

```json
{
  "sections": [
    {
      "items": [
        {
          "type": "TEXT",
          "content": "<h1>Tytuł produktu</h1><p>Opis produktu</p>"
        },
        {
          "type": "IMAGE",
          "url": "https://example.com/product.jpg"
        }
      ]
    }
  ]
}
```

### Obsługiwane typy elementów

- **TEXT** - zawartość HTML w postaci tekstu
- **IMAGE** - obraz z podanym URL

## 🧪 Testowanie

Projekt zawiera kompleksowy skrypt testowy:

```bash
# Uruchomienie wszystkich testów
./test.sh
```

## 🛠️ Rozwój

### Wymagania dla deweloperów

- Go 1.24.4+
- [templ](https://github.com/a-h/templ) CLI tool

### Instalacja narzędzi deweloperskich

```bash
# Instalacja templ CLI
go install github.com/a-h/templ/cmd/templ@latest
```

### Regeneracja szablonów

Po modyfikacji plików `.templ`:

```bash
templ generate
```

### Dodawanie nowych typów elementów

1. Rozszerz strukturę `Item` w `description.go`
2. Dodaj obsługę w szablonie `description.templ`
3. Regeneruj szablony: `templ generate`
4. Dodaj testy w `test.sh`

## 📄 Przykład użycia

```bash
# Tworzenie przykładowego pliku JSON
cat > example.json << EOF
{
  "sections": [
    {
      "items": [
        {
          "type": "TEXT",
          "content": "<h1>Fantastyczny Produkt</h1><p>To jest najlepszy produkt na rynku!</p>"
        },
        {
          "type": "IMAGE",
          "url": "https://via.placeholder.com/400x300"
        },
        {
          "type": "TEXT",
          "content": "<h2>Specyfikacja</h2><ul><li>Wysoka jakość</li><li>Trwały materiał</li></ul>"
        }
      ]
    }
  ]
}
EOF

# Renderowanie do pliku HTML
./render-sd --file example.json > output.html

# Otwarcie w przeglądarce (macOS)
open output.html
```

## 📝 Licencja

Ten projekt jest licencjonowany na licencji MIT - szczegóły w pliku [LICENSE](LICENSE).

## 🐛 Zgłaszanie błędów

Jeśli znajdziesz błąd, proszę [otwórz issue](https://github.com/marad/render-allegro-description/issues) z:

- Opisem problemu
- Krokami do odtworzenia
- Oczekiwanym zachowaniem
- Środowiskiem (OS, wersja Go)

---

⭐ **Jeśli ten projekt był pomocny, zostaw gwiazdkę!**
