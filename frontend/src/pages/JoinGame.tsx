import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Button from '../components/common/Button';

function JoinGame() {
  const [playerName, setPlayerName] = useState('');
  const [gameCode, setGameCode] = useState('');
  const navigate = useNavigate();

  const handleJoinGame = async (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement API call to join game
    navigate(`/lobby/${gameCode}`);
  };

  const handleCreateGame = async () => {
    // TODO: Implement API call to create game
    const newGameCode = 'GENERATED_CODE'; // This will come from the API
    navigate(`/lobby/${newGameCode}`);
  };

  return (
    <div className="flex flex-col items-center space-y-12 w-full px-6 md:px-0">
      {/* Logo Circle */}
      <div className="w-40 h-40 md:w-48 md:h-48 lg:w-56 lg:h-56 rounded-full bg-white/70 backdrop-blur-lg shadow-lg flex items-center justify-center border border-white/20">
        <div className="text-4xl md:text-5xl lg:text-6xl font-bold bg-gradient-to-r from-violet-500 to-fuchsia-500 text-transparent bg-clip-text logo-hover">
          TABOO
        </div>
      </div>

      {/* Form Container */}
      <div className="w-full max-w-xs sm:max-w-sm md:max-w-md lg:max-w-lg p-8 md:p-10 bg-white/70 backdrop-blur-lg rounded-3xl shadow-[0_8px_30px_rgb(0,0,0,0.06)] border border-white/20">
        <form onSubmit={handleJoinGame} className="space-y-8">
          <div className="space-y-4">
            <input
              type="text"
              id="playerName"
              placeholder="Player Name"
              value={playerName}
              onChange={(e) => setPlayerName(e.target.value)}
              className="w-full px-6 py-4 text-lg bg-white/50 border border-white/30 rounded-2xl 
                       focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-transparent
                       placeholder:text-gray-400 text-gray-600"
              required
            />
          </div>
          <div className="space-y-4">
            <input
              type="text"
              id="gameCode"
              placeholder="Game Code"
              value={gameCode}
              onChange={(e) => setGameCode(e.target.value)}
              className="w-full px-6 py-4 text-lg bg-white/50 border border-white/30 rounded-2xl 
                       focus:outline-none focus:ring-2 focus:ring-violet-500/50 focus:border-transparent
                       placeholder:text-gray-400 text-gray-600"
            />
          </div>
          <div className="space-y-4 pt-4">
            <Button type="submit" className="w-full py-4 text-lg bg-gradient-to-r from-violet-500 to-fuchsia-500 hover:from-violet-600 hover:to-fuchsia-600">
              Join Game
            </Button>
            <Button
              type="button"
              variant="secondary"
              className="w-full py-4 text-lg bg-white/50 hover:bg-white/70"
              onClick={handleCreateGame}
            >
              Create Game
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default JoinGame;