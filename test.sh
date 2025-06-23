#!/bin/bash

# Skrypt testowy dla narzędzia render-sd
# Użycie: ./test.sh
# 
# Skrypt testuje wszystkie główne funkcjonalności narzędzia render-sd:
# - Wczytywanie JSON ze stdin i z pliku
# - Walidację JSON i struktury danych
# - Obsługę błędów i kodów wyjścia
# - Generowanie poprawnego HTML
#
# W przypadku błędu test wyświetla szczegółowe informacje o przyczynie

# Kolory dla lepszej czytelności
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Liczniki testów
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Funkcja do wyświetlania nagłówka testu
test_header() {
    echo -e "\n${BLUE}=== TEST $1: $2 ===${NC}"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

# Funkcja do sprawdzania wyniku testu
check_result() {
    local expected_exit_code=$1
    local actual_exit_code=$2
    local test_name="$3"
    local expected_output="$4"
    local actual_output="$5"
    
    if [ "$expected_exit_code" -eq "$actual_exit_code" ]; then
        if [ -n "$expected_output" ] && [[ "$actual_output" != *"$expected_output"* ]]; then
            echo -e "${RED}❌ FAIL: $test_name${NC}"
            echo -e "${YELLOW}Oczekiwany fragment wyjścia:${NC} $expected_output"
            echo -e "${YELLOW}Rzeczywiste wyjście:${NC} $actual_output"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        else
            echo -e "${GREEN}✅ PASS: $test_name${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        fi
    else
        echo -e "${RED}❌ FAIL: $test_name${NC}"
        echo -e "${YELLOW}Oczekiwany kod wyjścia:${NC} $expected_exit_code"
        echo -e "${YELLOW}Rzeczywisty kod wyjścia:${NC} $actual_exit_code"
        if [ -n "$actual_output" ]; then
            echo -e "${YELLOW}Wyjście:${NC} $actual_output"
        fi
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Funkcja do uruchamiania testu
run_test() {
    local command="$1"
    local expected_exit_code=$2
    local test_name="$3"
    local expected_output="$4"
    
    # Uruchomienie polecenia
    local output
    output=$(bash -c "$command" 2>&1)
    local exit_code=$?
    
    check_result "$expected_exit_code" "$exit_code" "$test_name" "$expected_output" "$output"
}

echo -e "${BLUE}🧪 Rozpoczynanie testów narzędzia render-sd${NC}"

# Kompilacja aplikacji
echo -e "\n${BLUE}📦 Kompilacja aplikacji...${NC}"
if ! go build -o render-sd; then
    echo -e "${RED}❌ Błąd kompilacji aplikacji${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Aplikacja skompilowana pomyślnie${NC}"

# TEST 1: Podstawowy JSON ze stdin
test_header "1" "Podstawowy JSON ze stdin"
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"TEXT\",\"content\":\"<p>Test content</p>\"}]}]}' | ./render-sd" 0 "podstawowy JSON" "<!doctype html>"

# TEST 2: JSON z obrazem
test_header "2" "JSON z elementem IMAGE"
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"IMAGE\",\"url\":\"https://example.com/image.jpg\"}]}]}' | ./render-sd" 0 "element IMAGE" "<img src=\"https://example.com/image.jpg\""

# TEST 3: Tworzenie pliku testowego i test z flagą --file
test_header "3" "Test z flagą --file"
cat > test_file.json << 'EOF'
{
  "sections": [
    {
      "items": [
        {
          "type": "TEXT",
          "content": "<h1>Tytuł z pliku</h1>"
        }
      ]
    }
  ]
}
EOF

run_test "./render-sd --file test_file.json" 0 "wczytywanie z pliku --file" "<h1>Tytuł z pliku</h1>"

# TEST 4: Test ze skróconą flagą -f
test_header "4" "Test ze skróconą flagą -f"
run_test "./render-sd -f test_file.json" 0 "wczytywanie z pliku -f" "<h1>Tytuł z pliku</h1>"

# TEST 5: Pusty JSON
test_header "5" "Pusty obiekt JSON"
run_test "echo '{}' | ./render-sd" 0 "pusty JSON" "<!doctype html>"

# TEST 6: Puste sections
test_header "6" "Puste sections"
run_test "echo '{\"sections\":[]}' | ./render-sd" 0 "puste sections" "<main><section></section></main>"

# TEST 7: Nieistniejący plik
test_header "7" "Nieistniejący plik"
run_test "./render-sd --file nieistniejacy_plik.json" 1 "nieistniejący plik" "Error reading file"

# TEST 8: Nieprawidłowy JSON
test_header "8" "Nieprawidłowy JSON"
run_test "echo '{\"sections\":invalid}' | ./render-sd" 1 "nieprawidłowy JSON" "error parsing JSON"

# TEST 9: Pusty string
test_header "9" "Pusty string ze stdin"
run_test "echo -n '' | ./render-sd" 1 "pusty string" "Error reading description: error parsing JSON: no JSON data provided"

# TEST 10: Element TEXT bez zawartości
test_header "10" "Element TEXT bez zawartości"
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"TEXT\",\"content\":\"\"}]}]}' | ./render-sd" 1 "TEXT bez zawartości" "Error reading description: error parsing JSON: TEXT item in section 1 has no content"

# TEST 11: Element IMAGE bez URL
test_header "11" "Element IMAGE bez URL"
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"IMAGE\",\"url\":\"\"}]}]}' | ./render-sd" 1 "IMAGE bez URL" "Error reading description: error parsing JSON: IMAGE item in section 1 has no URL"

# TEST 12: Element bez typu
test_header "12" "Element bez typu"
run_test "echo '{\"sections\":[{\"items\":[{\"content\":\"test\"}]}]}' | ./render-sd" 1 "element bez typu" "Error reading description: error parsing JSON: item in section 1 has no type specified"

# TEST 13: Nieznany typ elementu (powinien działać z ostrzeżeniem)
test_header "13" "Nieznany typ elementu"
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"UNKNOWN\",\"content\":\"test\"}]}]}' | ./render-sd" 0 "nieznany typ" "Warning: unknown item type 'UNKNOWN' in section 1"

# TEST 14: Wielosekcyjny JSON z błędem
test_header "14" "Wielosekcyjny JSON z błędem w drugiej sekcji"
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"TEXT\",\"content\":\"OK\"}]},{\"items\":[{\"type\":\"IMAGE\",\"url\":\"\"}]}]}' | ./render-sd" 1 "błąd w drugiej sekcji" "Error reading description: error parsing JSON: IMAGE item in section 2 has no URL"

# TEST 15: Pomoc aplikacji
test_header "15" "Test pomocy --help"
run_test "./render-sd --help" 0 "pomoc aplikacji" "Usage:"

# TEST 16: Złożony poprawny JSON
test_header "16" "Złożony poprawny JSON"
cat > complex_test.json << 'EOF'
{
  "sections": [
    {
      "items": [
        {
          "type": "TEXT",
          "content": "<h1>Tytuł</h1><p>Opis</p>"
        },
        {
          "type": "IMAGE",
          "url": "https://example.com/img1.jpg"
        }
      ]
    },
    {
      "items": [
        {
          "type": "TEXT",
          "content": "<p>Druga sekcja</p>"
        }
      ]
    }
  ]
}
EOF

run_test "./render-sd -f complex_test.json" 0 "złożony JSON" "<div class=\"row\">"

# TEST 17: Test flagi --no-page ze stdin
test_header "17" "Test flagi --no-page ze stdin"
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"TEXT\",\"content\":\"<p>Test no-page</p>\"}]}]}' | ./render-sd --no-page" 0 "renderowanie bez pełnej strony" "<div class=\"row\"><div class=\"item\"><p>Test no-page</p></div></div>"

# TEST 18: Test flagi --no-page z plikiem
test_header "18" "Test flagi --no-page z plikiem"
cat > no_page_test.json << 'EOF'
{
  "sections": [
    {
      "items": [
        {
          "type": "TEXT",
          "content": "<h2>Bez pełnej strony</h2>"
        }
      ]
    }
  ]
}
EOF

run_test "./render-sd --no-page -f no_page_test.json" 0 "no-page z pliku" "<h2>Bez pełnej strony</h2>"

# TEST 19: Porównanie z pełną stroną vs --no-page
test_header "19" "Porównanie: pełna strona vs --no-page"
# Test że bez --no-page zawiera <!doctype html>
run_test "echo '{\"sections\":[{\"items\":[{\"type\":\"TEXT\",\"content\":\"<p>Test</p>\"}]}]}' | ./render-sd" 0 "pełna strona zawiera doctype" "<!doctype html>"

# TEST 20: Tryb --build-sections z płaską listą itemów
test_header "20" "Tryb --build-sections z płaską listą itemów"
cat > items_test.json << 'EOF'
[
  { "type": "TEXT", "content": "<h1>Tytuł testowy</h1><p>Opis testowy</p>" },
  { "type": "IMAGE", "url": "https://example.com/test.jpg" },
  { "type": "TEXT", "content": "<p>Kolejna sekcja tekstowa</p>" },
  { "type": "TEXT", "content": "<p>Jeszcze jeden tekst</p>" }
]
EOF

run_test "./render-sd --file items_test.json --build-sections" 0 "build-sections z płaską listą" "<h1>Tytuł testowy</h1>"
run_test "./render-sd --file items_test.json --build-sections" 0 "build-sections - sekcje" "<div class=\"row\">"

# Czyszczenie plików testowych
rm -f test_file.json complex_test.json no_page_test.json no_page_test.json items_test.json

# Podsumowanie wyników
echo -e "\n${BLUE}📊 PODSUMOWANIE TESTÓW${NC}"
echo -e "${BLUE}===================${NC}"
echo -e "Łącznie testów: $TOTAL_TESTS"
echo -e "${GREEN}Zakończone sukcesem: $PASSED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 Wszystkie testy przeszły pomyślnie!${NC}"
    exit 0
else
    echo -e "${RED}Zakończone błędem: $FAILED_TESTS${NC}"
    echo -e "\n${RED}💥 Niektóre testy nie powiodły się!${NC}"
    exit 1
fi
