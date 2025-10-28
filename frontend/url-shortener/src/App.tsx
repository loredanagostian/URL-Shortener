import { useState } from "react";
import "./App.css";

function App() {
  const [url, setUrl] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement URL shortening logic
    console.log("Shortening URL:", url);
  };

  return (
    <div className="app">
      <div className="container">
        <h1>URL Shortener</h1>
        <form onSubmit={handleSubmit} className="url-form">
          <input
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="Enter your URL here..."
            className="url-input"
            required
          />
          <button type="submit" className="shorten-btn">
            Shorten me
          </button>
        </form>
      </div>
    </div>
  );
}

export default App;
