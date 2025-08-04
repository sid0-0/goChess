This is temporary README. Will update it soon!

---

# 🧠 GoChess

A real-time chess game built with **Go (Golang)** on the backend and **HTMX** on the frontend. This project demonstrates how far you can push interactivity in the browser using HTMX, even for something as dynamic as a multiplayer game.

---

## ⚙️ Tech Stack

- **Backend**: Golang (`net/http`, `chi`, `html/template`)
- **Frontend**: HTMX, HTML templates
- **Realtime Communication**: WebSockets (for game state)
- **API**: RESTful endpoints for non-game interactions

---

## 🚀 Features

- ♟️ Full chess game logic
- 👥 2-player support (live games)
- 👁️ Spectator mode
- 📡 WebSocket-powered game communication
- 📮 HTMX-powered API for session, view updates, and more

---

## 🔧 Getting Started

To run the project locally:

1. Clone the repo:
   ```bash
   git clone https://github.com/yourusername/htmx-chess.git
   cd htmx-chess
   ```

2. Build and run the server (typical Go workflow):
   ```bash
   go run main.go
   ```

3. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

---

## 🛣️ Roadmap

Planned features:

- ⏳ Replay timers (per-move timing)
- 🏳️ Resign functionality
- 🔄 Start new game after game end  
  *(Currently requires manual cookie/session clearing)*
- 🚪 Spectator pool leaving
- 💬 In-game chat

---

## 🎯 Purpose

This project was created as a **proof of concept** to explore the extent of what’s possible using [HTMX](https://htmx.org/) for building **highly interactive, real-time** browser applications — traditionally a domain reserved for full SPA frameworks.

---

## 📬 Feedback & Contributions

Feel free to open issues or PRs if you'd like to contribute or suggest improvements!

---
```

Let me know if you'd like a section on folder structure, usage screenshots, or deployment instructions.
