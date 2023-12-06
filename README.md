# Go Stock Tracker

<p>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" height="25"/>
</p>

<p>
  <img src="./go.png" alt="Go Logo" width="250" align="right"/>
  <b>Project Overview</b><br>
  `Go Stock Tracker` is a Go-based application designed to track and display stock market data in real-time. Utilizing the Gin web framework for routing and WebSocket for real-time communication, this application fetches stock data from an external API and presents updated information at regular intervals.
</p>

### Features

- Real-time stock data tracking through WebSocket.
- Efficient web routing using Gin framework.
- Periodic updates of stock data at configurable intervals.
- Fetches data from external stock market APIs.
- Scalable architecture suitable for expanding to multiple stock symbols.

## Getting Started

### Prerequisites

- Go (1.x or later)
- External API key for stock data (e.g., Alpha Vantage)
- Git

### Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/your-username/go-stock-tracker.git
   ```
   API_KEY=<Your_Alpha_Vantage_API_Key>

go run app.go

## Testing WebSocket Endpoint with Postman

Postman allows you to test WebSocket connections and send messages. Here's how you can test the `/ws` WebSocket endpoint:

1. **Open Postman**

2. **Create a New WebSocket Request**

   - Click on the `New` button.
   - Select `WebSocket Request`.

3. **Configure the WebSocket URL**

   - Enter the WebSocket URL in the format: `ws://localhost:8080/ws` (adjust the domain and port as per your server configuration).

4. **Connect to the WebSocket**

   - Click `Connect` to establish a WebSocket connection.

5. **Send a Message**

   - To request stock data, send a JSON message in the following format:
     ```json
     {
       "action": "getStock",
       "ticker": "YOUR_STOCK_TICKER"
     }
     ```
   - Replace `YOUR_STOCK_TICKER` with the desired stock symbol (e.g., `IBM` or `AAPL`).

6. **Receive Updates**

   - After sending the message, you should receive periodic updates with the latest stock data.
   - The updates will appear in the `Messages` section of Postman.

7. **Close the Connection**
   - Once done, you can close the WebSocket connection by clicking the `Disconnect` button.

Note: Ensure your Go Stock Tracker application is running before testing with Postman.

Note:

- I have removed the sections commented out with `<!-- -->` as they might be related to features not currently implemented.
- The WebSocket section is now more focused on the real-time aspect of your application.
- Remember to replace `https://github.com/your-username/go-stock-tracker.git` with your actual repository URL.
- Ensure that the `.env` file and your API key are appropriately configured as per your project setup.
