# Go Stock Tracker

<p>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" height="25"/>
</p>

<p>
  <img src="./go.png" alt="Go Logo" width="250" align="right"/>
  <b>Project Overview</b><br>
  `Go Stock Tracker` is a Go-based application designed to track and display stock market data in real-time. Utilizing the Gin web framework for routing and WebSocket for real-time communication, this application fetches stock data from an external API and presents updated information periodically.
</p>

### Features

- Real-time stock data tracking.
- Utilizes Gin for efficient web routing.
- WebSocket implementation for real-time data updates.
- Fetches data from external stock market APIs.
- Periodically updates stock data at predefined intervals.

## Getting Started

### Prerequisites

- Go (1.x or later)
- External API key for stock data (e.g., Alpha Vantage)
- Git

### Installation

1. **Clone the Repository**
   git clone https://github.com/your-username/go-stock-tracker.git

2. **Configure Environment Variables**

Create a `.env` file in the root directory with your API key:

3. **Run the Application**
   go run main.go

## Usage

### WebSocket Endpoint

- **WebSocket Connection**: `GET /ws`
- Establishes a WebSocket connection for real-time data communication.

### Stock Data Fetching

- **Start Fetching Stocks**: `GET /startFetching`
- Initiates the process of fetching stock data.

### Example Requests

- Initiating Stock Data Fetching:
  curl http://localhost:8080/startFetching

## Contributing

Contributions to `Go Stock Tracker` are welcome. To contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.
