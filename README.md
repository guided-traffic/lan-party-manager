# LAN Party Manager

Eine Webanwendung für LAN-Partys, bei der sich Spieler gegenseitig mit Achievements bewerten können.

## ✨ Features

- 🎮 **Steam Login** - Authentifizierung über Steam OpenID
- 💰 **Credit System** - Spieler erhalten automatisch Credits über Zeit
- 🏆 **Achievement Voting** - Spieler bewerten sich gegenseitig mit vordefinierten Achievements
- 📺 **Live Timeline** - Alle Votes in Echtzeit via WebSocket
- 🥇 **Leaderboard** - Top 3 pro Achievement

## 🚀 Installation

### Voraussetzungen

- Kubernetes Cluster
- Helm 3.x
- Steam Web API Key ([hier beantragen](https://steamcommunity.com/dev/apikey))

### Helm Repository hinzufügen

```bash
helm repo add lan-party-manager https://guided-traffic.github.io/lan-party-manager
helm repo update
```

### Installation

```bash
helm install lan-party-manager lan-party-manager/lan-party-manager -f values.yaml
```

## ⚙️ Konfiguration

| Parameter | Beschreibung | Default |
|-----------|--------------|---------|
| `secrets.steamApiKey` | Steam Web API Key (erforderlich) | `""` |
| `secrets.jwtSecret` | JWT Secret für Token-Signierung (erforderlich) | `""` |
| `backend.env.CREDIT_INTERVAL_MINUTES` | Minuten zwischen Credit-Vergabe | `10` |
| `backend.env.CREDIT_MAX` | Maximale Credits pro Spieler | `10` |
| `backend.env.JWT_EXPIRATION_DAYS` | JWT Gültigkeit in Tagen | `7` |
| `ingress.enabled` | Ingress aktivieren | `false` |
| `ingress.hosts` | Ingress Hosts Konfiguration | `[]` |

Alle verfügbaren Optionen findest du in der [values.yaml](helm/lan-party-manager/values.yaml).

## 🎨 Credits

Achievement-Icons von [Game-icons.net](https://game-icons.net) unter [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/) Lizenz.

## 📄 Lizenz

Apache 2.0
