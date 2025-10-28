import { useState, useEffect } from "react";
import "./App.css";
import URLForm from "./components/URLForm";
import ResultDisplay from "./components/ResultDisplay";
import type { ShortenResponse } from "./services/api";
import testAPI from "./test-api";

function App() {
  const [result, setResult] = useState<ShortenResponse | null>(null);

  // Test API connection on component mount
  useEffect(() => {
    testAPI();
  }, []);

  const handleSuccess = (shortenResult: ShortenResponse) => {
    setResult(shortenResult);
  };

  const handleReset = () => {
    setResult(null);
  };

  return (
    <div className="app">
      <div className="container">
        <header className="app-header">
          <h1>🔗 URL Shortener</h1>
          <p>Transform long URLs into short, shareable links</p>
        </header>

        <main className="app-main">
          {result ? (
            <ResultDisplay result={result} onReset={handleReset} />
          ) : (
            <URLForm onSuccess={handleSuccess} />
          )}
        </main>

        <footer className="app-footer">
          <p>Built with React & Go • Made with ❤️</p>
        </footer>
      </div>
    </div>
  );
}

export default App;
