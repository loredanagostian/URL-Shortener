import { useState } from 'react';
import './App.css';
import URLForm from './components/URLForm';
import ResultDisplay from './components/ResultDisplay';
import History from './components/History';
import type { ShortenResponse } from './services/api';

type ViewType = 'form' | 'result' | 'history';

function App() {
    const [result, setResult] = useState<ShortenResponse | null>(null);
    const [currentView, setCurrentView] = useState<ViewType>('form');

    const handleSuccess = (shortenResult: ShortenResponse) => {
        setResult(shortenResult);
        setCurrentView('result');
    };

    const handleReset = () => {
        setResult(null);
        setCurrentView('form');
    };

    const handleViewChange = (view: ViewType) => {
        setCurrentView(view);
    };

    const renderContent = () => {
        switch (currentView) {
            case 'result':
                return result ? <ResultDisplay result={result} onReset={handleReset} /> : null;
            case 'history':
                return <History />;
            default:
                return <URLForm onSuccess={handleSuccess} />;
        }
    };

    return (
        <div className="app">
            <div className={`container ${currentView === 'history' ? 'history-view' : ''}`}>
                <header className="app-header">
                    <h1>🔗 URL Shortener</h1>
                    <p>Transform long URLs into short, shareable links</p>

                    <nav className="app-nav">
                        <button
                            className={`nav-btn ${currentView === 'form' ? 'active' : ''}`}
                            onClick={() => handleViewChange('form')}
                        >
                            Create Short URL
                        </button>
                        <button
                            className={`nav-btn ${currentView === 'history' ? 'active' : ''}`}
                            onClick={() => handleViewChange('history')}
                        >
                            History
                        </button>
                    </nav>
                </header>

                <main className="app-main">{renderContent()}</main>

                <footer className="app-footer">
                    <p>Built with React & Go • Made with ❤️</p>
                </footer>
            </div>
        </div>
    );
}

export default App;
