import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/common/Layout';
import JoinGame from './pages/JoinGame';
import GameLobby from './pages/GameLobby';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<JoinGame />} />
          <Route path="/lobby/:gameId" element={<GameLobby />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
