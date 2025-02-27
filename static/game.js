class Game2048 {
    constructor() {
        this.ws = new WebSocket(`ws://${window.location.host}/ws`);
        this.nextRpcId = 1;
        this.pendingRequests = new Map();
        
        this.tileContainer = document.querySelector('.tile-container');
        this.scoreDisplay = document.getElementById('score');
        this.bestScoreDisplay = document.getElementById('best-score');
        this.gameOverScreen = document.querySelector('.game-over');
        this.finalScoreDisplay = document.getElementById('final-score');

        this.ws.onmessage = (event) => this.handleServerMessage(event);
        this.ws.onopen = () => this.newGame();

        document.getElementById('new-game').addEventListener('click', () => this.newGame());
        document.getElementById('retry').addEventListener('click', () => this.newGame());
        document.addEventListener('keydown', (e) => this.handleInput(e));
    }

    async rpcCall(method, params = {}) {
        return new Promise((resolve, reject) => {
            const id = this.nextRpcId++;
            const request = {
                method,
                params,
                id
            };

            this.pendingRequests.set(id, { resolve, reject });
            this.ws.send(JSON.stringify(request));
        });
    }

    handleServerMessage(event) {
        const response = JSON.parse(event.data);
        const request = this.pendingRequests.get(response.id);
        if (request) {
            this.pendingRequests.delete(response.id);
            if (response.error) {
                request.reject(response.error);
            } else {
                request.resolve(response.result);
            }
        }
    }

    async newGame() {
        const gameState = await this.rpcCall('newGame');
        this.updateDisplay(gameState);
        this.gameOverScreen.classList.add('hidden');
    }

    async handleInput(e) {
        let direction;
        switch(e.key) {
            case 'ArrowUp': direction = 'up'; break;
            case 'ArrowDown': direction = 'down'; break;
            case 'ArrowLeft': direction = 'left'; break;
            case 'ArrowRight': direction = 'right'; break;
            default: return;
        }

        e.preventDefault();
        
        try {
            const result = await this.rpcCall('move', { direction });
            if (result.moved) {
                this.updateDisplay(result);
                if (result.gameOver) {
                    this.showGameOver();
                }
            }
        } catch (error) {
            console.error('Move failed:', error);
        }
    }

    updateDisplay(gameState) {
        this.tileContainer.innerHTML = '';
        
        for (let i = 0; i < 4; i++) {
            for (let j = 0; j < 4; j++) {
                const value = gameState.grid[i][j];
                if (value !== 0) {
                    this.addTile(i, j, value);
                }
            }
        }

        this.scoreDisplay.textContent = gameState.score;
        this.bestScoreDisplay.textContent = gameState.bestScore;
    }

    addTile(row, col, value) {
        const tile = document.createElement('div');
        tile.className = `tile tile-${value}`;
        tile.textContent = value;
        tile.style.transform = `translate(${col * 100 + col * 16}%, ${row * 100 + row * 16}%)`;
        this.tileContainer.appendChild(tile);
    }

    showGameOver() {
        this.finalScoreDisplay.textContent = this.scoreDisplay.textContent;
        this.gameOverScreen.classList.remove('hidden');
    }
}

// Запускаем игру после загрузки страницы
window.addEventListener('load', () => {
    new Game2048();
});
