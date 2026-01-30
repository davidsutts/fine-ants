# 🐜 Fine-Ants: Finance Done Properly

Fine-Ants is a minimalist, high-performance personal finance tracker designed to help you categorize and visualize your spending without the bloat. Built with a lightning-fast **Go (Fiber)** backend and a reactive **Svelte 5** frontend.

## ✨ Features

- **CSV Upload:** Drag and drop your bank statements to get started instantly.
- **Dual Modes:**
  - **List Mode:** A clean, chronological overview of all your transactions and running balances.
  - **Sort Mode:** An interactive, rapid-fire categorization interface to move through your "WIP" transactions using intuitive emoji-based categories.
- **Localized for AU:** Built-in support for AUD currency formatting and Australian date standards.
- **Privacy First:** Your financial data stays under your control, processed by your own local instance.

## 🛠️ Tech Stack

- **Frontend:** Svelte 5 (Runes), Tailwind CSS
- **Backend:** Go, Fiber Web Framework
- **Storage:** Integrated Datastore support

---

## 🚀 Getting Started

### Prerequisites
- [Go](https://go.dev/doc/install) (1.21+)
- [Node.js & npm](https://nodejs.org/)
- [Vite](https://vitejs.dev/)

### 1. Setup the Backend
Navigate to the root directory and install dependencies:
```bash
go mod tidy

```

Run the server:

```bash
# Default port 8080
go run main.go

# Or run in dev mode with a custom port
go run main.go -port 3000 -dev

```

### 2. Setup the Frontend

Navigate to your webapp directory:

```bash
cd webapp
npm install
npm run dev

```

The frontend will typically be available at `http://localhost:5173`.

## 📖 Roadmap & "What it could be"

Fine-Ants is currently in its early stages (`v0.0.1`). Future plans include:

* [ ] **AI-Assisted Auto-Categorization:** Learning from your previous sorts to suggest categories for new uploads.
* [ ] **Data Visualization:** Monthly spend charts and category breakdowns.
* [ ] **Export Options:** Clean PDF or CSV exports of categorized data for tax time.
* [ ] **Multi-Bank Support:** Pre-configured parsers for major Australian financial institutions.

---

## 📄 License

Copyright (C) 2026 - **David Sutton** ([dsutton1202@gmail.com](mailto:dsutton1202@gmail.com))
Distributed under the **GNU General Public License v3**. See `LICENSE` for more information.

---

*Finance shouldn't be a chore. Do it properly with Fine-Ants.*
