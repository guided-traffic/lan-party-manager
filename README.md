# LAN Party Manager

Eine Webanwendung zur Verwaltung von LAN-Party-Events.

## 🚀 Technologie-Stack

- **Frontend**: Angular 19+ mit TypeScript und SCSS
- **Backend**: Go 1.22+ mit Gin Framework
- **Deployment**: Kubernetes mit Helm Charts

## 📁 Projektstruktur

```
lan-party-manager/
├── frontend/                 # Angular Frontend
├── backend/                  # Go Backend
├── helm/                     # Helm Charts
│   └── lan-party-manager/
└── .github/                  # GitHub Konfiguration
```

## 🛠️ Entwicklung

### Voraussetzungen

- Node.js 20+
- Go 1.22+
- Docker (optional)
- Kubernetes & Helm (für Deployment)

### Frontend starten

```bash
cd frontend
npm install
npm start
```

Das Frontend ist unter http://localhost:4200 erreichbar.

### Backend starten

```bash
cd backend
go mod tidy
go run main.go
```

Das Backend ist unter http://localhost:8080 erreichbar.

## 🐳 Docker

### Images bauen

```bash
# Frontend
docker build -t lan-party-manager/frontend:latest ./frontend

# Backend
docker build -t lan-party-manager/backend:latest ./backend
```

## ☸️ Kubernetes Deployment

### Mit Helm installieren

```bash
helm install lan-party-manager ./helm/lan-party-manager
```

### Mit custom Values

```bash
helm install lan-party-manager ./helm/lan-party-manager -f custom-values.yaml
```

## 📡 API Endpoints

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| GET | `/health` | Health Check |
| GET | `/api/v1/events` | Alle Events abrufen |
| GET | `/api/v1/events/:id` | Einzelnes Event |
| POST | `/api/v1/events` | Event erstellen |
| PUT | `/api/v1/events/:id` | Event aktualisieren |
| DELETE | `/api/v1/events/:id` | Event löschen |
| GET | `/api/v1/participants` | Alle Teilnehmer |
| POST | `/api/v1/participants` | Teilnehmer erstellen |

## 🎨 Credits

Achievement-Icons von [Game-icons.net](https://game-icons.net) unter [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/) Lizenz.

## 📄 Lizenz

MIT
