<script>
	import { onMount, onDestroy } from 'svelte';
	
	let ws;
	let rivalName = 'Esperando...';
	let rivalRank = '';
	let isInGame = false;
	let matchStarted = false;
	let connectionStatus = 'Desconectado';
	let gifUrl = '/gifs/pet.gif'; // GIF por defecto desde carpeta local
	
	// Detectar si estamos en Electron
	onMount(() => {
		isElectron = navigator.userAgent.toLowerCase().indexOf('electron') > -1;
	});
	
	// Conectar al WebSocket
	function connectWebSocket() {
		ws = new WebSocket('ws://localhost:8080/ws');
		
		ws.onopen = () => {
			console.log('Conectado al servidor');
			connectionStatus = 'Conectado';
		};
		
		ws.onmessage = (event) => {
			const data = JSON.parse(event.data);
			rivalName = data.rivalName || 'Desconocido';
			rivalRank = data.rivalRank || '';
			isInGame = data.isInGame;
			matchStarted = data.matchStarted;
		};
		
		ws.onerror = (error) => {
			console.error('WebSocket error:', error);
			connectionStatus = 'Error';
		};
		
		ws.onclose = () => {
			console.log('Desconectado del servidor');
			connectionStatus = 'Desconectado';
			// Intentar reconectar después de 3 segundos
			setTimeout(connectWebSocket, 3000);
		};
	}
	
	// Fetch alternativo si WebSocket no funciona
	async function fetchGameData() {
		try {
			const response = await fetch('http://localhost:8080/api/gamedata');
			const data = await response.json();
			rivalName = data.rivalName || 'Desconocido';
			rivalRank = data.rivalRank || '';
			isInGame = data.isInGame;
			matchStarted = data.matchStarted;
		} catch (error) {
			console.error('Error fetching game data:', error);
		}
	}
	
	onMount(() => {
		connectWebSocket();
		// Backup: polling cada 5 segundos
		const interval = setInterval(fetchGameData, 5000);
		
		return () => {
			clearInterval(interval);
		};
	});
	
	onDestroy(() => {
		if (ws) {
			ws.close();
		}
	});
</script>

<main class="overlay-container">
	<!-- Panel principal con info del rival -->
	<div class="rival-info" class:in-game={isInGame}>
		<div class="status-indicator" class:active={matchStarted}>
			{matchStarted ? '🎮' : '⏸️'}
		</div>
		
		<div class="rival-details">
			<h2 class="rival-name">{rivalName}</h2>
			{#if rivalRank}
				<p class="rival-rank">{rivalRank}</p>
			{/if}
		</div>
		
		{#if isInGame && matchStarted}
			<div class="game-status">
				<span class="pulse-dot"></span>
				EN PARTIDA
			</div>
		{/if}
	</div>
	
	<!-- GIF animado en la esquina -->
	<div class="gif-container">
		<img src={gifUrl} alt="Animated GIF" class="overlay-gif" />
	</div>
	
	<!-- Indicador de conexión -->
	<div class="connection-status" class:connected={connectionStatus === 'Conectado'}>
		{connectionStatus}
	</div>
</main>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		overflow: hidden;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}
	
	.overlay-container {
		width: 100vw;
		height: 100vh;
		background: transparent;
		position: relative;
		pointer-events: none; /* IMPORTANTE: hace que el fondo no bloquee clicks */
	}
	
	.rival-info {
		position: absolute;
		top: 20px;
		left: 20px;
		background: rgba(0, 0, 0, 0.85);
		backdrop-filter: blur(10px);
		border: 2px solid #333;
		border-radius: 12px;
		padding: 20px 30px;
		min-width: 300px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
		transition: all 0.3s ease;
		pointer-events: auto;
	}
	
	.rival-info.in-game {
		border-color: #00ff88;
		box-shadow: 0 8px 32px rgba(0, 255, 136, 0.3);
	}
	
	.status-indicator {
		position: absolute;
		top: 10px;
		right: 10px;
		font-size: 24px;
	}
	
	.status-indicator.active {
		animation: bounce 1s infinite;
	}
	
	.rival-details {
		margin-bottom: 10px;
	}
	
	.rival-name {
		color: #fff;
		font-size: 28px;
		font-weight: bold;
		margin: 0 0 8px 0;
		text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.8);
	}
	
	.rival-rank {
		color: #00ff88;
		font-size: 18px;
		margin: 0;
		font-weight: 600;
	}
	
	.game-status {
		display: flex;
		align-items: center;
		gap: 8px;
		color: #00ff88;
		font-weight: bold;
		font-size: 14px;
		margin-top: 10px;
	}
	
	.pulse-dot {
		width: 10px;
		height: 10px;
		background: #00ff88;
		border-radius: 50%;
		animation: pulse 1.5s infinite;
	}
	
	.gif-container {
		position: absolute;
		bottom: 40px;
		right: 40px;
		pointer-events: none;
	}
	
	.overlay-gif {
		width: 80px;
		height: 80px;
		object-fit: contain;
		border-radius: 12px;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
		animation: float 3s ease-in-out infinite;
	}
	.connection-status {
		position: absolute;
		top: 20px;
		right: 20px;
		background: rgba(255, 0, 0, 0.8);
		color: white;
		padding: 8px 16px;
		border-radius: 8px;
		font-size: 12px;
		font-weight: bold;
		pointer-events: none;
	}
	
	.connection-status.connected {
		background: rgba(0, 255, 136, 0.8);
	}
	
	/* Animaciones */
	@keyframes pulse {
		0%, 100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.5;
			transform: scale(1.3);
		}
	}
	
	@keyframes bounce {
		0%, 100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-10px);
		}
	}
	
	@keyframes float {
		0%, 100% {
			transform: translateY(0px);
		}
		50% {
			transform: translateY(-20px);
		}
	}
</style>