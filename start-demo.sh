#!/bin/bash

# LAN Party Manager - Demo Start Script
# Startet Backend und Frontend für lokale Entwicklung

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"

# Farben für Output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     🎮 LAN Party Manager - Demo Mode 🎮    ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════╝${NC}"
echo ""

# Funktion zum Aufräumen bei Beendigung
cleanup() {
    echo ""
    echo -e "${YELLOW}Beende alle Prozesse...${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    echo -e "${GREEN}Auf Wiedersehen! 👋${NC}"
    exit 0
}

trap cleanup SIGINT SIGTERM

# Prüfe ob .env existiert
if [ ! -f "$BACKEND_DIR/.env" ]; then
    echo -e "${RED}❌ Fehler: Backend .env Datei nicht gefunden!${NC}"
    echo -e "${YELLOW}Bitte erstelle $BACKEND_DIR/.env mit den erforderlichen Umgebungsvariablen.${NC}"
    echo ""
    echo "Beispiel:"
    echo "  STEAM_API_KEY=dein-steam-api-key"
    echo "  JWT_SECRET=dein-jwt-secret"
    echo "  FRONTEND_URL=http://localhost:4200"
    echo "  BACKEND_URL=http://localhost:8080"
    exit 1
fi

# Prüfe ob Go installiert ist
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Fehler: Go ist nicht installiert!${NC}"
    exit 1
fi

# Prüfe ob Node/npm installiert ist
if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ Fehler: npm ist nicht installiert!${NC}"
    exit 1
fi

# Backend starten
echo -e "${BLUE}🚀 Starte Backend...${NC}"
cd "$BACKEND_DIR"

# Go dependencies laden falls nötig
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    echo -e "${YELLOW}   Lade Go Dependencies...${NC}"
    go mod tidy
fi

# Backend im Hintergrund starten
go run main.go &
BACKEND_PID=$!
echo -e "${GREEN}   ✓ Backend gestartet (PID: $BACKEND_PID)${NC}"

# Kurz warten damit Backend hochfahren kann
sleep 2

# Frontend starten
echo -e "${BLUE}🚀 Starte Frontend...${NC}"
cd "$FRONTEND_DIR"

# NPM dependencies installieren falls nötig
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}   Installiere npm Dependencies...${NC}"
    npm install
fi

# Frontend im Hintergrund starten
npm start &
FRONTEND_PID=$!
echo -e "${GREEN}   ✓ Frontend gestartet (PID: $FRONTEND_PID)${NC}"

echo ""
echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ LAN Party Manager läuft!${NC}"
echo ""
echo -e "   🌐 Frontend: ${BLUE}http://localhost:4200${NC}"
echo -e "   🔧 Backend:  ${BLUE}http://localhost:8080${NC}"
echo -e "   📊 Health:   ${BLUE}http://localhost:8080/health${NC}"
echo ""
echo -e "${YELLOW}Drücke Ctrl+C zum Beenden${NC}"
echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo ""

# Warte auf beide Prozesse
wait $BACKEND_PID $FRONTEND_PID
