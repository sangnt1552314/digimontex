# DigimonTex

A terminal-based Digimon viewer application built with Go and the tview library. This project displays Digimon information from the [Digimon Wikipedia](https://en.wikipedia.org/wiki/Digimon) universe, powered by the [Digi-API](https://digi-api.com/).

## About

DigimonTex is a learning project created to explore Go programming and the tview terminal UI library. It provides an interactive terminal interface to browse and view detailed information about various Digimon characters.

## Features

- **Interactive Terminal UI**: Clean, mouse-enabled interface built with tview
- **Digimon Browser**: Browse through a paginated list of Digimon
- **Detailed Information**: View comprehensive details including:
  - Digimon images and field symbols
  - Name, release date, levels, types, and attributes
  - Detailed descriptions in English
  - Skills and abilities
- **Navigation**: Easy navigation with keyboard shortcuts and mouse support
- **Real-time Data**: Fetches live data from the Digi-API

## Technology Stack

- **Language**: Go 1.24.2
- **UI Library**: [tview](https://github.com/rivo/tview) - Terminal UI library
- **Terminal**: [tcell](https://github.com/gdamore/tcell) - Terminal handling
- **API**: [Digi-API](https://digi-api.com/) - Digimon data source

## Installation

1. Clone the repository:
```bash
git clone https://github.com/sangnt1552314/digimontex.git
cd digimontex
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run cmd/main.go
```

## Usage

- **Navigation**: Use arrow keys to navigate through the interface
- **Browse Digimon**: Use the left panel to browse through available Digimon
- **Pagination**: Use `<<` and `>>` buttons to navigate between pages
- **View Details**: Click on any Digimon name to view detailed information
- **Exit**: Press `Ctrl+C` or click the "⏻ Exit" button to quit

## Project Structure

```
digimontex/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   └── digimontex.go    # Main application logic and UI setup
│   ├── models/
│   │   └── digimon.go       # Data models for API responses
│   └── services/
│       ├── common.go        # Common utilities
│       └── digimon.go       # API service functions
├── assets/
│   └── no-image.png         # Fallback image for missing images
└── storage/
    └── logs/                # Application logs
```

## API Reference

This project uses the [Digi-API](https://digi-api.com/api/v1/) which provides:
- Digimon list with pagination
- Detailed Digimon information by name or ID
- High-quality images and comprehensive data

## Learning Goals

This project was created to:
- Learn Go programming fundamentals
- Explore terminal UI development with tview
- Practice API integration and HTTP clients
- Understand Go project structure and organization
- Work with image handling and display in terminals

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Digi-API](https://digi-api.com/) for providing the Digimon data
- [tview](https://github.com/rivo/tview) library for the excellent terminal UI framework
- [Digimon](https://en.wikipedia.org/wiki/Digimon) franchise for the amazing characters and universe# digimontex