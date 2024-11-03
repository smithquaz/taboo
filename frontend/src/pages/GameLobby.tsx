import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Button from '../components/common/Button';
import LoadingSpinner from '../components/common/LoadingSpinner';

interface Player {
  id: string;
  name: string;
  team?: 'A' | 'B';
}

function GameLobby() {
  const { gameId } = useParams();
  const [players, setPlayers] = useState<Player[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // TODO: Implement WebSocket connection to get real-time updates
    // For now, we'll just simulate some players
    setPlayers([
      { id: '1', name: 'Player 1', team: 'A' },
      { id: '2', name: 'Player 2', team: 'B' },
    ]);
    setIsLoading(false);
  }, [gameId]);

  if (isLoading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-lg font-medium text-gray-900 mb-4">Game Lobby</h2>
        <div className="text-sm text-gray-500 mb-4">
          Game Code: <span className="font-mono font-bold">{gameId}</span>
        </div>
        
        <div className="grid grid-cols-2 gap-8">
          <div>
            <h3 className="text-md font-medium text-gray-900 mb-2">Team A</h3>
            <ul className="space-y-2">
              {players
                .filter((p) => p.team === 'A')
                .map((player) => (
                  <li key={player.id} className="text-gray-700">
                    {player.name}
                  </li>
                ))}
            </ul>
          </div>
          <div>
            <h3 className="text-md font-medium text-gray-900 mb-2">Team B</h3>
            <ul className="space-y-2">
              {players
                .filter((p) => p.team === 'B')
                .map((player) => (
                  <li key={player.id} className="text-gray-700">
                    {player.name}
                  </li>
                ))}
            </ul>
          </div>
        </div>
      </div>
      
      <div className="flex justify-end">
        <Button>Start Game</Button>
      </div>
    </div>
  );
}

export default GameLobby; 