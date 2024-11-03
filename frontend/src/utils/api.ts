const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export async function createGame(playerName: string): Promise<string> {
  const response = await fetch(`${API_BASE_URL}/games`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ playerName }),
  });
  
  if (!response.ok) {
    throw new Error('Failed to create game');
  }
  
  const data = await response.json();
  return data.gameId;
}

export async function joinGame(gameId: string, playerName: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/games/${gameId}/join`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ playerName }),
  });
  
  if (!response.ok) {
    throw new Error('Failed to join game');
  }
}

export async function getGameState(gameId: string): Promise<any> {
  const response = await fetch(`${API_BASE_URL}/games/${gameId}`);
  
  if (!response.ok) {
    throw new Error('Failed to get game state');
  }
  
  return response.json();
} 